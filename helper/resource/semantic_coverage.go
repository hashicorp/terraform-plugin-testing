// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

// stepCoverage records semantic coverage evidence for the steps of one
// TestCase. It wraps the WorkingDir operations whose results are evidence, so
// that plugintest doesn't need to know about coverage. Step functions get it
// as an explicit parameter.
//
// Only use the wrappers on the step's own operations, not in setup or cleanup.
// With a nil Collector (collection disabled) every method only performs the
// wrapped operation.
type stepCoverage struct {
	*semanticcoverage.Collector
}

// newStepCoverage returns the coverage recorder for the test t.
func newStepCoverage(t testing.T) stepCoverage {
	t.Helper()

	c := semanticcoverage.New(t.Name())
	if err := c.Ready(); err != nil {
		t.Fatalf("semantic coverage: %s", err)
	}
	return stepCoverage{c}
}

// beginStep writes the previous step, if any, and starts collecting for
// stepNumber. Reaching the next step means the previous step passed its
// checks.
func (s stepCoverage) beginStep(t testing.T, stepNumber int) {
	t.Helper()

	if err := s.FinishStep(semanticcoverage.OutcomeSuccess); err != nil {
		t.Errorf("writing semantic coverage finding: %s", err)
	}
	s.BeginStep(stepNumber)
}

// finishTest writes the last step, or the step that ended in t.Fatal. A step
// that ended in t.Skip is discarded. Defer it at the start of the test.
func (s stepCoverage) finishTest(t testing.T) {
	t.Helper()

	outcome := semanticcoverage.OutcomeSuccess
	switch {
	case t.Failed():
		outcome = semanticcoverage.OutcomeFailed
	case t.Skipped():
		s.DiscardStep()
		return
	}
	if err := s.FinishStep(outcome); err != nil {
		t.Errorf("writing semantic coverage finding: %s", err)
	}
}

// stepResult records how a step ended. err is the error that the step
// returned, before TestCase.ErrorCheck. completed records the step's
// completion evidence; it runs only when the step succeeded without
// ExpectError, and may be nil.
func (s stepCoverage) stepResult(step TestStep, err error, completed func()) {
	if step.ExpectError != nil {
		if err != nil && step.ExpectError.MatchString(err.Error()) {
			s.MatchedExpectedError()
		}
		return
	}
	if err == nil && completed != nil {
		completed()
	}
}

// setConfig sets the configuration of wd and records the configuration files
// as input evidence.
func (s stepCoverage) setConfig(ctx context.Context, wd *plugintest.WorkingDir, cfg teststep.Config, vars config.Variables) error {
	if err := wd.SetConfig(ctx, cfg, vars); err != nil {
		return err
	}
	s.Config(wd.BaseDir())
	return nil
}

// setQuery sets the query configuration of wd and records it as input
// evidence.
func (s stepCoverage) setQuery(ctx context.Context, wd *plugintest.WorkingDir, cfg teststep.Config, vars config.Variables) error {
	if err := wd.SetQuery(ctx, cfg, vars); err != nil {
		return err
	}
	s.QueryConfig(wd.BaseDir())
	return nil
}

// createPlan creates a saved plan in wd and records its planned operations.
// Call it inside runProviderCommand, because reading the saved plan needs the
// providers.
func (s stepCoverage) createPlan(ctx context.Context, wd *plugintest.WorkingDir, opts ...tfexec.PlanOption) error {
	_, err := s.plan(ctx, wd, opts...)
	return err
}

// createPreApplyPlan is createPlan for the plan that the next apply executes.
// The completion evidence of that apply comes from this plan.
func (s stepCoverage) createPreApplyPlan(ctx context.Context, wd *plugintest.WorkingDir, opts ...tfexec.PlanOption) error {
	plan, err := s.plan(ctx, wd, opts...)
	if err != nil {
		return err
	}
	s.PreApplyPlan(plan)
	return nil
}

func (s stepCoverage) plan(ctx context.Context, wd *plugintest.WorkingDir, opts ...tfexec.PlanOption) (*tfjson.Plan, error) {
	if err := wd.CreatePlan(ctx, opts...); err != nil {
		return nil, err
	}
	if !s.Enabled() {
		return nil, nil
	}
	// A nil plan records that the saved plan could not be read.
	plan, _ := wd.SavedPlan(ctx)
	s.Plan(plan)
	return plan, nil
}

// observeAppliedState records the state that a successful apply produced.
//
// Known and null attribute evidence must come from state that an apply
// actually produced; plan values are not used for it (see Collector.Plan).
// Without this read, a step with no Check would record no state evidence at
// all.
func (s stepCoverage) observeAppliedState(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, providers *providerFactories) {
	t.Helper()

	if !s.Enabled() {
		return
	}
	var state *tfjson.State
	err := runProviderCommand(ctx, t, wd, providers, func() error {
		var err error
		state, err = wd.State(ctx)
		return err
	})
	if err != nil {
		s.Warning(semanticcoverage.WarningStateUnavailable)
		return
	}
	s.ObserveState(state)
}
