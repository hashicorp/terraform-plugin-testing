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
	c, err := New("TestExample")
	if err != nil {
		t.Fatal(err)
	}
	c.BeginStep(1)
	c.observeConfigFile("test.tf", []byte(`resource "example_item" "test" { token = "secret-configuration" }`))
	plan := &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: `example_item.test["secret-instance-key"]`, Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}}}}
	c.ObservePlan(plan)
	c.PreApplyPlan(plan)
	resource := &tfjson.StateResource{Address: `example_item.test["secret-instance-key"]`, Type: "example_item", Mode: tfjson.ManagedResourceMode, AttributeValues: map[string]any{"id": "secret-value", "token": "secret-value", "map": map[string]any{"secret-map-key": "secret-value"}}}
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

func TestFinishWritesOneFindingPerStep(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvDirectory, dir)
	c, err := New("TestExample")
	if err != nil {
		t.Fatal(err)
	}
	c.BeginStep(1)
	c.observeConfigFile("test.tf", []byte(`resource "example_item" "test" { name = "first" }`))
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
	c, err := New("TestExample")
	if err != nil {
		t.Fatal(err)
	}
	c.BeginStep(1)
	c.observeConfigFile("test.tf", []byte(`resource "example_item" "test" {}`))
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

func contains(values set[string], wanted string) bool {
	_, ok := values[wanted]
	return ok
}

func TestNewFailsEarlyForUnusableDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvDirectory, filepath.Join(file, "findings"))
	if c, err := New("TestExample"); err == nil || c != nil {
		t.Fatalf("collector=%v err=%v", c, err)
	}
	t.Setenv(EnvDirectory, "")
	if c, err := New("TestExample"); err != nil || c != nil {
		t.Fatalf("disabled: collector=%v err=%v", c, err)
	}
}
