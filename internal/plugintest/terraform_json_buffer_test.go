package plugintest

import (
	"bufio"
	"os"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

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

	file, err := os.Open("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("cannot read file: %s", err)
	}
	defer file.Close()

	var fileEntries []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()

		_, err := tfJSON.Write([]byte(txt + "\n"))
		if err != nil {
			t.Errorf("cannot write to tfJSON: %s", err)
		}

		fileEntries = append(fileEntries, txt)
	}

	if err := scanner.Err(); err != nil {
		t.Errorf("scanner error: %s", err)
	}

	tfJSON.Parse()

	if diff := cmp.Diff(tfJSON.jsonOutput, fileEntries); diff != "" {
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

	file, err := os.Open("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("cannot read file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()

		_, err := tfJSON.Write([]byte(txt + "\n"))
		if err != nil {
			t.Errorf("cannot write to tfJSON: %s", err)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Errorf("scanner error: %s", err)
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

	file, err := os.Open("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("cannot read file: %s", err)
	}
	defer file.Close()

	var fileEntries []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()

		_, err := tfJSON.Write([]byte(txt + "\n"))
		if err != nil {
			t.Errorf("cannot write to tfJSON: %s", err)
		}

		fileEntries = append(fileEntries, txt)
	}

	if err := scanner.Err(); err != nil {
		t.Errorf("scanner error: %s", err)
	}

	tfJSON.Parse()

	if diff := cmp.Diff(tfJSON.JsonOutput(), fileEntries); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}