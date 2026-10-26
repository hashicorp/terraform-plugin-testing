// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestUnparsableConfigWarns(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestUnparsable"}
	c.BeginStep(1)
	c.Parse("broken.tf", []byte(`resource "example_item" {`))
	c.Parse("broken.tf.json", []byte(`{`))
	if !slices.Equal(c.warnings.sorted(), []Warning{WarningConfigUnavailable}) {
		t.Fatalf("warnings: %v", c.warnings)
	}
}

func TestQueryConfigDoesNotIncludePreviousStep(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestQuery"}
	c.BeginStep(1)
	c.Parse("step.tf", []byte(`resource "example_item" "prior" { name = "previous" }`))
	c.BeginStep(2)
	c.Parse("terraform_plugin_test.tfquery.hcl", []byte(`
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

func TestConfigBlocksDeclareFeaturesInHCLAndJSON(t *testing.T) {
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
	json := `{
  "provider": {"example": {"region": "x"}},
  "resource": {"example_item": {"a": {"name": "x"}}},
  "data": {"example_item": {"b": {"id": "x"}}},
  "ephemeral": {"example_secret": {"c": {"length": 1}}},
  "list": {"example_item": {"d": {"config": {"filter": "x"}}}},
  "action": {"example_action": {"e": {"config": {"message": "x"}}}},
  "state_store": {"example_store": {"config": {"path": "x"}}}
}`
	want := []Feature{
		{Kind: KindAction, Name: "example_action", Address: "action.example_action.e", Inputs: []string{"message"}},
		{Kind: KindDataSource, Name: "example_item", Address: "data.example_item.b", Inputs: []string{"id"}},
		{Kind: KindEphemeralResource, Name: "example_secret", Address: "ephemeral.example_secret.c", Inputs: []string{"length"}},
		{Kind: KindListResource, Name: "example_item", Address: "list.example_item.d", Inputs: []string{"filter"}},
		{Kind: KindProvider, Name: "example", Inputs: []string{"region"}},
		{Kind: KindResource, Name: "example_item", Address: "example_item.a", Inputs: []string{"name"}},
		{Kind: KindStateStore, Name: "example_store", Inputs: []string{"path"}},
	}
	for name, data := range map[string]string{"main.tf": hcl, "main.tf.json": json} {
		c := &Collector{}
		c.BeginStep(1)
		c.Parse(name, []byte(data))
		var got []Feature
		for _, f := range c.features {
			got = append(got, f.finding())
		}
		slices.SortFunc(got, func(a, b Feature) int {
			return strings.Compare(string(a.Kind)+a.Address, string(b.Kind)+b.Address)
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\ngot  %+v\nwant %+v", name, got, want)
		}
	}
}
