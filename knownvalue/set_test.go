// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestSetValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.SetValueExact([]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected []any value for SetValueExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.SetValueExact([]knownvalue.Check{}),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			expectedError: fmt.Errorf("expected []any value for SetValueExact check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for SetValueExact check, got: float64"),
		},
		"empty": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("expected 3 elements for SetValueExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("456"),
			},
			expectedError: fmt.Errorf("expected 3 elements for SetValueExact check, got 2 elements"),
		},
		"not-equal": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("456"),
				json.Number("654"),
			},
			expectedError: fmt.Errorf("missing value 789 for SetValueExact check"),
		},
		"equal-different-order": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("789"),
				json.Number("456"),
			},
		},
		"equal-same-order": {
			self: knownvalue.SetValueExact([]knownvalue.Check{
				knownvalue.Int64ValueExact(123),
				knownvalue.Int64ValueExact(456),
				knownvalue.Int64ValueExact(789),
			}),
			other: []any{
				json.Number("123"),
				json.Number("456"),
				json.Number("789"),
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
