// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

// Semantic coverage evidence of one test step is recorded on a
// *semanticcoverage.Step, which step functions get as an explicit coverage
// parameter.
//
// The step functions in helper/resource choose when to record evidence, and
// package semanticcoverage extracts it. Step code records inline, right after
// the WorkingDir operation that produced the evidence: it sets the
// configuration and then observes it, or reads a plan or state and then
// observes it. This keeps coverage code out of plugintest.
//
// Record only the step's own operations, not setup, cleanup, or configuration
// that the harness generates.

import (
	"context"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

// newCoverageCollector returns the coverage collector for the test t, or nil
// if collection is disabled.
func newCoverageCollector(ctx context.Context, t testing.T) *semanticcoverage.Collector {
	t.Helper()

	dir := os.Getenv(EnvTfAccSemanticCoverageDir)
	c, err := semanticcoverage.New(dir, t.Name())
	if err != nil {
		t.Fatalf("Error creating semantic coverage collector: %s", err)
	}
	if c != nil {
		logging.HelperResourceDebug(ctx, "Semantic coverage collection enabled", map[string]interface{}{logging.KeySemanticCoverageDirectory: dir})
	}
	return c
}

// recordPostTestDestroy writes the finding of the destroy that runs after the
// last step, unless the test has failed. number is the number after the last
// step's, served the addresses of the providers served from the test binary,
// and state the state before the destroy. Call it only after the destroy and
// TestCase.CheckDestroy succeeded.
func recordPostTestDestroy(ctx context.Context, t testing.T, c *semanticcoverage.Collector, number int, served []string, state *tfjson.State) {
	t.Helper()

	coverage := c.BeginPostTestDestroy(number)
	coverage.ObserveServedProviders(served)
	coverage.DestroyCompleted(state)
	coverage.End(logging.TestStepNumberContext(ctx, number), t)
}
