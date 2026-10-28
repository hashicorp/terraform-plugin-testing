// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilterWorkspaceNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "valid workspace names only",
			input:    []string{"default", "foo", "bar"},
			expected: []string{"default", "foo", "bar"},
		},
		{
			name:  "warning diagnostics from a state store method",
			input: []string{"This warning is intentionally returned by the test state store.", "Warning: test warning from the GetStates method", "default"},
			expected: []string{
				"default",
			},
		},
		{
			name:     "warning about no workspaces being present",
			input:    []string{`The "default" workspace is selected in your working directory. You can create`, "Warning: Terraform cannot find any existing workspaces.", `new" subcommand or by including the "-or-create" flag with the "terraform`, `this workspace by running "terraform init", by using the "terraform workspace`, `"workspace select" subcommand.`},
			expected: []string{},
		},
		{
			name:     "workspace names containing URL-escape characters are invalid",
			input:    []string{"foo%20bar", "default"},
			expected: []string{"default"},
		},
		{
			name:     "valid workspace names containing characters that URL-escaping allows",
			input:    []string{"foo-bar", "foo.bar", "foo_bar", "foo123", "FOO"},
			expected: []string{"foo-bar", "foo.bar", "foo_bar", "foo123", "FOO"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := filterWorkspaceNames(context.Background(), testCase.input)

			if diff := cmp.Diff(testCase.expected, actual); diff != "" {
				t.Errorf("unexpected difference (-want +got):\n%s", diff)
			}
		})
	}
}
