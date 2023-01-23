// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

func testStepNewRefreshState(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories, stdout *plugintest.TerraformJSONBuffer) error {
	t.Helper()

	var err error
	// Explicitly ensure prior state exists before refresh.
	err = runProviderCommand(ctx, t, func() error {
		_, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	err = runProviderCommand(ctx, t, func() error {
		return wd.RefreshJSON(ctx, stdout)
	}, wd, providers)
	if err != nil {
		target := &tfexec.ErrVersionMismatch{}
		if errors.As(err, &target) {
			err = runProviderCommand(ctx, t, func() error {
				return wd.Refresh(ctx)
			}, wd, providers)
			if err != nil {
				return fmt.Errorf("Error running refresh: %w", err)
			}
		} else {
			return fmt.Errorf("Error running refresh: %s", stdout.GetJSONOutputStr())
		}
	}

	var refreshState *terraform.State
	err = runProviderCommand(ctx, t, func() error {
		refreshState, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providers)
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// Go through the refreshed state and verify
	if step.Check != nil {
		logging.HelperResourceDebug(ctx, "Calling TestStep Check for RefreshState")

		if err := step.Check(refreshState); err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep Check for RefreshState")
	}

	// do a plan
	err = runProviderCommand(ctx, t, func() error {
		return wd.CreatePlanJSON(ctx, stdout)
	}, wd, providers)
	if err != nil {
		target := &tfexec.ErrVersionMismatch{}
		if errors.As(err, &target) {
			err = runProviderCommand(ctx, t, func() error {
				return wd.CreatePlan(ctx)
			}, wd, providers)
			if err != nil {
				return fmt.Errorf("Error running post-apply plan: %w", err)
			}
		} else {
			return fmt.Errorf("Error running post-apply plan: %s", stdout.GetJSONOutputStr())
		}
	}

	var plan *tfjson.Plan
	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error retrieving post-apply plan: %w", err)
	}

	if !planIsEmpty(plan) && !step.ExpectNonEmptyPlan {
		var stdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			stdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error retrieving formatted plan output: %w", err)
		}
		return fmt.Errorf("After refreshing state during this test step, a followup plan was not empty.\nstdout:\n\n%s", stdout)
	}

	return nil
}
