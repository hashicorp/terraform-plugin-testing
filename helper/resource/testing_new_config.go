// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

type testStepNewConfigResponse struct {
	tfJSONDiags plugintest.TerraformJSONDiagnostics
	stdout      string
}

func testStepNewConfig(ctx context.Context, t testing.T, c TestCase, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories) (testStepNewConfigResponse, error) {
	t.Helper()

	var tfJSONDiags plugintest.TerraformJSONDiagnostics
	var stdout string

	err := wd.SetConfig(ctx, step.mergedConfig(ctx, c))
	if err != nil {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("Error setting config: %w", err)
	}

	// require a refresh before applying
	// failing to do this will result in data sources not being updated
	err = runProviderCommand(ctx, t, func() error {
		refreshResponse, err := wd.Refresh(ctx)

		tfJSONDiags = append(tfJSONDiags, refreshResponse.Diagnostics...)
		stdout += refreshResponse.Stdout

		return err
	}, wd, providers)
	if err != nil {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("Error running pre-apply refresh: %w", err)
	}

	// If this step is a PlanOnly step, skip over this first Plan and
	// subsequent Apply, and use the follow-up Plan that checks for
	// permadiffs
	if !step.PlanOnly {
		logging.HelperResourceDebug(ctx, "Running Terraform CLI plan and apply")

		// Plan!
		err := runProviderCommand(ctx, t, func() error {
			if step.Destroy {
				createDestroyPlanResponse, err := wd.CreateDestroyPlan(ctx)

				tfJSONDiags = append(tfJSONDiags, createDestroyPlanResponse.Diagnostics...)
				stdout += createDestroyPlanResponse.Stdout

				return err
			}
			createPlanResponse, err := wd.CreatePlan(ctx)

			tfJSONDiags = append(tfJSONDiags, createPlanResponse.Diagnostics...)
			stdout += createPlanResponse.Stdout

			return err
		}, wd, providers)
		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error running pre-apply plan: %w", err)
		}

		// We need to keep a copy of the state prior to destroying such
		// that the destroy steps can verify their behavior in the
		// check function
		var stateBeforeApplication *terraform.State
		err = runProviderCommand(ctx, t, func() error {
			stateBeforeApplication, err = getState(ctx, t, wd)
			if err != nil {
				return err
			}
			return nil
		}, wd, providers)
		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error retrieving pre-apply state: %w", err)
		}

		// Apply the diff, creating real resources
		err = runProviderCommand(ctx, t, func() error {
			applyResponse, err := wd.Apply(ctx)

			tfJSONDiags = append(tfJSONDiags, applyResponse.Diagnostics...)
			stdout += applyResponse.Stdout

			return err
		}, wd, providers)
		if err != nil {
			if step.Destroy {
				return testStepNewConfigResponse{
					tfJSONDiags: tfJSONDiags,
					stdout:      stdout,
				}, fmt.Errorf("Error running destroy: %w", err)
			}

			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error running apply: %w", err)
		}

		// Get the new state
		var state *terraform.State
		err = runProviderCommand(ctx, t, func() error {
			state, err = getState(ctx, t, wd)
			if err != nil {
				return err
			}
			return nil
		}, wd, providers)
		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error retrieving state after apply: %w", err)
		}

		// Run any configured checks
		if step.Check != nil {
			logging.HelperResourceTrace(ctx, "Using TestStep Check")

			state.IsBinaryDrivenTest = true
			if step.Destroy {
				if err := step.Check(stateBeforeApplication); err != nil {
					return testStepNewConfigResponse{
						tfJSONDiags: tfJSONDiags,
						stdout:      stdout,
					}, fmt.Errorf("Check failed: %w", err)
				}
			} else {
				if err := step.Check(state); err != nil {
					return testStepNewConfigResponse{
						tfJSONDiags: tfJSONDiags,
						stdout:      stdout,
					}, fmt.Errorf("Check failed: %w", err)
				}
			}
		}
	}

	// Test for perpetual diffs by performing a plan, a refresh, and another plan
	logging.HelperResourceDebug(ctx, "Running Terraform CLI plan to check for perpetual differences")

	// do a plan
	err = runProviderCommand(ctx, t, func() error {
		if step.Destroy {
			createDestroyPlanResponse, err := wd.CreateDestroyPlan(ctx)

			tfJSONDiags = append(tfJSONDiags, createDestroyPlanResponse.Diagnostics...)
			stdout += createDestroyPlanResponse.Stdout

			return err
		}
		createPlanResponse, err := wd.CreatePlan(ctx)

		tfJSONDiags = append(tfJSONDiags, createPlanResponse.Diagnostics...)
		stdout += createPlanResponse.Stdout

		return err
	}, wd, providers)
	if err != nil {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("Error running post-apply plan: %w", err)
	}

	var plan *tfjson.Plan
	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("Error retrieving post-apply plan: %w", err)
	}

	if !planIsEmpty(plan) && !step.ExpectNonEmptyPlan {
		var savedPlanRawStdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			savedPlanRawStdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error retrieving formatted plan output: %w", err)
		}
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("After applying this test step, the plan was not empty.\nstdout:\n\n%s", savedPlanRawStdout)
	}

	// do a refresh
	if !step.Destroy || (step.Destroy && !step.PreventPostDestroyRefresh) {
		err = runProviderCommand(ctx, t, func() error {
			refreshResponse, err := wd.Refresh(ctx)

			tfJSONDiags = append(tfJSONDiags, refreshResponse.Diagnostics...)
			stdout += refreshResponse.Stdout

			return err
		}, wd, providers)
		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error running post-apply refresh: %w", err)
		}
	}

	// do another plan
	err = runProviderCommand(ctx, t, func() error {
		if step.Destroy {
			createDestroyPlanResponse, err := wd.CreateDestroyPlan(ctx)

			tfJSONDiags = append(tfJSONDiags, createDestroyPlanResponse.Diagnostics...)
			stdout += createDestroyPlanResponse.Stdout

			return err
		}
		createPlanResponse, err := wd.CreatePlan(ctx)

		tfJSONDiags = append(tfJSONDiags, createPlanResponse.Diagnostics...)
		stdout += createPlanResponse.Stdout

		return err
	}, wd, providers)
	if err != nil {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("Error running second post-apply plan: %w", err)
	}

	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("Error retrieving second post-apply plan: %w", err)
	}

	// check if plan is empty
	if !planIsEmpty(plan) && !step.ExpectNonEmptyPlan {
		var savedPlanRawStdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			savedPlanRawStdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, fmt.Errorf("Error retrieving formatted second plan output: %w", err)
		}
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, fmt.Errorf("After applying this test step and performing a `terraform refresh`, the plan was not empty.\nstdout\n\n%s", savedPlanRawStdout)
	} else if step.ExpectNonEmptyPlan && planIsEmpty(plan) {
		return testStepNewConfigResponse{
			tfJSONDiags: tfJSONDiags,
			stdout:      stdout,
		}, errors.New("Expected a non-empty plan, but got an empty plan")
	}

	// ID-ONLY REFRESH
	// If we've never checked an id-only refresh and our state isn't
	// empty, find the first resource and test it.
	if c.IDRefreshName != "" {
		logging.HelperResourceTrace(ctx, "Using TestCase IDRefreshName")

		var state *terraform.State

		err = runProviderCommand(ctx, t, func() error {
			state, err = getState(ctx, t, wd)
			if err != nil {
				return err
			}
			return nil
		}, wd, providers)

		if err != nil {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, err
		}

		if state.Empty() {
			return testStepNewConfigResponse{
				tfJSONDiags: tfJSONDiags,
				stdout:      stdout,
			}, nil
		}

		var idRefreshCheck *terraform.ResourceState

		// Find the first non-nil resource in the state
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[c.IDRefreshName]; ok {
					idRefreshCheck = v
				}

				break
			}
		}

		// If we have an instance to check for refreshes, do it
		// immediately. We do it in the middle of another test
		// because it shouldn't affect the overall state (refresh
		// is read-only semantically) and we want to fail early if
		// this fails. If refresh isn't read-only, then this will have
		// caught a different bug.
		if idRefreshCheck != nil {
			if err := testIDRefresh(ctx, t, c, wd, step, idRefreshCheck, providers); err != nil {
				return testStepNewConfigResponse{
					tfJSONDiags: tfJSONDiags,
					stdout:      stdout,
				}, fmt.Errorf("[ERROR] Test: ID-only test failed: %s", err)
			}
		}
	}

	return testStepNewConfigResponse{
		tfJSONDiags: tfJSONDiags,
		stdout:      stdout,
	}, nil
}
