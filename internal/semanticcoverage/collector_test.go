// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestCollectorKeepsOnlyMetadataAndSeparatesPlansFromApply(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvDirectory, dir)
	c := New("TestExample")
	c.BeginStep(1)
	c.Parse("test.tf", []byte(`resource "example_item" "test" { token = "secret-configuration" }`))
	plan := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: `example_item.test["secret-instance-key"]`, Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	c.Plan(plan)
	c.PreApplyPlan()
	resource := &tfjson.StateResource{Address: `example_item.test["secret-instance-key"]`, Type: "example_item", AttributeValues: map[string]any{"id": "secret-value", "token": "secret-value", "map": map[string]any{"secret-map-key": "secret-value"}}}
	c.ObserveState(&tfjson.State{Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{Resources: []*tfjson.StateResource{resource}}}})
	c.ApplyCompleted(false)
	if err := c.FinishStep(OutcomeSuccess); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d findings", len(entries))
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"secret-configuration", "secret-instance-key", "secret-map-key", "secret-value"} {
		if strings.Contains(string(data), s) {
			t.Errorf("finding persisted test value %q", s)
		}
	}
	for _, s := range []string{`"planned":["create"]`, `"completed":["apply","create"]`, `"known":["id","map","token"]`} {
		if !strings.Contains(string(data), s) {
			t.Errorf("missing %s in %s", s, data)
		}
	}
}

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

func TestStateRequiresExplicitPostOperationObservation(t *testing.T) {
	t.Parallel()

	c := &Collector{}
	c.BeginStep(1)
	c.Plan(&tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}, After: map[string]any{"id": "secret"}}}}})
	f := c.feature("resource", "example_item", "example_item.test")
	if len(f.Known) != 0 {
		t.Fatalf("plan claimed post-operation state: %v", f.Known)
	}
	c.ObserveState(&tfjson.State{Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{Resources: []*tfjson.StateResource{{Address: "example_item.test", Type: "example_item", AttributeValues: map[string]any{"id": "secret"}}}}}})
	if !contains(f.Known, "id") {
		t.Fatalf("missing post-operation state: %v", f.Known)
	}
}

func TestFinishWritesOneFindingPerStep(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvDirectory, dir)
	c := New("TestExample")
	c.BeginStep(1)
	c.Parse("test.tf", []byte(`resource "example_item" "test" { name = "first" }`))
	if err := c.FinishStep(OutcomeSuccess); err != nil {
		t.Fatal(err)
	}
	// No active step: the end-of-test FinishStep must not write again.
	if err := c.FinishStep(OutcomeSuccess); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d findings", len(entries))
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"step":1`) || strings.Contains(string(data), "first") {
		t.Fatalf("unexpected finding: %s", data)
	}
}

func TestDiscardStepWritesNoFinding(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvDirectory, dir)
	c := New("TestExample")
	c.BeginStep(1)
	c.Parse("test.tf", []byte(`resource "example_item" "test" {}`))
	c.DiscardStep()
	// The discarded step is no longer active, so finishing it is a no-op.
	if err := c.FinishStep(OutcomeSuccess); err != nil {
		t.Fatal(err)
	}
	c.BeginStep(2)
	if err := c.FinishStep(OutcomeSuccess); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d findings", len(entries))
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"step":2`) {
		t.Fatalf("unexpected finding: %s", data)
	}
}

func TestCompletedLifecycleUsesAppliedPlan(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestLifecycle"}
	c.BeginStep(1)
	create := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.Plan(create)
	c.PreApplyPlan()
	c.Plan(noOp)
	c.ApplyCompleted(false)
	f := c.feature("resource", "example_item", "example_item.test")
	if !contains(f.Completed, "create") || !contains(f.Completed, "apply") || contains(f.Completed, "no_op") {
		t.Fatalf("completed actions: %v", f.Completed)
	}
	c.BeginStep(2)
	c.Plan(noOp)
	c.PreApplyPlan()
	c.ApplyCompleted(true)
	f = c.feature("resource", "example_item", "example_item.test")
	if len(f.Completed) != 0 {
		t.Fatalf("no-op destroy step claimed completion: %v", f.Completed)
	}
}

func TestPreApplyPlanWithoutReadablePlanWarns(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestUnavailable"}
	c.BeginStep(1)
	c.Plan(nil)
	c.PreApplyPlan()
	c.ApplyCompleted(false)
	if !contains(c.warnings, "plan_unavailable") || !contains(c.warnings, "pre_apply_plan_unavailable") {
		t.Fatalf("warnings: %v", c.warnings)
	}
	for _, f := range c.features {
		if len(f.Completed) != 0 {
			t.Fatalf("unreadable plan claimed completion: %v", f)
		}
	}
}

func TestNoOpApplyIsStillAnApply(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestNoOp"}
	c.BeginStep(1)
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.Plan(noOp)
	c.PreApplyPlan()
	c.ApplyCompleted(false)
	f := c.feature("resource", "example_item", "example_item.test")
	if !contains(f.Completed, "apply") || !contains(f.Completed, "no_op") {
		t.Fatalf("completed actions: %v", f.Completed)
	}
}

func TestSafeAddressHidesInstanceKeysContainingBrackets(t *testing.T) {
	t.Parallel()

	input := `module.test["secret]tail"].example_item.test["other-key"]`
	if got := safeAddress(input); got != "module.test[*].example_item.test[*]" {
		t.Fatalf("normalized address %q", got)
	}
}

func TestQueryConfigDoesNotIncludePreviousStep(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestQuery"}
	c.BeginStep(1)
	c.Parse("step.tf", []byte(`resource "example_item" "prior" { name = "previous" }`))
	c.BeginStep(2)
	c.ResetConfig()
	c.Parse("terraform_plugin_test.tfquery.hcl", []byte(`provider "example" {} list "example_item" "query" { config { name = "current" } }`))
	for _, f := range c.features {
		if f.Kind == "resource" || f.Kind == "provider" {
			t.Fatalf("stale or query-only feature: %+v", f)
		}
		if f.Kind == "list_resource" && !contains(f.Inputs, "name") {
			t.Fatalf("missing list config: %+v", f)
		}
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
