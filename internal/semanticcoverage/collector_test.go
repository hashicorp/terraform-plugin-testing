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
	c.Begin(1)
	c.Parse("test.tf", []byte(`resource "example_item" "test" { token = "secret-configuration" }`))
	plan := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: `example_item.test["secret-instance-key"]`, Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	c.Plan(plan)
	c.PreApplyPlan(plan)
	resource := &tfjson.StateResource{Address: `example_item.test["secret-instance-key"]`, Type: "example_item", AttributeValues: map[string]any{"id": "secret-value", "token": "secret-value", "map": map[string]any{"secret-map-key": "secret-value"}}}
	c.ObserveState(&tfjson.State{Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{Resources: []*tfjson.StateResource{resource}}}})
	c.ApplyCompleted(false)
	if err := c.Finish("success"); err != nil {
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
	c := &Collector{}
	if err := c.Begin(1); err != nil {
		t.Fatal(err)
	}
	c.Plan(&tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}})
	f := c.feature("resource", "example_item", "example_item.test")
	if len(f.Planned) != 1 || len(f.Completed) != 0 || len(f.Known) != 0 {
		t.Fatalf("planned=%v completed=%v known=%v", f.Planned, f.Completed, f.Known)
	}
}

func TestStateRequiresExplicitPostOperationObservation(t *testing.T) {
	c := &Collector{}
	if err := c.Begin(1); err != nil {
		t.Fatal(err)
	}
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

func TestBeginFinishesPreviousStep(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvDirectory, dir)
	c := New("TestExample")
	if err := c.Begin(1); err != nil {
		t.Fatal(err)
	}
	c.Parse("test.tf", []byte(`resource "example_item" "test" { name = "first" }`))
	if err := c.Begin(2); err != nil {
		t.Fatal(err)
	}
	if err := c.FinishStep(1, false, false); err != nil {
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

func TestCompletedLifecycleUsesAppliedPlan(t *testing.T) {
	c := &Collector{directory: t.TempDir(), test: "TestLifecycle"}
	if err := c.Begin(1); err != nil {
		t.Fatal(err)
	}
	create := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.PreApplyPlan(create)
	c.Plan(noOp)
	c.ApplyCompleted(false)
	f := c.feature("resource", "example_item", "example_item.test")
	if !contains(f.Completed, "create") || !contains(f.Completed, "apply") || contains(f.Completed, "no_op") {
		t.Fatalf("completed actions: %v", f.Completed)
	}
	if err := c.Begin(2); err != nil {
		t.Fatal(err)
	}
	c.PreApplyPlan(noOp)
	c.ApplyCompleted(true)
	f = c.feature("resource", "example_item", "example_item.test")
	if len(f.Completed) != 0 {
		t.Fatalf("no-op destroy step claimed completion: %v", f.Completed)
	}
}

func TestNoOpApplyIsStillAnApply(t *testing.T) {
	c := &Collector{directory: t.TempDir(), test: "TestNoOp"}
	if err := c.Begin(1); err != nil {
		t.Fatal(err)
	}
	noOp := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}}}}
	c.PreApplyPlan(noOp)
	c.ApplyCompleted(false)
	f := c.feature("resource", "example_item", "example_item.test")
	if !contains(f.Completed, "apply") || !contains(f.Completed, "no_op") {
		t.Fatalf("completed actions: %v", f.Completed)
	}
}

func TestSafeAddressHidesInstanceKeysContainingBrackets(t *testing.T) {
	input := `module.test["secret]tail"].example_item.test["other-key"]`
	if got := safeAddress(input); got != "module.test[*].example_item.test[*]" {
		t.Fatalf("normalized address %q", got)
	}
}

func TestQueryConfigDoesNotIncludePreviousStep(t *testing.T) {
	c := &Collector{directory: t.TempDir(), test: "TestQuery"}
	if err := c.Begin(1); err != nil {
		t.Fatal(err)
	}
	c.Parse("step.tf", []byte(`resource "example_item" "prior" { name = "previous" }`))
	if err := c.Begin(2); err != nil {
		t.Fatal(err)
	}
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
