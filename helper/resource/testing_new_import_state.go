// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-version"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func testStepNewImportState(ctx context.Context, t testing.T, helper *plugintest.Helper, wd *plugintest.WorkingDir, step TestStep, cfgRaw string, providers *providerFactories, stepNumber int) error {
	t.Helper()

	// step.ImportStateKind implicitly defaults to the zero-value (ImportCommandWithID) for backward compatibility
	kind := step.ImportStateKind
	if kind.plannable() {
		// Instead of calling [t.Fatal], return an error. This package's unit tests can use [TestStep.ExpectError] to match on the error message.
		// An alternative, [plugintest.TestExpectTFatal], does not have access to logged error messages, so it is open to false positives on this
		// complex code path.
		if err := requirePlannableImport(t, *helper.TerraformVersion()); err != nil {
			return err
		}
	}

	configRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       step.Config,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	testStepConfig := teststep.Configuration(configRequest)

	resourceName := step.ResourceName
	if resourceName == "" {
		t.Fatal("ResourceName is required for an import state test")
	}

	// get state from check sequence
	var state *terraform.State
	var stateJSON *tfjson.State
	var err error

	err = runProviderCommand(ctx, t, func() error {
		stateJSON, state, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// TODO: this statement is a placeholder -- it simply prevents stateJSON from being unused
	logging.HelperResourceTrace(ctx, fmt.Sprintf("State before import: values %v", stateJSON.Values != nil))

	// Determine the ID to import
	var importId string
	switch {
	case step.ImportStateIdFunc != nil:
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateIdFunc for import identifier")

		var err error

		logging.HelperResourceDebug(ctx, "Calling TestStep ImportStateIdFunc")

		importId, err = step.ImportStateIdFunc(state)

		if err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep ImportStateIdFunc")
	case step.ImportStateId != "":
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateId for import identifier")

		importId = step.ImportStateId
	default:
		logging.HelperResourceTrace(ctx, "Using resource identifier for import identifier")

		resource, err := testResource(resourceName, state)
		if err != nil {
			t.Fatal(err)
		}
		importId = resource.Primary.ID
	}

	if step.ImportStateIdPrefix != "" {
		logging.HelperResourceTrace(ctx, "Prepending TestStep ImportStateIdPrefix for import identifier")

		importId = step.ImportStateIdPrefix + importId
	}

	logging.HelperResourceTrace(ctx, fmt.Sprintf("Using import identifier: %s", importId))

	// Append to previous step config unless using ConfigFile or ConfigDirectory
	if testStepConfig == nil || step.Config != "" {
		importConfig := step.Config
		if importConfig == "" {
			logging.HelperResourceTrace(ctx, "Using prior TestStep Config for import")
			importConfig = cfgRaw
		}

		if kind.plannable() {
			importConfig = appendImportBlock(importConfig, resourceName, importId)
		}

		confRequest := teststep.PrepareConfigurationRequest{
			Directory: step.ConfigDirectory,
			File:      step.ConfigFile,
			Raw:       importConfig,
			TestStepConfigRequest: config.TestStepConfigRequest{
				StepNumber: stepNumber,
				TestName:   t.Name(),
			},
		}.Exec()

		testStepConfig = teststep.Configuration(confRequest)
		if testStepConfig == nil {
			t.Fatal("Cannot import state with no specified config")
		}
	}

	var importWd *plugintest.WorkingDir

	// Use the same working directory to persist the state from import
	if step.ImportStatePersist {
		importWd = wd
	} else {
		importWd = helper.RequireNewWorkingDir(ctx, t, "")
		defer importWd.Close() //nolint:errcheck
	}

	err = importWd.SetConfig(ctx, testStepConfig, step.ConfigVariables)
	if err != nil {
		t.Fatalf("Error setting test config: %s", err)
	}

	if !step.ImportStatePersist {
		err = runProviderCommand(ctx, t, func() error {
			logging.HelperResourceDebug(ctx, "Run terraform init")
			return importWd.Init(ctx)
		}, importWd, providers)
		if err != nil {
			t.Fatalf("Error running init: %s", err)
		}
	}

	var plan *tfjson.Plan
	if kind.plannable() {
		var opts []tfexec.PlanOption

		err = runProviderCommand(ctx, t, func() error {
			logging.HelperResourceDebug(ctx, "Run terraform plan")
			return importWd.CreatePlan(ctx, opts...)
		}, importWd, providers)
		if err != nil {
			return err
		}

		err = runProviderCommand(ctx, t, func() error {
			var err error
			logging.HelperResourceDebug(ctx, "Run terraform show")
			plan, err = importWd.SavedPlan(ctx)
			return err
		}, importWd, providers)
		if err != nil {
			return err
		}

		if plan.ResourceChanges != nil {
			logging.HelperResourceDebug(ctx, fmt.Sprintf("ImportBlockWithId: %d resource changes", len(plan.ResourceChanges)))

			for _, rc := range plan.ResourceChanges {
				if rc.Address != resourceName {
					// we're only interested in the changes for the resource being imported
					continue
				}
				if rc.Change != nil && rc.Change.Actions != nil {
					// should this be length checked and used as a condition, if it's a no-op then there shouldn't be any other changes here
					for _, action := range rc.Change.Actions {
						if action != "no-op" {
							var stdout string
							err = runProviderCommand(ctx, t, func() error {
								var err error
								stdout, err = importWd.SavedPlanRawStdout(ctx)
								return err
							}, importWd, providers)
							if err != nil {
								return fmt.Errorf("retrieving formatted plan output: %w", err)
							}

							return fmt.Errorf("importing resource %s: expected a no-op resource action, got %q action with plan \nstdout:\n\n%s", rc.Address, action, stdout)
						}
					}
				}
			}
		}

		// TODO compare plan to state from previous step

		if err := runPlanChecks(ctx, t, plan, step.ImportPlanChecks.PreApply); err != nil {
			return err
		}
	} else {
		err = runProviderCommand(ctx, t, func() error {
			return importWd.Import(ctx, resourceName, importId)
		}, importWd, providers)
		if err != nil {
			return err
		}
	}

	var importState *terraform.State
	err = runProviderCommand(ctx, t, func() error {
		_, importState, err = getState(ctx, t, importWd)
		if err != nil {
			return err
		}
		return nil
	}, importWd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	logging.HelperResourceDebug(ctx, fmt.Sprintf("State after import: %d resources in the root module", len(importState.RootModule().Resources)))

	// Go through the imported state and verify
	if step.ImportStateCheck != nil {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateCheck")
		runImportStateCheckFunction(ctx, t, importState, step)
	}

	// Verify that all the states match
	if step.ImportStateVerify {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateVerify")

		// Ensure that we do not match against data sources as they
		// cannot be imported and are not what we want to verify.
		// Mode is not present in ResourceState so we use the
		// stringified ResourceStateKey for comparison.
		newResources := make(map[string]*terraform.ResourceState)
		for k, v := range importState.RootModule().Resources {
			if !strings.HasPrefix(k, "data.") {
				newResources[k] = v
			}
		}
		oldResources := make(map[string]*terraform.ResourceState)
		for k, v := range state.RootModule().Resources {
			if !strings.HasPrefix(k, "data.") {
				oldResources[k] = v
			}
		}

		identifierAttribute := step.ImportStateVerifyIdentifierAttribute

		if identifierAttribute == "" {
			identifierAttribute = "id"
		}

		for _, r := range newResources {
			rIdentifier, ok := r.Primary.Attributes[identifierAttribute]

			if !ok {
				t.Fatalf("ImportStateVerify: New resource missing identifier attribute %q, ensure attribute value is properly set or use ImportStateVerifyIdentifierAttribute to choose different attribute", identifierAttribute)
			}

			// Find the existing resource
			var oldR *terraform.ResourceState
			for _, r2 := range oldResources {
				if r2.Primary == nil || r2.Type != r.Type || r2.Provider != r.Provider {
					continue
				}

				r2Identifier, ok := r2.Primary.Attributes[identifierAttribute]

				if !ok {
					t.Fatalf("ImportStateVerify: Old resource missing identifier attribute %q, ensure attribute value is properly set or use ImportStateVerifyIdentifierAttribute to choose different attribute", identifierAttribute)
				}

				if r2Identifier == rIdentifier {
					oldR = r2
					break
				}
			}
			if oldR == nil || oldR.Primary == nil {
				t.Fatalf(
					"Failed state verification, resource with ID %s not found",
					rIdentifier)
			}

			// don't add empty flatmapped containers, so we can more easily
			// compare the attributes
			skipEmpty := func(k, v string) bool {
				if strings.HasSuffix(k, ".#") || strings.HasSuffix(k, ".%") {
					if v == "0" {
						return true
					}
				}
				return false
			}

			// Compare their attributes
			actual := make(map[string]string)
			for k, v := range r.Primary.Attributes {
				if skipEmpty(k, v) {
					continue
				}
				actual[k] = v
			}

			expected := make(map[string]string)
			for k, v := range oldR.Primary.Attributes {
				if skipEmpty(k, v) {
					continue
				}
				expected[k] = v
			}

			// Remove fields we're ignoring
			for _, v := range step.ImportStateVerifyIgnore {
				for k := range actual {
					if strings.HasPrefix(k, v) {
						delete(actual, k)
					}
				}
				for k := range expected {
					if strings.HasPrefix(k, v) {
						delete(expected, k)
					}
				}
			}

			// timeouts are only _sometimes_ added to state. To
			// account for this, just don't compare timeouts at
			// all.
			for k := range actual {
				if strings.HasPrefix(k, "timeouts.") {
					delete(actual, k)
				}
				if k == "timeouts" {
					delete(actual, k)
				}
			}
			for k := range expected {
				if strings.HasPrefix(k, "timeouts.") {
					delete(expected, k)
				}
				if k == "timeouts" {
					delete(expected, k)
				}
			}

			if !reflect.DeepEqual(actual, expected) {
				// Determine only the different attributes
				// go-cmp tries to show surrounding identical map key/value for
				// context of differences, which may be confusing.
				for k, v := range expected {
					if av, ok := actual[k]; ok && v == av {
						delete(expected, k)
						delete(actual, k)
					}
				}

				if diff := cmp.Diff(expected, actual); diff != "" {
					return fmt.Errorf("ImportStateVerify attributes not equivalent. Difference is shown below. The - symbol indicates attributes missing after import.\n\n%s", diff)
				}
			}
		}
	}

	return nil
}

func appendImportBlock(config string, resourceName string, importID string) string {
	return config + fmt.Sprintf(``+"\n"+
		`import {`+"\n"+
		`	to = %s`+"\n"+
		`	id = %q`+"\n"+
		`}`,
		resourceName, importID)
}

func requirePlannableImport(t testing.T, versionUnderTest version.Version) error {
	t.Helper()

	if versionUnderTest.LessThan(tfversion.Version1_5_0) {
		return fmt.Errorf(
			`ImportState steps using plannable import blocks require Terraform 1.5.0 or later. Either ` +
				`upgrade the Terraform version running the test or add a ` + "`TerraformVersionChecks`" + ` to ` +
				`the test case to skip this test.` + "\n\n" +
				`https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/tfversion-checks#skip-version-checks`)
	}

	return nil
}

func runImportStateCheckFunction(ctx context.Context, t testing.T, importState *terraform.State, step TestStep) {
	t.Helper()

	var states []*terraform.InstanceState
	for address, r := range importState.RootModule().Resources {
		if strings.HasPrefix(address, "data.") {
			continue
		}

		if r.Primary == nil {
			continue
		}

		is := r.Primary.DeepCopy() //nolint:staticcheck // legacy usage
		is.Ephemeral.Type = r.Type // otherwise the check function cannot see the type
		states = append(states, is)
	}

	logging.HelperResourceTrace(ctx, "Calling TestStep ImportStateCheck")

	if err := step.ImportStateCheck(states); err != nil {
		t.Fatal(err)
	}

	logging.HelperResourceTrace(ctx, "Called TestStep ImportStateCheck")
}
