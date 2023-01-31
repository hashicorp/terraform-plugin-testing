package plugintest

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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

func TestTerraformJSONDiagnostics_Errors(t *testing.T) {
	testCases := map[string]struct {
		diags    []tfjson.Diagnostic
		expected TerraformJSONDiagnostics
	}{
		"errors-found": {
			diags: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error 1 summary",
					Detail:   "error 1 detail",
				},
				{
					Severity: tfjson.DiagnosticSeverityWarning,
					Summary:  "warning summary",
					Detail:   "warning detail",
				},
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error 2 summary",
					Detail:   "error 2 detail",
				},
			},
			expected: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error 1 summary",
					Detail:   "error 1 detail",
				},
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error 2 summary",
					Detail:   "error 2 detail",
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			var tfJSONDiagnostics TerraformJSONDiagnostics = testCase.diags

			actual := tfJSONDiagnostics.Errors()

			if diff := cmp.Diff(actual, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTerraformJSONDiagnostics_Warnings(t *testing.T) {
	testCases := map[string]struct {
		diags    []tfjson.Diagnostic
		expected TerraformJSONDiagnostics
	}{
		"warnings-found": {
			diags: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error 1 summary",
					Detail:   "error 1 detail",
				},
				{
					Severity: tfjson.DiagnosticSeverityWarning,
					Summary:  "warning summary",
					Detail:   "warning detail",
				},
				{
					Severity: tfjson.DiagnosticSeverityError,
					Summary:  "error 2 summary",
					Detail:   "error 2 detail",
				},
			},
			expected: []tfjson.Diagnostic{
				{
					Severity: tfjson.DiagnosticSeverityWarning,
					Summary:  "warning summary",
					Detail:   "warning detail",
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			var tfJSONDiagnostics TerraformJSONDiagnostics = testCase.diags

			actual := tfJSONDiagnostics.Warnings()

			if diff := cmp.Diff(actual, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTerraformJSONBuffer_Parse(t *testing.T) {
	t.Parallel()

	tfJSON := NewTerraformJSONBuffer()
	err := tfJSON.ReadFile("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("ReadFile err: %s", err)
	}

	entries, err := fileEntries("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("fileEntries error: %s", err)
	}

	err = tfJSON.Parse()
	if err != nil {
		t.Errorf("parse error: %s", err)
	}

	if diff := cmp.Diff(strings.Split(strings.Trim(tfJSON.rawOutput, "\n"), "\n"), entries); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}

	if diff := cmp.Diff(tfJSON.jsonOutput, entries); diff != "" {
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

	tfJSON, err := NewTerraformJSONBufferFromFile("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("NewTerraformJSONBufferFromFile err: %s", err)
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

	diags, err := tfJSON.Diagnostics()
	if err != nil {
		t.Errorf("Diagnostics error: %s", err)
	}

	if diff := cmp.Diff(diags, tfJSONDiagnostics); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func TestTerraformJSONBuffer_JsonOutput(t *testing.T) {
	t.Parallel()

	tfJSON := NewTerraformJSONBuffer()

	err := tfJSON.ReadFile("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("ReadFile err: %s", err)
	}

	entries, err := fileEntries("../testdata/terraform-json-output.txt")
	if err != nil {
		t.Errorf("fileEntries error: %s", err)
	}

	jsonOutput, err := tfJSON.JsonOutput()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if diff := cmp.Diff(jsonOutput, entries); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func fileEntries(filePath string) ([]string, error) {
	var fileEntries []string

	file, err := os.Open(filePath)
	if err != nil {
		return fileEntries, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileEntries = append(fileEntries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fileEntries, fmt.Errorf("scanner error: %s", err)
	}

	return fileEntries, nil
}
