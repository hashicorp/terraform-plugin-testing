// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mitchellh/go-testing-interface"
)

func testStepNewBackendSmokeTest(ctx context.Context, t testing.T, c TestCase, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories, stepIndex int, helper *plugintest.Helper) error {
	t.Helper()

	stepNumber := stepIndex + 1

	configRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       step.Config,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	cfg := teststep.Configuration(configRequest)

	var hasTerraformBlock bool
	var hasProviderBlock bool

	if cfg != nil {
		var err error

		hasTerraformBlock, err = cfg.HasTerraformBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains terraform block",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains terraform block: %s", err)
		}

		hasProviderBlock, err = cfg.HasProviderBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains provider block",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains provider block: %s", err)
		}
	}

	// TODO: currently, this won't write the provider block because the terraform block exists (backends), since
	// it needs to write to terraform -> required_providers.
	//
	// - Honestly just need to refactor everything above this, who knows how much is actually needed.
	backendConfig, err := step.mergedConfig(ctx, c, hasTerraformBlock, hasProviderBlock, helper.TerraformVersion())

	if err != nil {
		logging.HelperResourceError(ctx,
			"Error generating merged configuration",
			map[string]interface{}{logging.KeyError: err},
		)
		t.Fatalf("Error generating merged configuration: %s", err)
	}

	confRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       backendConfig,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	testStepConfig := teststep.Configuration(confRequest)

	err = wd.SetConfig(ctx, testStepConfig, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// 1. Validate and configure the backend
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	// 2. Grab all the workspaces from the backend
	workspaces := make([]string, 0)
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		workspaces, err = wd.Workspaces(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Error getting workspaces: %w", err)
	}

	// 3. Assert the only workspace is the default one
	//
	// TODO: apparently some backends don't support the default workspace, so maybe this
	// assertion needs to be controlled by the test step?
	//
	// https://github.com/hashicorp/terraform/blob/643266dc90523ae794b586ad42e8b30864b61aaa/internal/backend/testing.go#L81-L88
	if len(workspaces) != 1 || workspaces[0] != "default" {
		t.Fatalf("Expected only the default workspace when initialized: %#v", workspaces)
	}

	// 4. Create "foo" workspace and assert the state returned is empty
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.CreateWorkspace(ctx, "foo")
	})
	if err != nil {
		return fmt.Errorf("Error creating \"foo\" workspace: %w", err)
	}

	var fooState *tfjson.State
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		fooState, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"foo\" state: %w", err)
	}

	if fooState.Values != nil && fooState.Values.RootModule != nil && len(fooState.Values.RootModule.Resources) > 0 {
		t.Fatalf("Expected the newly created \"foo\" state to be empty. Found %d resources.", len(fooState.Values.RootModule.Resources))
	}

	// 5. Create "bar" workspace and assert the state returned is empty
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.CreateWorkspace(ctx, "bar")
	})
	if err != nil {
		return fmt.Errorf("Error creating \"bar\" workspace: %w", err)
	}

	var barState *tfjson.State
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		barState, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"bar\" state: %w", err)
	}

	if barState.Values != nil && barState.Values.RootModule != nil && len(barState.Values.RootModule.Resources) > 0 {
		t.Fatalf("Expected the newly created \"bar\" state to be empty. Found %d resources.", len(barState.Values.RootModule.Resources))
	}

	// 6. Add a fake resource to the bar workspace
	barConfig := backendConfig + `
resource "terraform_data" "tf_plugin_testing_resource_bar" {
  input = "this resource was injected by terraform-plugin-testing"
}`

	barOneResourceConfigRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       barConfig,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	barOneResourceConfig := teststep.Configuration(barOneResourceConfigRequest)

	err = wd.SetConfig(ctx, barOneResourceConfig, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// 6a. Apply
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Apply(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error creating fake resource in \"bar\" workspace: %w", err)
	}

	// 6b. Grab the bar state again
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		barState, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"bar\" state: %w", err)
	}

	// 6c. Check if the resource exists in the "bar" state (I'm lazy lol)
	checkOutput := statecheck.ExpectKnownValue(
		"terraform_data.tf_plugin_testing_resource_bar",
		tfjsonpath.New("output"),
		knownvalue.StringExact("this resource was injected by terraform-plugin-testing"),
	)

	checkResp := statecheck.CheckStateResponse{}
	checkOutput.CheckState(ctx, statecheck.CheckStateRequest{State: barState}, &checkResp)

	if checkResp.Error != nil {
		return fmt.Errorf("After writing a test resource instance object to \"bar\" and re-reading it, the object has vanished: %w", err)
	}

	// 7. Switch to the "foo" workspace and grab the state, ensuring it's still empty.
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "foo")
	})
	if err != nil {
		return fmt.Errorf("Error creating \"foo\" workspace: %w", err)
	}

	// TODO: Do I need to apply an empty state to foo? The test helper kind of does this, but unsure whether that's actually useful or needed
	// https://github.com/hashicorp/terraform/blob/643266dc90523ae794b586ad42e8b30864b61aaa/internal/backend/testing.go#L134-L140
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		fooState, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"foo\" state: %w", err)
	}

	if fooState.Values != nil && fooState.Values.RootModule != nil && len(fooState.Values.RootModule.Resources) > 0 {
		t.Fatalf("After writing a resource to \"bar\" state, expected the \"foo\" state to be empty. Found %d resources in \"foo\".", len(fooState.Values.RootModule.Resources))
	}

	// 8. Verify when we list the workspaces we get back "default", "foo" (created during this test), and "bar" (created during this test).
	workspaces = make([]string, 0)
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		workspaces, err = wd.Workspaces(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Error getting workspaces: %w", err)
	}

	slices.Sort(workspaces)
	expected := []string{"bar", "default", "foo"}
	if !slices.Equal(expected, workspaces) {
		t.Fatalf("Expected workspaces to be %#v, got: %#v", expected, workspaces)
	}

	// 9. Delete "bar" workspace (the original helper deletes foo, but that is already empty)
	// Switch to "foo" workspace so we can actually delete bar.
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "foo")
	})
	if err != nil {
		return fmt.Errorf("Error selecting \"foo\" workspace: %w", err)
	}
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, "bar", tfexec.Force(true))
	})
	if err != nil {
		return fmt.Errorf("Error deleting \"bar\" workspace: %w", err)
	}

	// 10. Attempt to delete "default" workspace, assert error
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, "default", tfexec.Force(true))
	})
	if err == nil {
		return errors.New("Expected error when deleting \"default\" workspace")
	}

	// 11. Recreate the "bar" workspace
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.CreateWorkspace(ctx, "bar")
	})
	if err != nil {
		return fmt.Errorf("Error creating \"bar\" workspace: %w", err)
	}

	// 12. Grab "bar" state and assert it is empty (i.e. no left over artifacts)
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		barState, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"bar\" state: %w", err)
	}

	if barState.Values != nil && barState.Values.RootModule != nil && len(barState.Values.RootModule.Resources) > 0 {
		t.Fatalf("Expected the newly created \"bar\" state to be empty. Found %d resources.", len(barState.Values.RootModule.Resources))
	}

	// 13. Delete "bar" workspace again, force=true
	// Switch to "foo" workspace so we can actually delete bar.
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "foo")
	})
	if err != nil {
		return fmt.Errorf("Error selecting \"foo\" workspace: %w", err)
	}
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, "bar", tfexec.Force(true))
	})
	if err != nil {
		return fmt.Errorf("Error deleting \"bar\" workspace: %w", err)
	}

	// 14. List workspaces and verify it's just "foo" and "default"
	workspaces = make([]string, 0)
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		workspaces, err = wd.Workspaces(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Error getting workspaces: %w", err)
	}

	slices.Sort(workspaces)
	expected = []string{"default", "foo"}
	if !slices.Equal(expected, workspaces) {
		t.Fatalf("Expected workspaces to be %#v, got: %#v", expected, workspaces)
	}

	return nil
}

