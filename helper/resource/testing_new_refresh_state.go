// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

type testStepNewRefreshStateResponse struct {
	tfJSONDiags plugintest.TerraformJSONDiagnostics
	stdout      string
}

func testStepNewRefreshState(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories) (testStepNewRefreshStateResponse, error) {
	t.Helper()

	var tfJSONDiags plugintest.TerraformJSONDiagnostics

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
		refreshResponse, err := wd.Refresh(ctx)

		tfJSONDiags = append(tfJSONDiags, refreshResponse.Diagnostics...)

		return err
	}, wd, providers)
	if err != nil {
		return testStepNewRefreshStateResponse{
			tfJSONDiags: tfJSONDiags,
		}, fmt.Errorf("Error running refresh: %w", err)
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
		createPlanResponse, err := wd.CreatePlan(ctx)

		tfJSONDiags = append(tfJSONDiags, createPlanResponse.Diagnostics...)

		return err
	}, wd, providers)
	if err != nil {
		return testStepNewRefreshStateResponse{
			tfJSONDiags: tfJSONDiags,
		}, fmt.Errorf("Error running post-apply plan: %w", err)
	}

	var plan *tfjson.Plan
	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		return testStepNewRefreshStateResponse{
			tfJSONDiags: tfJSONDiags,
		}, fmt.Errorf("Error retrieving post-apply plan: %w", err)
	}

	if !planIsEmpty(plan) && !step.ExpectNonEmptyPlan {
		var stdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			stdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return testStepNewRefreshStateResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error retrieving formatted plan output: %w", err)
		}
		return testStepNewRefreshStateResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("After refreshing state during this test step, a followup plan was not empty.\nstdout:\n\n%s", stdout)
	}

	return testStepNewRefreshStateResponse{
		tfJSONDiags: tfJSONDiags,
	}, nil
}
