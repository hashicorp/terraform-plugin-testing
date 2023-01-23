package plugintest

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

var terraformJSONOutput = []string{
	`{"@level":"info","@message":"Terraform 1.3.2","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.232751Z","terraform":"1.3.2","type":"version","ui":"1.0"}`,
	`{"@level":"warn","@message":"Warning: Empty or non-existent state","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.250725Z","diagnostic":{"severity":"warning","summary":"Empty or non-existent state","detail":"There are currently no remote objects tracked in the state, so there is nothing to refresh."},"type":"diagnostic"}`,
	`{"@level":"info","@message":"Outputs: 0","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.251120Z","outputs":{},"type":"outputs"}`,
	`{"@level":"info","@message":"random_password.test: Plan to create","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.282021Z","change":{"resource":{"addr":"random_password.test","module":"","resource":"random_password.test","implied_provider":"random","resource_type":"random_password","resource_name":"test","resource_key":null},"action":"create"},"type":"planned_change"}`,
	`{"@level":"info","@message":"Plan: 1 to add, 0 to change, 0 to destroy.","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.282054Z","changes":{"add":1,"change":0,"remove":0,"operation":"plan"},"type":"change_summary"}`,
	`{"@level":"error","@message":"Error: error diagnostic - summary","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.626572Z","diagnostic":{"severity":"error","summary":"error diagnostic - summary","detail":""},"type":"diagnostic"}`,
}

func TestTerraformJSONDiagnostics_Contains(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		diags    []tfjson.Diagnostic
		regex    *regexp.Regexp
		severity tfjson.DiagnosticSeverity
		expected bool
	}{
		"severity-not-found": {
			diags: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error summary",
					Detail:   "error detail",
				},
			},
			regex:    regexp.MustCompile("error summary"),
			severity: tfjson.DiagnosticSeverityWarning,
			expected: false,
		},
		"summary-not-found": {
			diags: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "warning summary",
					Detail:   "warning detail",
				},
			},
			regex:    regexp.MustCompile("error detail"),
			severity: tfjson.DiagnosticSeverityError,
			expected: false,
		},
		"summary-found": {
			diags: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error summary",
					Detail:   "error detail",
				},
			},
			regex:    regexp.MustCompile("error summary"),
			severity: tfjson.DiagnosticSeverityError,
			expected: true,
		},
		"detail-found": {
			diags: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error summary",
					Detail:   "error detail",
				},
			},
			regex:    regexp.MustCompile("error detail"),
			severity: tfjson.DiagnosticSeverityError,
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			var tfJSONDiagnostics TerraformJSONDiagnostics = testCase.diags

			isFound := tfJSONDiagnostics.Contains(testCase.regex, testCase.severity)

			if diff := cmp.Diff(isFound, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTerraformJSONBuffer_Parse(t *testing.T) {
	t.Parallel()

	tfJSON := NewTerraformJSONBuffer()

	for _, v := range terraformJSONOutput {
		_, err := tfJSON.Write([]byte(v + "\n"))
		if err != nil {
			t.Fatalf("cannot write to tfJSON: %s", err)
		}
	}

	tfJSON.Parse()

	if diff := cmp.Diff(tfJSON.jsonOutput, terraformJSONOutput); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}

	var tfJSONDiagnostics TerraformJSONDiagnostics = []tfjson.Diagnostic{
		{
			Severity: tfjson.DiagnosticSeverityWarning,
			Summary:  "Empty or non-existent state",
			Detail:   "There are currently no remote objects tracked in the state, so there is nothing to refresh.",
		},
		{
			Severity: tfjson.DiagnosticSeverityError,
			Summary:  "error diagnostic - summary",
		},
	}

	if diff := cmp.Diff(tfJSON.diagnostics, tfJSONDiagnostics); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func TestTerraformJSONBuffer_Diagnostics(t *testing.T) {
	t.Parallel()

	tfJSON := NewTerraformJSONBuffer()

	for _, v := range terraformJSONOutput {
		_, err := tfJSON.Write([]byte(v + "\n"))
		if err != nil {
			t.Fatalf("cannot write to tfJSON: %s", err)
		}
	}

	tfJSON.Parse()

	var tfJSONDiagnostics TerraformJSONDiagnostics = []tfjson.Diagnostic{
		{
			Severity: tfjson.DiagnosticSeverityWarning,
			Summary:  "Empty or non-existent state",
			Detail:   "There are currently no remote objects tracked in the state, so there is nothing to refresh.",
		},
		{
			Severity: tfjson.DiagnosticSeverityError,
			Summary:  "error diagnostic - summary",
		},
	}

	if diff := cmp.Diff(tfJSON.Diagnostics(), tfJSONDiagnostics); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func TestTerraformJSONBuffer_JsonOutput(t *testing.T) {
	t.Parallel()

	tfJSON := NewTerraformJSONBuffer()

	for _, v := range terraformJSONOutput {
		_, err := tfJSON.Write([]byte(v + "\n"))
		if err != nil {
			t.Fatalf("cannot write to tfJSON: %s", err)
		}
	}

	tfJSON.Parse()

	if diff := cmp.Diff(tfJSON.JsonOutput(), terraformJSONOutput); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
