// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compare_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

func TestValuesSame_CompareValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            []any
		expectedError error
	}{
		"nil": {},
		"single-value": {
			in: []any{"str"},
		},
		"matching-values": {
			in: []any{"str", "str", "str"},
		},
		"non-matching-values-string": {
			in:            []any{"str", "str", "other_str"},
			expectedError: fmt.Errorf("expected values to be the same, but they differ: str != other_str"),
		},
		"non-matching-values-slice": {
			in: []any{
				[]any{"str"},
				[]any{"str"},
				[]any{"other_str"},
			},
			expectedError: fmt.Errorf("expected values to be the same, but they differ: [str] != [other_str]"),
		},
		"non-matching-values-map": {
			in: []any{
				map[string]any{"a": "str"},
				map[string]any{"a": "str"},
				map[string]any{"a": "other_str"},
			},
			expectedError: fmt.Errorf("expected values to be the same, but they differ: map[a:str] != map[a:other_str]"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := compare.ValuesSame().CompareValues(testCase.in...)

			if diff := cmp.Diff(err, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

// equateErrorMessage reports errors to be equal if both are nil
// or both have the same message.
var equateErrorMessage = cmp.Comparer(func(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	return x.Error() == y.Error()
})
