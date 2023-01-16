package resource

import (
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

var stdoutJSON = []string{
	`{"@level":"info","@message":"Terraform 1.3.2","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.232751Z","terraform":"1.3.2","type":"version","ui":"1.0"}`,
	`{"@level":"warn","@message":"Warning: Empty or non-existent state","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.250725Z","diagnostic":{"severity":"warning","summary":"Empty or non-existent state","detail":"There are currently no remote objects tracked in the state, so there is nothing to refresh."},"type":"diagnostic"}`,
	`{"@level":"info","@message":"Outputs: 0","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.251120Z","outputs":{},"type":"outputs"}`,
	`{"@level":"info","@message":"random_password.test: Plan to create","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.282021Z","change":{"resource":{"addr":"random_password.test","module":"","resource":"random_password.test","implied_provider":"random","resource_type":"random_password","resource_name":"test","resource_key":null},"action":"create"},"type":"planned_change"}`,
	`{"@level":"info","@message":"Plan: 1 to add, 0 to change, 0 to destroy.","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.282054Z","changes":{"add":1,"change":0,"remove":0,"operation":"plan"},"type":"change_summary"}`,
	`{"format_version":"1.0"}`,
	`{"@level":"error","@message":"Error: error diagnostic - summary","@module":"terraform.ui","@timestamp":"2023-01-16T17:02:14.626572Z","diagnostic":{"severity":"error","summary":"error diagnostic - summary","detail":""},"type":"diagnostic"}`,
}

func TestStdout_GetJSONOutputStr(t *testing.T) {
	t.Parallel()

	stdout := NewStdout()

	for _, v := range stdoutJSON {
		stdout.Write([]byte(v + "\n"))
	}

	jsonOutputStr := stdout.GetJSONOutputStr()

	if diff := cmp.Diff(jsonOutputStr, strings.Join(stdoutJSON, "\n")); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

func TestStdout_DiagnosticFound(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		regex        *regexp.Regexp
		severity     tfjson.DiagnosticSeverity
		expected     bool
		expectedJSON []string
	}{
		"found": {
			regex:        regexp.MustCompile(`.*error diagnostic - summary`),
			severity:     tfjson.DiagnosticSeverityError,
			expected:     true,
			expectedJSON: stdoutJSON,
		},
		"not-found": {
			regex:        regexp.MustCompile(`.*warning diagnostic - summary`),
			severity:     tfjson.DiagnosticSeverityError,
			expected:     false,
			expectedJSON: stdoutJSON,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			stdout := NewStdout()

			for _, v := range stdoutJSON {
				_, err := stdout.Write([]byte(v + "\n"))
				if err != nil {
					t.Errorf("error writing to stdout: %s", err)
				}
			}

			isFound, output := stdout.DiagnosticFound(testCase.regex, testCase.severity)

			if diff := cmp.Diff(isFound, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
			if diff := cmp.Diff(output, testCase.expectedJSON); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}

}
