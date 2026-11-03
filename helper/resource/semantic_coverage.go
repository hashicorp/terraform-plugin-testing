// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

// stepCoverage records the semantic coverage evidence of one test step. Step
// functions get it as an explicit parameter.
//
// helper/resource decides when to record, semanticcoverage decides what. The
// wrappers below are for operations whose evidence needs an extra read, so
// that plugintest doesn't need to know about coverage. Documents that the
// harness reads anyway, such as plans for plan checks, are passed to the
// embedded Step directly.
//
// Only record the step's own operations, not setup, cleanup, or configuration
// that the harness generates. With collection disabled (a nil Step) every
// method only performs the wrapped operation.
type stepCoverage struct {
	*semanticcoverage.Step
}

// newCoverageCollector returns the coverage collector for the test t, or nil
// if collection is disabled.
func newCoverageCollector(t testing.T) *semanticcoverage.Collector {
	t.Helper()

	c, err := semanticcoverage.New(t.Name())
	if err != nil {
		t.Fatalf("semantic coverage: %s", err)
	}
	return c
}

// end writes the step's finding, if a step is being recorded, and stops
// recording.
func (s *stepCoverage) end(t testing.T, failed bool) {
	t.Helper()

	if err := s.Finish(failed); err != nil {
		t.Errorf("writing semantic coverage finding: %s", err)
	}
	s.Step = nil
}

// setConfig sets the configuration of wd and records it.
func (s stepCoverage) setConfig(ctx context.Context, wd *plugintest.WorkingDir, cfg teststep.Config, vars config.Variables) error {
	if err := wd.SetConfig(ctx, cfg, vars); err != nil {
		return err
	}
	s.ObserveConfigDir(wd.BaseDir())
	return nil
}

// setQuery sets the query configuration of wd and records it.
func (s stepCoverage) setQuery(ctx context.Context, wd *plugintest.WorkingDir, cfg teststep.Config, vars config.Variables) error {
	if err := wd.SetQuery(ctx, cfg, vars); err != nil {
		return err
	}
	s.ObserveQueryConfig(filepath.Join(wd.BaseDir(), plugintest.QueryFileName))
	return nil
}

// apply applies the saved plan in wd. It records the plan's operations as
// planned and, once the apply succeeded, as completed, and records the
// resulting state. Call it inside runProviderCommand, because reading the plan
// and the state needs the providers.
//
// Without the state read, a step with no Check would record no populated
// attributes at all.
func (s stepCoverage) apply(ctx context.Context, wd *plugintest.WorkingDir, opts ...tfexec.ApplyOption) error {
	if !s.Enabled() {
		return wd.Apply(ctx, opts...)
	}
	// A nil plan or state records that it could not be read. terraform-exec
	// returns nil with every error.
	plan, _ := wd.SavedPlan(ctx)
	s.ObservePlan(plan)
	if err := wd.Apply(ctx, opts...); err != nil {
		return err
	}
	s.ApplyCompleted(plan)
	state, _ := wd.State(ctx)
	s.ObserveState(state)
	return nil
}
