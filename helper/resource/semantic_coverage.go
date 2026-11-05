// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

// stepCoverage records the semantic coverage evidence of one test step. Step
// functions get it as an explicit parameter.
//
// The step functions in helper/resource choose when to record evidence, and
// package semanticcoverage extracts it. Step code records inline, right after
// the WorkingDir operation that produced the evidence: it sets the
// configuration and then observes it, or reads a plan or state and then
// observes it. This keeps coverage code out of plugintest.
//
// Record only the step's own operations, not setup, cleanup, or configuration
// that the harness generates. With collection disabled (a nil Step), every
// method is a no-op; check Enabled before running Terraform commands only to
// feed the step.
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

// end stops recording the step and writes its finding, unless the test has
// failed. A test that reports an error without stopping fails all of its
// remaining steps too.
func (s *stepCoverage) end(t testing.T) {
	t.Helper()

	if !t.Failed() {
		if err := s.Finish(); err != nil {
			t.Errorf("writing semantic coverage finding: %s", err)
		}
	}
	s.Step = nil
}
