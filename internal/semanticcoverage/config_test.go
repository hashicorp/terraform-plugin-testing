// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"slices"
	"testing"
)

func TestUnparsableConfigWarns(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestUnparsable"}
	c.BeginStep(1)
	c.Parse("broken.tf", []byte(`resource "example_item" {`))
	c.Parse("broken.tf.json", []byte(`{`))
	if !slices.Equal(c.warnings, []Warning{WarningConfigUnavailable}) {
		t.Fatalf("warnings: %v", c.warnings)
	}
}

func TestQueryConfigDoesNotIncludePreviousStep(t *testing.T) {
	t.Parallel()

	c := &Collector{directory: t.TempDir(), test: "TestQuery"}
	c.BeginStep(1)
	c.Parse("step.tf", []byte(`resource "example_item" "prior" { name = "previous" }`))
	c.BeginStep(2)
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
