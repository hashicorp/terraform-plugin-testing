// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestCompletedLifecycleUsesAppliedPlan(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestLifecycle"}
	c.BeginStep(1)
	create := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.Plan(create)
	c.PreApplyPlan(create)
	c.Plan(noOp)
	c.ApplyCompleted(false)
	f := c.feature("resource", "example_item", "example_item.test")
	if !contains(f.Completed, "create") || !contains(f.Completed, "apply") || contains(f.Completed, "no_op") {
		t.Fatalf("completed actions: %v", f.Completed)
	}
	c.BeginStep(2)
	c.Plan(noOp)
	c.PreApplyPlan(noOp)
	c.ApplyCompleted(true)
	f = c.feature("resource", "example_item", "example_item.test")
	if len(f.Completed) != 0 {
		t.Fatalf("no-op destroy step claimed completion: %v", f.Completed)
	}
}

func TestNoOpApplyIsStillAnApply(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestNoOp"}
	c.BeginStep(1)
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.Plan(noOp)
	c.PreApplyPlan(noOp)
	c.ApplyCompleted(false)
	f := c.feature("resource", "example_item", "example_item.test")
	if !contains(f.Completed, "apply") || !contains(f.Completed, "no_op") {
		t.Fatalf("completed actions: %v", f.Completed)
	}
}
