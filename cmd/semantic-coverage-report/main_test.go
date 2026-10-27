// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

func TestAggregateKeepsLifecycleAndAttributeEvidenceDistinct(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, exampleSchema)
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Address: "example_item.test", Inputs: []string{"name"}, Planned: []string{"create"}}}}, {Version: 1, Test: "TestExample", Step: 2, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Address: "example_item.test", Known: []string{"id"}, Completed: []string{"apply", "create"}}}}})
	if result.Configurable.Covered != 1 || result.Computed.Covered != 1 {
		t.Fatalf("coverage: %+v", result)
	}
	target := result.Targets[0]
	if len(target.Planned["create"]) != 1 || len(target.Completed["create"]) != 1 {
		t.Fatalf("lifecycle: %+v", target)
	}
}

func TestExpectedErrorDoesNotClaimCompletedOperationOrComputedOutput(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, exampleSchema)
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeExpectedError, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Inputs: []string{"name"}, Known: []string{"id"}, Planned: []string{"create"}, Completed: []string{"apply", "create"}}}}})
	if result.Configurable.Covered != 1 || result.Computed.Covered != 0 {
		t.Fatalf("configurable/computed coverage: %+v", result)
	}
	target := result.Targets[0]
	if len(target.Planned["create"]) != 1 || len(target.Completed) != 0 {
		t.Fatalf("planned/completed: %+v", target)
	}
}

func TestAggregateMergesFindingsForTheSameStep(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, exampleSchema)
	feature := semanticcoverage.Feature{Kind: semanticcoverage.KindResource, Name: "example_item", Inputs: []string{"name"}, Known: []string{"id"}}
	result := aggregate(doc, []semanticcoverage.Finding{
		{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeFailed, Features: []semanticcoverage.Feature{feature}},
		{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{feature}},
	})
	if result.Configurable.Covered != 1 || result.Computed.Covered != 1 || !result.Targets[0].Exercised {
		t.Fatalf("a failed finding must not cancel another finding's coverage: %+v", result)
	}
}

func TestTargetMatches(t *testing.T) {
	t.Parallel()

	const aws = "registry.terraform.io/hashicorp/aws"
	provider := &target{Provider: aws, Kind: semanticcoverage.KindProvider}
	resource := &target{Provider: aws, Kind: semanticcoverage.KindResource, Name: "aws_instance"}
	function := &target{Provider: aws, Kind: semanticcoverage.KindFunction, Name: "arn_parse"}
	tests := map[string]struct {
		target *target
		f      semanticcoverage.Feature
		want   bool
	}{
		"provider block by local name":     {provider, semanticcoverage.Feature{Kind: semanticcoverage.KindProvider, Name: "aws"}, true},
		"provider block of other provider": {provider, semanticcoverage.Feature{Kind: semanticcoverage.KindProvider, Name: "google"}, false},
		"resource from config":             {resource, semanticcoverage.Feature{Kind: semanticcoverage.KindResource, Name: "aws_instance"}, true},
		"resource from plan":               {resource, semanticcoverage.Feature{Kind: semanticcoverage.KindResource, Name: "aws_instance", Provider: aws}, true},
		"resource of fork from plan":       {resource, semanticcoverage.Feature{Kind: semanticcoverage.KindResource, Name: "aws_instance", Provider: "registry.terraform.io/example/aws"}, false},
		"data source with resource name":   {resource, semanticcoverage.Feature{Kind: semanticcoverage.KindDataSource, Name: "aws_instance"}, false},
		"function by local name":           {function, semanticcoverage.Feature{Kind: semanticcoverage.KindFunction, Name: "arn_parse", Provider: "aws"}, true},
		"function of other provider":       {function, semanticcoverage.Feature{Kind: semanticcoverage.KindFunction, Name: "arn_parse", Provider: "google"}, false},
	}
	for name, test := range tests {
		if got := test.target.matches(test.f); got != test.want {
			t.Errorf("%s: matches = %t, want %t", name, got, test.want)
		}
	}
}

const optionalComputedSchema = `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {"attributes": {
	"tags": {"type": "string", "optional": true, "computed": true}
}}}}}}}`

func TestAggregateNeedsOptionalComputedAttributesSetAndUnset(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, optionalComputedSchema)
	unset := semanticcoverage.Finding{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Address: "example_item.test", Known: []string{"tags"}}}}
	set := semanticcoverage.Finding{Version: 1, Test: "TestExample", Step: 2, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Address: "example_item.test", Inputs: []string{"tags"}, Known: []string{"tags"}}}}

	result := aggregate(doc, []semanticcoverage.Finding{unset, set})
	if result.Configurable != (metric{Covered: 1, Total: 1}) || result.Computed != (metric{Covered: 1, Total: 1}) {
		t.Fatalf("configurable=%+v computed=%+v, want both 1/1", result.Configurable, result.Computed)
	}
	if a := result.Targets[0].Attributes[0]; !a.Covered {
		t.Fatalf("attribute should be covered: %+v", a)
	}

	result = aggregate(doc, []semanticcoverage.Finding{set})
	if result.Configurable != (metric{Covered: 1, Total: 1}) || result.Computed != (metric{Total: 1}) {
		t.Fatalf("a configured value echoed into state is not computed: configurable=%+v computed=%+v", result.Configurable, result.Computed)
	}
	if a := result.Targets[0].Attributes[0]; a.Covered || strings.Join(a.missing(), ",") != "state_known" {
		t.Fatalf("attribute should only miss state_known: %+v", a)
	}
}

