// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestListValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.ListExact([]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected []any value for ListExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.ListExact([]knownvalue.Check{}),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			expectedError: fmt.Errorf("expected []any value for ListExact check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for ListExact check, got: float64"),
		},
		"empty": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("expected 3 elements for ListExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			other: []any{
				int64(123),
				int64(456),
			},
			expectedError: fmt.Errorf("expected 3 elements for ListExact check, got 2 elements"),
		},
		"not-equal": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("456"),
				json.Number("654"),
			},
			expectedError: fmt.Errorf("list element index 2: expected value 789 for Int64Exact check, got: 654"),
		},
		"wrong-order": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("789"),
				json.Number("456"),
			},
			expectedError: fmt.Errorf("list element index 1: expected value 456 for Int64Exact check, got: 789"),
		},
		"equal": {
			self: knownvalue.ListExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Int64Exact(456),
				knownvalue.Int64Exact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("456"),
				json.Number("789"),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.self.CheckValue(testCase.other)

			if diff := cmp.Diff(got, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.ListExact([]knownvalue.Check{
		knownvalue.Int64Exact(123),
		knownvalue.Int64Exact(456),
		knownvalue.Int64Exact(789),
	}).String()

	if diff := cmp.Diff(got, "[123 456 789]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
