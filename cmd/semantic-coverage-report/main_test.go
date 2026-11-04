// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	sc "github.com/hashicorp/terraform-plugin-testing/internal/semanticcoverage"
)

func TestAggregateKeepsOperationAndAttributeEvidenceDistinct(t *testing.T) {
	t.Parallel()

	result := aggregate(parseSchema(t, exampleSchema), []sc.Finding{
		finding(1, item(func(f *sc.Feature) {
			f.Configured = []string{"name"}
			f.Planned = []sc.Operation{sc.OperationCreate}
		})),
		finding(2, item(func(f *sc.Feature) {
			f.Populated = []string{"id"}
			f.Completed = []sc.Operation{sc.OperationCreate}
		})),
	})
	if result.Totals.Configurable != (metric{Covered: 1, Total: 1}) || result.Totals.Computed != (metric{Covered: 1, Total: 1}) {
		t.Fatalf("totals: %+v", result.Totals)
	}
	f := result.Features[0]
	want := map[sc.Operation][]stepRef{sc.OperationCreate: {ref(1)}}
	if diff := cmp.Diff(want, f.Planned); diff != "" {
		t.Errorf("planned mismatch (-want +got):\n%s", diff)
	}
	want = map[sc.Operation][]stepRef{sc.OperationCreate: {ref(2)}}
	if diff := cmp.Diff(want, f.Completed); diff != "" {
		t.Errorf("completed mismatch (-want +got):\n%s", diff)
	}
}

func TestFeatureCoverage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		feature sc.Feature
		want    bool
	}{
		"used":                      {item(nil), true},
		"only another feature used": {sc.Feature{Kind: sc.KindDataSource, Name: "example_item"}, false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := aggregate(parseSchema(t, exampleSchema), []sc.Finding{finding(1, test.feature)})
			if got := result.Features[0].Covered; got != test.want {
				t.Errorf("covered = %t, want %t", got, test.want)
			}
		})
	}
}