func TestAggregateIgnoresStateOfTypesConfiguredInTheSameStep(t *testing.T) {
	t.Parallel()

	// With count or for_each, the configured address has no instance key and
	// the state address has one, so they are different features.
	doc := parseSchema(t, optionalComputedSchema)
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{
		{Kind: semanticcoverage.KindResource, Name: "example_item", Address: "example_item.test", Inputs: []string{"tags"}},
		{Kind: semanticcoverage.KindResource, Name: "example_item", Address: "example_item.test[*]", Known: []string{"tags"}},
	}}})
	if result.Computed.Covered != 0 {
		t.Fatalf("computed = %+v, want no coverage", result.Computed)
	}
}

func TestReadFindingsRejectsUnknownWarning(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	finding := `{"version":1,"test":"TestExample","step":1,"outcome":"success","warnings":["step_skipped"]}`
	if err := os.WriteFile(filepath.Join(dir, "finding.json"), []byte(finding), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readFindings(dir); err == nil || !strings.Contains(err.Error(), `invalid warning "step_skipped"`) {
		t.Fatalf("err = %v", err)
	}
}

func TestReadFindingsRejectsUnknownKind(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	finding := `{"version":1,"test":"TestExample","step":1,"outcome":"success","features":[{"kind":"test_step","name":"x"}]}`
	if err := os.WriteFile(filepath.Join(dir, "finding.json"), []byte(finding), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readFindings(dir); err == nil || !strings.Contains(err.Error(), `invalid feature kind "test_step"`) {
		t.Fatalf("err = %v", err)
	}
}

func TestAggregateReadsActionsAndLeavesIdentityOutOfAttributeTotals(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {
		"action_schemas": {"example_action": {"block": {"attributes": {"message": {"type": "string", "required": true}}}}},
		"resource_identity_schemas": {"example_item": {"version": 0, "attributes": {"id": {"type": "string", "required_for_import": true}}}},
		"functions": {"example_function": {"return_type": "string"}}
	}}}`)
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{
		{Kind: semanticcoverage.KindAction, Name: "example_action", Inputs: []string{"message"}},
		{Kind: semanticcoverage.KindIdentity, Name: "example_item", Known: []string{"id"}},
	}}})
	if result.Configurable != (metric{Covered: 1, Total: 1}) || result.Computed != (metric{}) {
		t.Fatalf("configurable=%+v computed=%+v", result.Configurable, result.Computed)
	}
	kinds := map[semanticcoverage.Kind]*target{}
	for _, target := range result.Targets {
		kinds[target.Kind] = target
	}
	if kinds[semanticcoverage.KindIdentity] == nil || len(kinds[semanticcoverage.KindIdentity].Attributes) != 0 || !kinds[semanticcoverage.KindIdentity].Exercised {
		t.Fatalf("identity target: %+v", kinds[semanticcoverage.KindIdentity])
	}
	if kinds[semanticcoverage.KindFunction] == nil {
		t.Fatalf("missing function target: %+v", result.Targets)
	}
}

func TestAggregateCountsDeprecatedAttributes(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {
		"attributes": {"old_name": {"type": "string", "optional": true, "deprecated": true}},
		"deprecated": true
	}}}}}}`)
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Inputs: []string{"old_name"}}}}})
	if result.Configurable != (metric{Covered: 1, Total: 1}) {
		t.Fatalf("deprecated attribute of a deprecated resource should count: %+v", result.Configurable)
	}
}

func TestAggregateExcludesNestedAttributesFromCoverageTotals(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {
		"attributes": {
			"settings": {"optional": true, "nested_type": {"nesting_mode": "single", "attributes": {
				"mode":   {"type": "string", "optional": true},
				"status": {"type": "string", "computed": true}
			}}},
			"id": {"type": "string", "computed": true}
		},
		"block_types": {
			"connection": {"nesting_mode": "list", "block": {"attributes": {
				"host":  {"type": "string", "required": true},
				"state": {"type": "string", "computed": true}
			}}}
		}
	}}}}}}`)
	result := aggregate(doc, []semanticcoverage.Finding{{Version: 1, Test: "TestExample", Step: 1, Outcome: semanticcoverage.OutcomeSuccess, Features: []semanticcoverage.Feature{{Kind: semanticcoverage.KindResource, Name: "example_item", Inputs: []string{"settings"}, Known: []string{"id"}}}}})
	if result.AttributeScope != "top_level" {
		t.Fatalf("attribute scope = %q, want top_level", result.AttributeScope)
	}
	if result.Configurable.Covered != 1 || result.Configurable.Total != 2 {
		t.Fatalf("configurable coverage should include only settings and connection: %+v", result.Configurable)
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

const exampleSchema = `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {"attributes": {
	"name": {"type": "string", "required": true},
	"id":   {"type": "string", "computed": true}
}}}}}}}`

func parseSchema(t *testing.T, data string) *tfjson.ProviderSchemas {
	t.Helper()

	var doc tfjson.ProviderSchemas
	if err := json.Unmarshal([]byte(data), &doc); err != nil {
		t.Fatal(err)
	}
	return &doc
}
