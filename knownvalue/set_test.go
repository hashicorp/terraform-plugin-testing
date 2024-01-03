// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestSetValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.SetValue
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is []Check"),
		},
		"zero-other": {
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is []Check"),
		},
		"wrong-type": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("wrong type: float64, known value type is []Check"),
		},
		"empty": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("wrong length: 0, known value length is 3"),
		},
		"wrong-length": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				int64(123),
				int64(456),
			},
			expectedError: fmt.Errorf("wrong length: 2, known value length is 3"),
		},
		"not-equal": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				int64(123),
				int64(456),
				int64(654),
			},
			expectedError: fmt.Errorf("expected value not found: 789"),
		},
		"equal-different-order": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				int64(123),
				int64(789),
				int64(456),
			},
		},
		"equal-same-order": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				int64(123),
				int64(456),
				int64(789),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.self.CheckValue(testCase.other)

			if diff := cmp.Diff(got, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSetValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.SetValueExact([]knownvalue.Check{
		knownvalue.Int64ValueExact(123),
		knownvalue.Int64ValueExact(456),
		knownvalue.Int64ValueExact(789),
	}).String()

	if diff := cmp.Diff(got, "[123 456 789]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
