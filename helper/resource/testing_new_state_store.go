package resource

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mitchellh/go-testing-interface"
)

// TODO:PSS: I should go back through and review this logic, refactor, maybe add/remove any details we don't need
// It's likely we will add lock testing somewhere here as well, so might be worth refactoring at the top-level
//
// 1. StateStore (true) + nothing else == Run basic smoke test
// 2. StateStore (true) + VerifyLock (true) == Run concurrent lock test
// 2. StateStore (true) + VerifyLock (true) + ForceUnlock (true) == Run concurrent lock test, then finish with a force unlock (maybe do this by default?)
// 2. StateStore (true) + ??? == Run lock soak test (configurable)
func testStepNewStateStore(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories, cfg teststep.Config) error {
	t.Helper()

	err := wd.SetConfig(ctx, cfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// ----- Validate and configure the state store
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	// ----- Retrieve all the workspaces
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

	// Assert the only workspace created is the "default" one
	// TODO:PSS: Might need to revisit this assertion for state stores. Not sure what the behavior will be during init, is this expected still?
	// TODO:PSS: Is it possible to not support this for state stores or TF core backends? The cloud backend doesn't but that one is special :P
	if len(workspaces) != 1 || workspaces[0] != "default" {
		t.Fatalf("Expected a single workspace named \"default\" after initialization, got: %#v", workspaces)
	}

	// ----- Create "foo" workspace and assert the state returned is empty
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

	// ----- Create "bar" workspace and assert the state returned is empty
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

	// ----- Add a fake resource to the bar workspace
	barConfig := `
resource "terraform_data" "tf_plugin_testing_resource_bar" {
  input = "this resource was injected by terraform-plugin-testing"
}`

	// TODO:PSS: I'm 99% sure this should work with all of config file/directory implementations
	cfgWithBar := cfg.Append(barConfig)
	err = wd.SetConfig(ctx, cfgWithBar, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// ----- Apply bar workspace
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.Apply(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error creating fake resource in \"bar\" workspace: %w", err)
	}

	// ----- Grab the bar state again
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		barState, err = wd.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"bar\" state: %w", err)
	}

	// ----- Check if the resource exists in the "bar" state
	// TODO:PSS: Should we not use a statecheck here? Could just read the state manually...
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

	// ----- Switch to the "foo" workspace and grab the state, ensuring it's still empty.
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "foo")
	})
	if err != nil {
		return fmt.Errorf("Error creating \"foo\" workspace: %w", err)
	}

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

	// ----- Verify when we list the workspaces we get back "default", "foo" (created during this test), and "bar" (created during this test).
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

	// ----- Delete "bar" workspace

	// Switch to "foo" workspace so we can delete bar.
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

	// ----- Attempt to delete "default" workspace, assert error
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, "default", tfexec.Force(true))
	})
	if err == nil {
		return errors.New("Expected error when deleting \"default\" workspace")
	}

	// ----- Recreate the "bar" workspace
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.CreateWorkspace(ctx, "bar")
	})
	if err != nil {
		return fmt.Errorf("Error creating \"bar\" workspace: %w", err)
	}

	// ----- Grab "bar" state and assert it is empty (i.e. no left over artifacts)
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

	// ----- Delete "bar" workspace again, force=true

	// Switch to "foo" workspace so we can delete bar.
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

	// ----- List workspaces and verify it's just "foo" and "default"
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

	// ----- Delete "foo" workspace, force=true

	// Switch to "default" workspace so we can delete bar.
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.SelectWorkspace(ctx, "default")
	})
	if err != nil {
		return fmt.Errorf("Error selecting \"default\" workspace: %w", err)
	}
	err = runProviderCommand(ctx, t, wd, providers, func() error {
		return wd.DeleteWorkspace(ctx, "foo", tfexec.Force(true))
	})
	if err != nil {
		return fmt.Errorf("Error deleting \"foo\" workspace: %w", err)
	}

	// ----- List workspaces and verify it's just "default" (which we did not modify)
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
	expected = []string{"default"}
	if !slices.Equal(expected, workspaces) {
		t.Fatalf("Expected workspaces to be %#v, got: %#v", expected, workspaces)
	}

	return nil
}
