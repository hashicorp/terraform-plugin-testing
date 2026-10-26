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
	c.Plan(&tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}})
	f := c.feature("resource", "example_item", "example_item.test")
	if len(f.Planned) != 1 || len(f.Completed) != 0 || len(f.Known) != 0 {
		t.Fatalf("planned=%v completed=%v known=%v", f.Planned, f.Completed, f.Known)
	}
}

func TestPreApplyPlanWithoutReadablePlanWarns(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestUnavailable"}
	c.BeginStep(1)
	c.PreApplyPlan(nil)
	if !slices.Equal(c.warnings, []Warning{WarningPlanUnavailable}) {
		t.Fatalf("unreadable pre-apply plan, warnings: %v", c.warnings)
	}
	c.BeginStep(2)
	c.Plan(nil)
	c.PreApplyPlan(nil)
	c.ApplyCompleted(false)
	if !slices.Equal(c.warnings, []Warning{WarningPlanUnavailable}) {
		t.Fatalf("warnings: %v", c.warnings)
	}
	for _, f := range c.features {
		if len(f.Completed) != 0 {
			t.Fatalf("unreadable plan claimed completion: %v", f)
		}
	}
}
