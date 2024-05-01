// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

func TestValuesSameAny_CompareValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            []any
		expectedError error
	}{
		"nil": {},
		"single-value": {
			in: []any{"str"},
		},
		"non-sequential-matching-values": {
			in: []any{"str", "other_str", "str"},
		},
		"non-matching-values-string": {
			in:            []any{"str", "other_str", "another_str"},
			expectedError: fmt.Errorf("expected at least two values to be the same, but all values differ"),
		},
		"non-matching-values-slice": {
			in: []any{
				[]any{"str"},
				[]any{"other_str"},
				[]any{"another_str"},
			},
			expectedError: fmt.Errorf("expected at least two values to be the same, but all values differ"),
		},
		"non-matching-values-map": {
			in: []any{
				map[string]any{"a": "str"},
				map[string]any{"a": "other_str"},
				map[string]any{"a": "another_str"},
			},
			expectedError: fmt.Errorf("expected at least two values to be the same, but all values differ"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := compare.ValuesSameAny().CompareValues(testCase.in...)

			if diff := cmp.Diff(err, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