func TestAggregateMatchesFeaturesToTheSchema(t *testing.T) {
	t.Parallel()

	const aws = "registry.terraform.io/hashicorp/aws"
	const google = "registry.terraform.io/hashicorp/google"
	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {
		"registry.terraform.io/hashicorp/aws": {
			"provider": {"block": {}},
			"resource_schemas": {"aws_instance": {"block": {}}},
			"functions": {"arn_parse": {"return_type": "string"}}
		},
		"registry.terraform.io/hashicorp/google": {
			"provider": {"block": {}},
			"functions": {"arn_parse": {"return_type": "string"}}
		}
	}}`)
	tests := map[string]struct {
		feature sc.Feature
		want    []string
	}{
		"provider block by local name":   {sc.Feature{Kind: sc.KindProvider, Name: "aws"}, []string{aws + " provider "}},
		"provider from a plan":           {sc.Feature{Kind: sc.KindProvider, Name: aws}, []string{aws + " provider "}},
		"resource from configuration":    {sc.Feature{Kind: sc.KindResource, Name: "aws_instance"}, []string{aws + " resource aws_instance"}},
		"resource from a plan":           {sc.Feature{Kind: sc.KindResource, Name: "aws_instance", Provider: aws}, []string{aws + " resource aws_instance"}},
		"resource of another provider":   {sc.Feature{Kind: sc.KindResource, Name: "aws_instance", Provider: "registry.terraform.io/example/aws"}, nil},
		"data source with resource name": {sc.Feature{Kind: sc.KindDataSource, Name: "aws_instance"}, nil},
		"function by local name":         {sc.Feature{Kind: sc.KindFunction, Name: "arn_parse", Provider: "google"}, []string{google + " function arn_parse"}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := aggregate(doc, []sc.Finding{finding(1, test.feature)})
			var covered []string
			for _, f := range result.Features {
				if f.Covered {
					covered = append(covered, f.Provider+" "+string(f.Kind)+" "+f.Name)
				}
			}
			if diff := cmp.Diff(test.want, covered); diff != "" {
				t.Errorf("covered features mismatch (-want +got):\n%s", diff)
			}
			if len(result.Warnings) != 0 {
				t.Errorf("warnings: %v", result.Warnings)
			}
		})
	}
}

// publishedSchema is exampleSchema under the address the provider is
// published under, which differs from the address the test harness serves it
// under.
const publishedSchema = `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/example-org/example": {
	"provider": {"block": {}},
	"resource_schemas": {"example_item": {"block": {"attributes": {
		"name": {"type": "string", "required": true},
		"id":   {"type": "string", "computed": true}
	}}}}
}}}`

// servedFinding returns a finding of step 1 of TestExample that served the
// given providers and used example_item, as a plan and state record it.
func servedFinding(served ...string) sc.Finding {
	const address = "registry.terraform.io/hashicorp/example"
	f := finding(1,
		sc.Feature{Kind: sc.KindProvider, Name: address},
		item(func(f *sc.Feature) {
			f.Provider = address
			f.Configured = []string{"name"}
			f.Populated = []string{"id", "name"}
			f.Completed = []sc.Operation{sc.OperationCreate}
		}),
	)
	f.ServedProviders = append([]string{}, served...)
	return f
}

func TestServedProviderMapsToThePublishedAddress(t *testing.T) {
	t.Parallel()

	result := aggregate(parseSchema(t, publishedSchema), []sc.Finding{servedFinding("registry.terraform.io/hashicorp/example")})
	want := totals{Features: metric{2, 2}, Configurable: metric{1, 1}, Computed: metric{1, 1}}
	if result.Totals != want {
		t.Errorf("totals = %+v, want %+v", result.Totals, want)
	}
	if got := result.Features[1].Completed[sc.OperationCreate]; len(got) != 1 {
		t.Errorf("completed create = %v, want step 1", got)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("warnings: %v", result.Warnings)
	}
}

func TestProvidersNotServedFromTheTestBinaryAreExcluded(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		served []string
		want   []string
	}{
		// An upgrade step runs a released provider that Terraform
		// installs. Its address can equal the schema address exactly.
		"no provider served": {
			served: nil,
			want:   []string{"no provider served from the test binary matches a single schema provider (served: none); generate the schema from the provider under test"},
		},
		// Another provider of the same type, such as an external one,
		// isn't the provider under test.
		"another provider served": {
			served: []string{"registry.terraform.io/hashicorp/other"},
			want:   []string{"no provider served from the test binary matches a single schema provider (served: registry.terraform.io/hashicorp/other); generate the schema from the provider under test"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := aggregate(parseSchema(t, exampleSchema), []sc.Finding{servedFinding(test.served...)})
			if result.Totals.Features.Covered != 0 {
				t.Errorf("features = %+v, want none covered", result.Totals.Features)
			}
			if diff := cmp.Diff(test.want, result.Warnings); diff != "" {
				t.Errorf("warnings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestServedHelperProvidersAreIgnored(t *testing.T) {
	t.Parallel()

	const echo = "registry.terraform.io/hashicorp/echo"
	f := servedFinding("registry.terraform.io/hashicorp/example", echo)
	f.Features = append(f.Features,
		sc.Feature{Kind: sc.KindProvider, Name: "echo"},
		sc.Feature{Kind: sc.KindProvider, Name: echo},
		sc.Feature{Kind: sc.KindResource, Name: "echo", Provider: echo, Configured: []string{"data"}},
	)
	result := aggregate(parseSchema(t, publishedSchema), []sc.Finding{f})
	if result.Totals.Features != (metric{2, 2}) {
		t.Errorf("features = %+v, want both covered", result.Totals.Features)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("warnings: %v", result.Warnings)
	}
}

func TestServedProviderOfAnAmbiguousTypeIsNotMapped(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {
		"registry.terraform.io/example-org/example": {"resource_schemas": {"example_item": {"block": {}}}},
		"registry.terraform.io/fork/example": {"resource_schemas": {"example_item": {"block": {}}}}
	}}`)
	result := aggregate(doc, []sc.Finding{servedFinding("registry.terraform.io/hashicorp/example")})
	if result.Totals.Features.Covered != 0 {
		t.Errorf("features = %+v, want none covered", result.Totals.Features)
	}
	want := []string{"no provider served from the test binary matches a single schema provider (served: registry.terraform.io/hashicorp/example); generate the schema from the provider under test"}
	if diff := cmp.Diff(want, result.Warnings); diff != "" {
		t.Errorf("warnings mismatch (-want +got):\n%s", diff)
	}
}

