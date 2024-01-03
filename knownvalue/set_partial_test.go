// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestSetValuePartial_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.SetValuePartial
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("expected []any value for SetValuePartial check, got: <nil>"),
		},
		"zero-other": {
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.SetValuePartialMatch([]knownvalue.Check{
				knownvalue.Float64ValueExact(1.23),
				knownvalue.Float64ValueExact(4.56),
				knownvalue.Float64ValueExact(7.89),
			}),
			expectedError: fmt.Errorf("expected []any value for SetValuePartial check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.SetValuePartialMatch([]knownvalue.Check{
				knownvalue.Float64ValueExact(1.23),
				knownvalue.Float64ValueExact(4.56),
				knownvalue.Float64ValueExact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for SetValuePartial check, got: float64"),
		},
		"equal-empty": {
			self: knownvalue.SetValuePartialMatch([]knownvalue.Check{
				knownvalue.Float64ValueExact(1.23),
				knownvalue.Float64ValueExact(4.56),
				knownvalue.Float64ValueExact(7.89),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("missing value 1.23 for SetValuePartial check"),
		},
		"not-equal": {
			self: knownvalue.SetValuePartialMatch([]knownvalue.Check{
				knownvalue.Float64ValueExact(1.23),
				knownvalue.Float64ValueExact(4.56),
				knownvalue.Float64ValueExact(7.89),
			}),
			other:         []any{1.23, 4.56, 6.54, 5.46},
			expectedError: fmt.Errorf("missing value 7.89 for SetValuePartial check"),
		},
		"equal-different-order": {
			self: knownvalue.SetValuePartialMatch([]knownvalue.Check{
				knownvalue.Float64ValueExact(1.23),
				knownvalue.Float64ValueExact(4.56),
				knownvalue.Float64ValueExact(7.89),
			}),
			other: []any{1.23, 0.00, 7.89, 4.56},
		},
		"equal-same-order": {
			self: knownvalue.SetValuePartialMatch([]knownvalue.Check{
				knownvalue.Float64ValueExact(1.23),
				knownvalue.Float64ValueExact(4.56),
				knownvalue.Float64ValueExact(7.89),
			}),
			other: []any{1.23, 0.00, 4.56, 7.89},
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

func TestSetValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.SetValuePartialMatch([]knownvalue.Check{
		knownvalue.Float64ValueExact(1.23),
		knownvalue.Float64ValueExact(4.56),
		knownvalue.Float64ValueExact(7.89),
	}).String()

	if diff := cmp.Diff(got, "[1.23 4.56 7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
