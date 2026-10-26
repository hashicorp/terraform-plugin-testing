// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestCompletedLifecycleUsesAppliedPlan(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), testName: "TestLifecycle"}
	c.BeginStep(1)
	create := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.ObservePlan(create)
	c.PreApplyPlan(create)
	c.ObservePlan(noOp)
	c.ApplyCompleted(false)
	f := c.feature(KindResource, "example_item", "example_item.test")
	if !contains(f.completed, "create") || !contains(f.completed, "apply") || contains(f.completed, "no_op") {
		t.Fatalf("completed actions: %v", f.completed)
	}
	c.BeginStep(2)
	c.ObservePlan(noOp)
	c.PreApplyPlan(noOp)
	c.ApplyCompleted(true)
	f = c.feature(KindResource, "example_item", "example_item.test")
	if len(f.completed) != 0 {
		t.Fatalf("no-op destroy step claimed completion: %v", f.completed)
	}
}

func TestNoOpApplyIsStillAnApply(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), testName: "TestNoOp"}
	c.BeginStep(1)
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.ObservePlan(noOp)
	c.PreApplyPlan(noOp)
	c.ApplyCompleted(false)
	f := c.feature(KindResource, "example_item", "example_item.test")
	if !contains(f.completed, "apply") || !contains(f.completed, "no_op") {
		t.Fatalf("completed actions: %v", f.completed)
	}
}
