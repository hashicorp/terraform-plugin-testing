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
)

func requirePlannableImport(t testing.T, versionUnderTest version.Version) error {
	t.Helper()

	minVersion, err := version.NewVersion("1.5.0")
	if err != nil {
		panic("failed to parse version string")
	}

	if versionUnderTest.LessThan(minVersion) {
		return fmt.Errorf("ImportState steps require Terraform 1.5.0 or later")
	}

	return nil
}

func testStepNewImportState(ctx context.Context, t testing.T, helper *plugintest.Helper, wd *plugintest.WorkingDir, step TestStep, cfgRaw string, providers *providerFactories, stepNumber int) error {
	t.Helper()

	// Currently import modes `ImportBlockWithId` and `ImportBlockWithResourceIdentity` cannot support config file or directory
	// since these modes append the import block to the configuration automatically
	if step.ImportStateKind != ImportCommandWithId {
		if step.ConfigFile != nil || step.ConfigDirectory != nil {
			t.Fatalf("ImportStateKind %q is not supported for config file or directory", step.ImportStateKind)
		}
	}

	{
		err := requirePlannableImport(t, *helper.TerraformVersion())
		if err != nil {
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

	if step.ResourceName == "" {
		t.Fatal("ResourceName is required for an import state test")
	}

	// get state from check sequence
	var state *terraform.State
	var err error

	err = runProviderCommand(ctx, t, func() error {
		state, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

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

		resource, err := testResource(step, state)
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

	if testStepConfig == nil || step.Config != "" {
		importConfig := step.Config
		if importConfig == "" {
			logging.HelperResourceTrace(ctx, "Using prior TestStep Config for import")
			importConfig = cfgRaw
		}

		// Update the test config dependent on the kind of import test being performed
		switch step.ImportStateKind {
		case ImportBlockWithResourceIdentity:
			t.Fatalf("TODO implement me")
		case ImportBlockWithId:
			importConfig += fmt.Sprintf(`
			import {
				to = %s
				id = %q
			}
		`, step.ResourceName, importId)
		default:
			// Not an import block test so nothing to do here
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
		defer importWd.Close()
	}

	err = importWd.SetConfig(ctx, testStepConfig, step.ConfigVariables)
	if err != nil {
		t.Fatalf("Error setting test config: %s", err)
	}

	logging.HelperResourceDebug(ctx, "Running Terraform CLI init and import")

	if !step.ImportStatePersist {
		err = runProviderCommand(ctx, t, func() error {
			return importWd.Init(ctx)
		}, importWd, providers)
		if err != nil {
			t.Fatalf("Error running init: %s", err)
		}
	}

	if step.ImportStateKind == ImportBlockWithResourceIdentity || step.ImportStateKind == ImportBlockWithId {
		var opts []tfexec.PlanOption

		err = runProviderCommand(ctx, t, func() error {
			return importWd.CreatePlan(ctx, opts...)
		}, importWd, providers)
		if err != nil {
			return err
		}

		var plan *tfjson.Plan
		err = runProviderCommand(ctx, t, func() error {
			var err error
			plan, err = importWd.SavedPlan(ctx)
			return err
		}, importWd, providers)
		if err != nil {
			return err
		}

		if plan.ResourceChanges != nil {
			for _, rc := range plan.ResourceChanges {
				if rc.Address != step.ResourceName {
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

							return fmt.Errorf("importing resource %s should be a no-op, but got action %s with plan \\nstdout:\\n\\n%s", rc.Address, action, stdout)
						}
					}
				}
			}
		}

		// TODO compare plan to state from previous step
	} else {
		err = runProviderCommand(ctx, t, func() error {
			return importWd.Import(ctx, step.ResourceName, importId)
		}, importWd, providers)
		if err != nil {
			return err
		}
	}

	var importState *terraform.State
	err = runProviderCommand(ctx, t, func() error {
		importState, err = getState(ctx, t, importWd)
		if err != nil {
			return err
		}
		return nil
	}, importWd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// Go through the imported state and verify
	if step.ImportStateCheck != nil {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateCheck")

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

		logging.HelperResourceDebug(ctx, "Calling TestStep ImportStateCheck")

		if err := step.ImportStateCheck(states); err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep ImportStateCheck")
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