func TestAmbiguousConfigurationIsNotCounted(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {
		"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {}}}},
		"registry.terraform.io/fork/example": {"resource_schemas": {"example_item": {"block": {}}}}
	}}`)
	result := aggregate(doc, []sc.Finding{finding(1, sc.Feature{Kind: sc.KindResource, Name: "example_item"})})
	if result.Totals.Features.Covered != 0 {
		t.Errorf("features = %+v, want none covered", result.Totals.Features)
	}
	want := []string{"TestExample step 1: resource example_item has 2 schema matches"}
	if diff := cmp.Diff(want, result.Warnings); diff != "" {
		t.Errorf("warnings mismatch (-want +got):\n%s", diff)
	}
}

func TestOptionalComputedAttributesNeedATestThatSetsThemAndOneThatDoesNot(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {"attributes": {
		"tags": {"type": "string", "optional": true, "computed": true}
	}}}}}}}`)
	unset := finding(1, item(func(f *sc.Feature) { f.Populated = []string{"tags"} }))
	set := finding(2, item(func(f *sc.Feature) {
		f.Configured = []string{"tags"}
		f.Populated = []string{"tags"}
	}))

	got := aggregate(doc, []sc.Finding{unset, set}).Features[0].Attributes[0]
	want := attributeResult{Path: "tags", Configurable: true, Computed: true, Covered: true, ConfiguredBy: []stepRef{ref(2)}, ComputedBy: []stepRef{ref(1)}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("attribute mismatch (-want +got):\n%s", diff)
	}

	// A configured value is echoed into state, so it isn't computed.
	got = aggregate(doc, []sc.Finding{set}).Features[0].Attributes[0]
	want = attributeResult{Path: "tags", Configurable: true, Computed: true, ConfiguredBy: []stepRef{ref(2)}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("attribute mismatch (-want +got):\n%s", diff)
	}
}

func TestAggregateCountsDeprecatedAttributes(t *testing.T) {
	t.Parallel()

	doc := parseSchema(t, `{"format_version": "1.0", "provider_schemas": {"registry.terraform.io/hashicorp/example": {"resource_schemas": {"example_item": {"block": {
		"attributes": {"old_name": {"type": "string", "optional": true, "deprecated": true}},
		"deprecated": true
	}}}}}}`)
	result := aggregate(doc, []sc.Finding{finding(1, item(func(f *sc.Feature) { f.Configured = []string{"old_name"} }))})
	if result.Totals.Configurable != (metric{Covered: 1, Total: 1}) {
		t.Fatalf("deprecated attribute of a deprecated resource should count: %+v", result.Totals.Configurable)
	}
}

func TestAggregateExcludesNestedAttributes(t *testing.T) {
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
	result := aggregate(doc, []sc.Finding{finding(1, item(func(f *sc.Feature) {
		f.Configured = []string{"settings"}
		f.Populated = []string{"id"}
	}))})
	if result.AttributeScope != "top_level" {
		t.Fatalf("attribute scope = %q, want top_level", result.AttributeScope)
	}
	var paths []string
	for _, a := range result.Features[0].Attributes {
		paths = append(paths, a.Path)
	}
	if diff := cmp.Diff([]string{"connection", "id", "settings"}, paths); diff != "" {
		t.Errorf("attributes mismatch (-want +got):\n%s", diff)
	}
	if result.Totals.Configurable != (metric{Covered: 1, Total: 2}) || result.Totals.Computed != (metric{Covered: 1, Total: 1}) {
		t.Fatalf("totals: %+v", result.Totals)
	}
}

