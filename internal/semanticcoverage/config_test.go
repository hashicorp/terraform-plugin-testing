// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnparsableConfigWarns(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), testName: "TestUnparsable"}
	c.BeginStep(1)
	c.observeConfigFile("broken.tf", []byte(`resource "example_item" {`))
	if !slices.Equal(c.warnings.sorted(), []Warning{WarningConfigUnavailable}) {
		t.Fatalf("warnings: %v", c.warnings)
	}
}

func TestJSONConfigWarnsAndIsNotRead(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.tf.json"), []byte(`{"resource": {"example_item": {"a": {"name": "x"}}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte(`resource "example_item" "b" { name = "x" }`), 0600); err != nil {
		t.Fatal(err)
	}
	c := &Collector{}
	c.BeginStep(1)
	c.ObserveConfigDir(dir)
	if !slices.Equal(c.warnings.sorted(), []Warning{WarningConfigUnavailable}) {
		t.Fatalf("warnings: %v", c.warnings)
	}
	var addresses []string
	for _, f := range c.features {
		addresses = append(addresses, f.address)
	}
	if !slices.Equal(addresses, []string{"example_item.b"}) {
		t.Fatalf("only the HCL file should be read, got features at %v", addresses)
	}
}

func TestQueryConfigDoesNotIncludePreviousStep(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), testName: "TestQuery"}
	c.BeginStep(1)
	c.observeConfigFile("step.tf", []byte(`resource "example_item" "prior" { name = "previous" }`))
	c.BeginStep(2)
	c.observeConfigFile("terraform_plugin_test.tfquery.hcl", []byte(`
provider "example" {}
list "example_item" "query" {
  config { name = "current" }
}
`))
	if len(c.features) != 1 || len(c.warnings) != 0 {
		t.Fatalf("features=%v warnings=%v", c.features, c.warnings)
	}
	for _, f := range c.features {
		if f.kind == KindResource || f.kind == KindProvider {
			t.Fatalf("stale or query-only feature: %+v", f)
		}
		if f.kind == KindListResource && !contains(f.inputs, "name") {
			t.Fatalf("missing list config: %+v", f)
		}
	}
}

func TestConfigBlocksDeclareFeatures(t *testing.T) {
	t.Parallel()

	hcl := `
provider "example" { region = "x" }
resource "example_item" "a" { name = "x" }
data "example_item" "b" { id = "x" }
ephemeral "example_secret" "c" { length = 1 }
list "example_item" "d" {
  config { filter = "x" }
}
action "example_action" "e" {
  config { message = "x" }
}
state_store "example_store" {
  config { path = "x" }
}
`
	want := []Feature{
		{Kind: KindAction, Name: "example_action", Address: "action.example_action.e", Inputs: []string{"message"}},
		{Kind: KindDataSource, Name: "example_item", Address: "data.example_item.b", Inputs: []string{"id"}},
		{Kind: KindEphemeralResource, Name: "example_secret", Address: "ephemeral.example_secret.c", Inputs: []string{"length"}},
		{Kind: KindListResource, Name: "example_item", Address: "list.example_item.d", Inputs: []string{"filter"}},
		{Kind: KindProvider, Name: "example", Inputs: []string{"region"}},
		{Kind: KindResource, Name: "example_item", Address: "example_item.a", Inputs: []string{"name"}},
		{Kind: KindStateStore, Name: "example_store", Inputs: []string{"path"}},
	}
	c := &Collector{}
	c.BeginStep(1)
	c.observeConfigFile("main.tf", []byte(hcl))
	var got []Feature
	for _, f := range c.features {
		got = append(got, f.finding())
	}
	sortFeatures(got)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("features mismatch (-want +got):\n%s", diff)
	}
}
