// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"slices"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestPlanOnlyDoesNotCompleteApply(t *testing.T) {
	t.Parallel()

	c := &Collector{}
	c.BeginStep(1)
	c.ObservePlan(&tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}})
	f := c.feature(KindResource, "example_item", "example_item.test")
	if len(f.planned) != 1 || len(f.completed) != 0 || len(f.known) != 0 {
		t.Fatalf("planned=%v completed=%v known=%v", f.planned, f.completed, f.known)
	}
}

func TestPreApplyPlanWithoutReadablePlanWarns(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), testName: "TestUnavailable"}
	c.BeginStep(1)
	c.PreApplyPlan(nil)
	if !slices.Equal(c.warnings.sorted(), []Warning{WarningPlanUnavailable}) {
		t.Fatalf("unreadable pre-apply plan, warnings: %v", c.warnings)
	}
	c.BeginStep(2)
	c.ObservePlan(nil)
	c.PreApplyPlan(nil)
	c.ApplyCompleted(false)
	if !slices.Equal(c.warnings.sorted(), []Warning{WarningPlanUnavailable}) {
		t.Fatalf("warnings: %v", c.warnings)
	}
	for _, f := range c.features {
		if len(f.completed) != 0 {
			t.Fatalf("unreadable plan claimed completion: %v", f)
		}
	}
}