func TestWarningsAreInStepOrder(t *testing.T) {
	t.Parallel()

	step := func(n int) sc.Finding {
		f := finding(n)
		f.Package = "example.com/provider"
		f.Warnings = []sc.Warning{sc.WarningStateUnavailable}
		return f
	}
	result := aggregate(parseSchema(t, exampleSchema), []sc.Finding{step(10), step(2)})
	want := []string{
		"example.com/provider TestExample step 2: state_unavailable",
		"example.com/provider TestExample step 10: state_unavailable",
	}
	if diff := cmp.Diff(want, result.Warnings); diff != "" {
		t.Errorf("warnings mismatch (-want +got):\n%s", diff)
	}
}

func TestTextReport(t *testing.T) {
	t.Parallel()

	unplanned := finding(2)
	unplanned.Warnings = []sc.Warning{sc.WarningPlanUnavailable}
	result := aggregate(parseSchema(t, exampleSchema), []sc.Finding{
		finding(1, item(func(f *sc.Feature) {
			f.Configured = []string{"name"}
			f.Planned = []sc.Operation{sc.OperationCreate}
			f.Completed = []sc.Operation{sc.OperationCreate}
		})),
		unplanned,
	})
	want := `Test-derived semantic coverage: 2 findings
Attribute scope: top-level schema attributes and blocks only; nested descendants are excluded.
Features 1/1; configurable 1/1; computed 0/1
registry.terraform.io/hashicorp/example resource example_item: covered=true attributes=1/2 planned=[create] completed=[create]
  uncovered id (needs a test that leaves it to the provider)
warning: TestExample step 2: plan_unavailable
`
	if diff := cmp.Diff(want, textReport(result)); diff != "" {
		t.Errorf("text report mismatch (-want +got):\n%s", diff)
	}
}

func TestReadFindingsRejectsInvalidFindings(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		finding string
		want    string
	}{
		"not JSON":          {`{`, "unexpected end of JSON input"},
		"other version":     {`{"version":2,"test":"TestExample","step":1}`, "finding version 2, this report reads version 1"},
		"no test name":      {`{"version":1,"step":1}`, "missing test name or step"},
		"no step":           {`{"version":1,"test":"TestExample"}`, "missing test name or step"},
		"unknown warning":   {`{"version":1,"test":"TestExample","step":1,"warnings":["step_skipped"]}`, `invalid warning "step_skipped"`},
		"unknown kind":      {`{"version":1,"test":"TestExample","step":1,"features":[{"kind":"identity","name":"x"}]}`, `invalid feature kind "identity"`},
		"unknown operation": {`{"version":1,"test":"TestExample","step":1,"features":[{"kind":"resource","name":"x","completed":["apply"]}]}`, `invalid operation "apply"`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			path := filepath.Join(dir, "finding.json")
			if err := os.WriteFile(path, []byte(test.finding), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := readFindings(dir)
			if err == nil || !strings.HasPrefix(err.Error(), path+": ") || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("err = %v, want %q", err, test.want)
			}
		})
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

// testProviders are the providers that findings returned by finding served
// from the test binary: those of the schemas in these tests.
var testProviders = []string{
	"registry.terraform.io/hashicorp/aws",
	"registry.terraform.io/hashicorp/example",
	"registry.terraform.io/hashicorp/google",
}

// finding returns a finding of step n of TestExample that served
// testProviders.
func finding(n int, features ...sc.Feature) sc.Finding {
	return sc.Finding{Version: sc.FindingVersion, Test: "TestExample", Step: n, ServedProviders: testProviders, Features: features}
}

// item returns the resource example_item of exampleSchema, changed by fn.
func item(fn func(f *sc.Feature)) sc.Feature {
	f := sc.Feature{Kind: sc.KindResource, Name: "example_item"}
	if fn != nil {
		fn(&f)
	}
	return f
}

func ref(step int) stepRef {
	return stepRef{Test: "TestExample", Step: step}
}