func testStepNewBackendLockTest(ctx context.Context, t testing.T, c TestCase, _ *plugintest.WorkingDir, step TestStep, providers *providerFactories, stepIndex int, helper *plugintest.Helper) error {
	t.Helper()

	stepNumber := stepIndex + 1

	configRequest := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       step.Config,
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	cfg := teststep.Configuration(configRequest)

	var hasTerraformBlock bool
	var hasProviderBlock bool

	if cfg != nil {
		var err error

		hasTerraformBlock, err = cfg.HasTerraformBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains terraform block",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains terraform block: %s", err)
		}

		hasProviderBlock, err = cfg.HasProviderBlock(ctx)

		if err != nil {
			logging.HelperResourceError(ctx,
				"Error determining whether configuration contains provider block",
				map[string]interface{}{logging.KeyError: err},
			)
			t.Fatalf("Error determining whether configuration contains provider block: %s", err)
		}
	}

	// TODO: currently, this won't write the provider block because the terraform block exists (backends), since
	// it needs to write to terraform -> required_providers.
	//
	// - Honestly just need to refactor everything above this, who knows how much is actually needed.
	backendConfig, err := step.mergedConfig(ctx, c, hasTerraformBlock, hasProviderBlock, helper.TerraformVersion())
	backendConfigWithResource := backendConfig + `
resource "examplecloud_thing" "tf_plugin_testing_resource_foo" {}
`

	if err != nil {
		logging.HelperResourceError(ctx,
			"Error generating merged configuration",
			map[string]interface{}{logging.KeyError: err},
		)
		t.Fatalf("Error generating merged configuration: %s", err)
	}

	confRequestA := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       backendConfigWithResource, // gets the resource + backend config
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	testStepConfigA := teststep.Configuration(confRequestA)

	confRequestB := teststep.PrepareConfigurationRequest{
		Directory: step.ConfigDirectory,
		File:      step.ConfigFile,
		Raw:       backendConfig, // just the backend config
		TestStepConfigRequest: config.TestStepConfigRequest{
			StepNumber: stepNumber,
			TestName:   t.Name(),
		},
	}.Exec()

	testStepConfigB := teststep.Configuration(confRequestB)

	finishApply := make(chan bool)

	providers.protov6 = providers.protov6.merge(protov6ProviderFactories{
		"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
			Resources: map[string]testprovider.Resource{
				"examplecloud_thing": examplecloudResource(
					func(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
						fmt.Println("waiting for channel to indicate finish apply")
						<-finishApply
						fmt.Println("received indication that apply can finish")
					},
				),
			},
		}),
	})

	// 1. Create two TF working directories (using the same backend config and run terraform init on both
	workingDirA := helper.RequireNewWorkingDir(ctx, t, "")
	workingDirB := helper.RequireNewWorkingDir(ctx, t, "")
	defer workingDirA.Close()
	defer workingDirB.Close()

	err = workingDirA.SetConfig(ctx, testStepConfigA, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, workingDirA, providers, func() error {
		return workingDirA.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	err = workingDirB.SetConfig(ctx, testStepConfigB, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, workingDirB, providers, func() error {
		return workingDirB.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	go func() {
		// 2. Run terraform apply on working directory 1 (let it pause with a custom resource)
		err = runProviderCommand(ctx, t, workingDirA, providers, func() error {
			return workingDirA.Apply(ctx)
		})
		// if err != nil {
		// 	// TODO: idk, panic? lol
		// }
	}()

	// 3. Get state of working directory 2, should be no error and empty
	var defaultState *tfjson.State
	err = runProviderCommand(ctx, t, workingDirB, providers, func() error {
		defaultState, err = workingDirB.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"default\" state: %w", err)
	}

	if defaultState.Values != nil && defaultState.Values.RootModule != nil && len(defaultState.Values.RootModule.Resources) > 0 {
		t.Fatalf("Expected the newly created \"default\" state to be empty. Found %d resources.", len(defaultState.Values.RootModule.Resources))
	}

	// 4. Run terraform apply on working directory 2, assert that an error occurs, i.e. it's locked still
	err = runProviderCommand(ctx, t, workingDirA, providers, func() error {
		return workingDirA.Apply(ctx)
	})

	// todo: probably use wait group to ensure all applies are guaranteed to finish (this is kind of timing based)
	// 5. Send message to custom resource to let working directory 1 finish successfully it's apply from step #2
	finishApply <- true

	var lockID string
	if err != nil {
		lockErr := regexp.MustCompile(`Error acquiring the state lock`)
		if !lockErr.MatchString(err.Error()) {
			return fmt.Errorf("Expected lock error when running, received different error: %w", err)
		}

		// TODO: This is probably good enough for now, but this really should be a structured CLI output
		// since it will be up to the provider to produce this error message. Framework should probably create some form
		// of lock info on behalf of the provider, then produce error messaging that is uniform, ideally, keeping the exact
		// same format as below.
		//
		//    Error: Error acquiring the state lock
		//
		//     Error message: operation error S3: PutObject, https response error
		//     StatusCode: 412, RequestID: TRNCB4JS3Z46S65Z, HostID:
		//     okvqgtxa8+jVOz/WvleaTbhsNIQccxRD0kYDlnGIKiJsV4BRFadHx3RNdGQR05P4b7Rq6umUKj0=,
		//     api error PreconditionFailed: At least one of the pre-conditions you
		//     specified did not hold
		//     Lock Info:
		//       ID:        21992728-31bd-1d8c-648e-815b75e861a6
		//       Path:      test-pss-backend/av-terraform-state
		//       Operation: OperationTypeApply
		//       Who:       austin.valle@austin.valle-K6YK19LNPP
		//       Version:   1.12.1
		//       Created:   2025-06-05 15:02:29.816176 +0000 UTC
		//       Info:
		//
		//     Terraform acquires a state lock to protect the state from being written
		//     by multiple users at the same time. Please resolve the issue above and try
		//     again. For most commands, you can disable locking with the "-lock=false"
		//     flag, but this is not recommended.
		//
		lockIDExtract := regexp.MustCompile(`\sID:\s+(\S+)`)
		matches := lockIDExtract.FindStringSubmatch(err.Error())

		if len(matches) != 2 {
			return fmt.Errorf("Failed to extract lock ID from error message, matches in error: %v", matches)
		}

		lockID = matches[1]
	} else {
		return errors.New("Expected error when attempting to apply to \"default\" state that should be locked")
	}

	fmt.Println(lockID) // we don't need this here, but just keeping it for now

	// 6. Get state and ensure the resource is populated from step #2
	var populatedState *tfjson.State
	err = runProviderCommand(ctx, t, workingDirA, providers, func() error {
		populatedState, err = workingDirA.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"default\" state: %w", err)
	}

	// 6a. Check if the resource exists in the "default" state (I'm lazy lol)
	checkOutput := statecheck.ExpectKnownValue(
		"examplecloud_thing.tf_plugin_testing_resource_foo",
		tfjsonpath.New("id"),
		knownvalue.StringExact("id-123"),
	)

	checkResp := statecheck.CheckStateResponse{}
	checkOutput.CheckState(ctx, statecheck.CheckStateRequest{State: populatedState}, &checkResp)

	if checkResp.Error != nil {
		return fmt.Errorf("After writing a test resource instance object to \"default\" and re-reading it, the object has vanished: %w", err)
	}

	// 7. Run terraform apply on working directory 2, no pause
	err = runProviderCommand(ctx, t, workingDirB, providers, func() error {
		return workingDirB.Apply(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error removing fake resource in \"default\" workspace: %w", err)
	}

	// 8. Get state and ensure the resources are deleted from step #7
	var emptyState *tfjson.State
	err = runProviderCommand(ctx, t, workingDirB, providers, func() error {
		emptyState, err = workingDirB.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"default\" state: %w", err)
	}

	if emptyState.Values != nil && emptyState.Values.RootModule != nil && len(emptyState.Values.RootModule.Resources) > 0 {
		t.Fatalf("Expected the \"default\" state to be empty. Found %d resources.", len(emptyState.Values.RootModule.Resources))
	}

	// TODO: implement
	if step.ForceUnlockTest {
		// -------------------------------
		// --- If testing force unlock ---
		// -------------------------------
		// 9. Run terraform apply on working directory 1 (let it pause with a custom resource)
		// 10. Run terraform apply on working directory 2, assert that an error occurs, i.e. it's locked still
		// 11. Run force unlock with lock ID from step #10 error message
		// 12. Run terraform apply on working directory 2, no pause, assert no error
		// 13. Send message to custom resource to let working directory 1 finish successfully it's apply from step #2
	}

	return nil
}

func examplecloudResource(createFunc func(context.Context, resource.CreateRequest, *resource.CreateResponse)) testprovider.Resource {
	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "id-123"),
				},
			),
		},
		CreateFunc: createFunc,
		ReadResponse: &resource.ReadResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "id-123"),
				},
			),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						{
							Name:     "id",
							Type:     tftypes.String,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
