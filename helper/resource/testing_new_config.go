// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

func testStepNewConfig(ctx context.Context, t testing.T, c TestCase, wd *plugintest.WorkingDir, step TestStep, providers *providerFactories) error {
	t.Helper()

	err := wd.SetConfig(ctx, step.mergedConfig(ctx, c))
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	// require a refresh before applying
	// failing to do this will result in data sources not being updated
	err = runProviderCommand(ctx, t, func() error {
		return wd.Refresh(ctx)
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error running pre-apply refresh: %w", err)
	}

	// If this step is a PlanOnly step, skip over this first Plan and
	// subsequent Apply, and use the follow-up Plan that checks for
	// permadiffs
	if !step.PlanOnly {
		logging.HelperResourceDebug(ctx, "Running Terraform CLI plan and apply")

		// Plan!
		err := runProviderCommand(ctx, t, func() error {
			if step.Destroy {
				return wd.CreateDestroyPlan(ctx)
			}
			return wd.CreatePlan(ctx)
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error running pre-apply plan: %w", err)
		}

		// If there are expectations for the planned diffs of resources, we assert those with the saved plan file
		if len(step.ExpectedResourceChanges) > 0 {
			var plan *tfjson.Plan
			err = runProviderCommand(ctx, t, func() error {
				var err error
				plan, err = wd.SavedPlan(ctx)
				return err
			}, wd, providers)
			if err != nil {
				return fmt.Errorf("Error retrieving pre-apply plan: %w", err)
			}

			err = assertExpectedResourceChanges(step, plan)
			if err != nil {
				return fmt.Errorf("Error asserting ExpectedResourceChanges: %w", err)
			}
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
			return fmt.Errorf("Error retrieving pre-apply state: %w", err)
		}

		// Apply the diff, creating real resources
		err = runProviderCommand(ctx, t, func() error {
			return wd.Apply(ctx)
		}, wd, providers)
		if err != nil {
			if step.Destroy {
				return fmt.Errorf("Error running destroy: %w", err)
			}
			return fmt.Errorf("Error running apply: %w", err)
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
			return fmt.Errorf("Error retrieving state after apply: %w", err)
		}

		// Run any configured checks
		if step.Check != nil {
			logging.HelperResourceTrace(ctx, "Using TestStep Check")

			state.IsBinaryDrivenTest = true
			if step.Destroy {
				if err := step.Check(stateBeforeApplication); err != nil {
					return fmt.Errorf("Check failed: %w", err)
				}
			} else {
				if err := step.Check(state); err != nil {
					return fmt.Errorf("Check failed: %w", err)
				}
			}
		}
	}

	// Test for perpetual diffs by performing a plan, a refresh, and another plan
	logging.HelperResourceDebug(ctx, "Running Terraform CLI plan to check for perpetual differences")

	// do a plan
	err = runProviderCommand(ctx, t, func() error {
		if step.Destroy {
			return wd.CreateDestroyPlan(ctx)
		}
		return wd.CreatePlan(ctx)
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error running post-apply plan: %w", err)
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
		return fmt.Errorf("After applying this test step, the plan was not empty.\nstdout:\n\n%s", stdout)
	}

	// do a refresh
	if !step.Destroy || (step.Destroy && !step.PreventPostDestroyRefresh) {
		err := runProviderCommand(ctx, t, func() error {
			return wd.Refresh(ctx)
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error running post-apply refresh: %w", err)
		}
	}

	// do another plan
	err = runProviderCommand(ctx, t, func() error {
		if step.Destroy {
			return wd.CreateDestroyPlan(ctx)
		}
		return wd.CreatePlan(ctx)
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error running second post-apply plan: %w", err)
	}

	err = runProviderCommand(ctx, t, func() error {
		var err error
		plan, err = wd.SavedPlan(ctx)
		return err
	}, wd, providers)
	if err != nil {
		return fmt.Errorf("Error retrieving second post-apply plan: %w", err)
	}

	// check if plan is empty
	if !planIsEmpty(plan) && !step.ExpectNonEmptyPlan {
		var stdout string
		err = runProviderCommand(ctx, t, func() error {
			var err error
			stdout, err = wd.SavedPlanRawStdout(ctx)
			return err
		}, wd, providers)
		if err != nil {
			return fmt.Errorf("Error retrieving formatted second plan output: %w", err)
		}
		return fmt.Errorf("After applying this test step and performing a `terraform refresh`, the plan was not empty.\nstdout\n\n%s", stdout)
	} else if step.ExpectNonEmptyPlan && planIsEmpty(plan) {
		return errors.New("Expected a non-empty plan, but got an empty plan")
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
			return err
		}

		if state.Empty() {
			return nil
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
				return fmt.Errorf(
					"[ERROR] Test: ID-only test failed: %s", err)
			}
		}
	}

	return nil
}

// TODO: refactor this logic?
func assertExpectedResourceChanges(step TestStep, plan *tfjson.Plan) error {
	var result *multierror.Error
	for resource, expectedChange := range step.ExpectedResourceChanges {
		foundResource := false

		for _, rc := range plan.ResourceChanges {
			if resource == rc.Address {
				switch expectedChange {
				case DiffNoop:
					if !rc.Change.Actions.NoOp() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected NoOp, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffCreate:
					if !rc.Change.Actions.Create() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected Create, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffRead:
					if !rc.Change.Actions.Read() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected Read, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffUpdate:
					if !rc.Change.Actions.Update() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected Update, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffDestroy:
					if !rc.Change.Actions.Delete() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected Destroy, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffDestroyBeforeCreate:
					if !rc.Change.Actions.DestroyBeforeCreate() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected DestroyBeforeCreate, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffCreateBeforeDestroy:
					if !rc.Change.Actions.CreateBeforeDestroy() {
						result = multierror.Append(result, fmt.Errorf("'%s' - expected CreateBeforeDestroy, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				case DiffReplace:
					if !rc.Change.Actions.Replace() {
						result = multierror.Append(result, fmt.Errorf("%s - expected Replace, got action(s): %v", rc.Address, rc.Change.Actions))
					}
				default:
					result = multierror.Append(result, fmt.Errorf("%s - unexpected DiffChangeType byte: %d", rc.Address, expectedChange))
				}

				foundResource = true
				break
			}
		}

		if !foundResource {
			result = multierror.Append(result, fmt.Errorf("%s - Resource not found in planned ResourceChanges", resource))
		}
	}

	return result.ErrorOrNil()
}
