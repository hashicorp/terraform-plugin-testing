// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package compare_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

func TestValuesDiffer_CompareValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            []any
		expectedError error
	}{
		"nil": {},
		"single-value": {
			in: []any{"str"},
		},
		"non-matching-sequential-values": {
			in: []any{"str", "other_str", "str"},
		},
		"matching-values-string": {
			in:            []any{"str", "other_str", "other_str"},
			expectedError: fmt.Errorf("expected values to differ, but they are the same: other_str == other_str"),
		},
		"matching-values-slice": {
			in: []any{
				[]any{"other_str"},
				[]any{"other_str"},
			},
			expectedError: fmt.Errorf("expected values to differ, but they are the same: [other_str] == [other_str]"),
		},
		"matching-values-map": {
			in: []any{
				map[string]any{"a": "other_str"},
				map[string]any{"a": "other_str"},
			},
			expectedError: fmt.Errorf("expected values to differ, but they are the same: map[a:other_str] == map[a:other_str]"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := compare.ValuesDiffer().CompareValues(testCase.in...)

			if diff := cmp.Diff(err, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
