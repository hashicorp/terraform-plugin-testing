// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

func TestAggregateKeepsLifecycleAndAttributeEvidenceDistinct(t *testing.T) {
	t.Parallel()

	doc := schemaDoc{Providers: map[string]providerSchema{"registry.terraform.io/hashicorp/example": {Resources: map[string]schema{"example_item": {Block: &block{Attributes: map[string]attribute{"name": {Required: true}, "id": {Computed: true}}}}}}}}
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: "resource", Name: "example_item", Address: "example_item.test", Inputs: []string{"name"}, Planned: []string{"create"}}}}, {Version: 1, Test: "TestExample", Step: 2, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: "resource", Name: "example_item", Address: "example_item.test", Known: []string{"id"}, Completed: []string{"apply", "create"}}}}})
	if result.Inputs.Covered != 1 || result.Computed.Covered != 1 {
		t.Fatalf("coverage: %+v", result)
	}
	target := result.Targets[0]
	if len(target.Planned["create"]) != 1 || len(target.Completed["create"]) != 1 {
		t.Fatalf("lifecycle: %+v", target)
	}
}

func TestExpectedErrorDoesNotClaimCompletedOperationOrComputedOutput(t *testing.T) {
	t.Parallel()

	doc := schemaDoc{Providers: map[string]providerSchema{"registry.terraform.io/hashicorp/example": {Resources: map[string]schema{"example_item": {Block: &block{Attributes: map[string]attribute{"name": {Required: true}, "id": {Computed: true}}}}}}}}
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeExpectedError, Features: []semanticcoverage.Feature{{Kind: "resource", Name: "example_item", Inputs: []string{"name"}, Known: []string{"id"}, Planned: []string{"create"}, Completed: []string{"apply", "create"}}}}})
	if result.Inputs.Covered != 1 || result.Computed.Covered != 0 {
		t.Fatalf("input/computed coverage: %+v", result)
	}
	target := result.Targets[0]
	if len(target.Planned["create"]) != 1 || len(target.Completed) != 0 {
		t.Fatalf("planned/completed: %+v", target)
	}
}

func TestAggregateExcludesNestedAttributesFromCoverageTotals(t *testing.T) {
	t.Parallel()

	doc := schemaDoc{Providers: map[string]providerSchema{"registry.terraform.io/hashicorp/example": {
		Resources: map[string]schema{"example_item": {Block: &block{
			Attributes: map[string]attribute{
				"settings": {Optional: true, Nested: &nestedAttribute{Nesting: "single", Attributes: map[string]attribute{
					"mode":   {Optional: true},
					"status": {Computed: true},
				}}},
				"id": {Computed: true},
			},
			Blocks: map[string]nestedBlock{
				"connection": {Nesting: "list", Block: &block{Attributes: map[string]attribute{
					"host":  {Required: true},
					"state": {Computed: true},
				}}},
			},
		}}},
	}}}
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: "resource", Name: "example_item", Inputs: []string{"settings"}, Known: []string{"id"}}}}})
	if result.AttributeScope != "top_level" {
		t.Fatalf("attribute scope = %q, want top_level", result.AttributeScope)
	}
	if result.Inputs.Covered != 1 || result.Inputs.Total != 2 {
		t.Fatalf("input coverage should include only settings and connection: %+v", result.Inputs)
	}
	if result.Computed.Covered != 1 || result.Computed.Total != 1 {
		t.Fatalf("computed coverage should include only id: %+v", result.Computed)
	}
	if got, want := len(result.Targets[0].Attributes), 3; got != want {
		t.Fatalf("top-level attributes = %d, want %d: %+v", got, want, result.Targets[0].Attributes)
	}
	for _, a := range result.Targets[0].Attributes {
		if a.Path == "settings.mode" || a.Path == "settings.status" || a.Path == "connection[*].host" || a.Path == "connection[*].state" {
			t.Errorf("nested path %q should be excluded", a.Path)
		}
	}
}
