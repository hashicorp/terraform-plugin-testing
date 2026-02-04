// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestSetValuePartial_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.SetPartial([]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected []any value for SetPartial check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.SetPartial([]knownvalue.Check{}),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.SetPartial([]knownvalue.Check{
				knownvalue.Float64Exact(1.23),
				knownvalue.Float64Exact(4.56),
				knownvalue.Float64Exact(7.89),
			}),
			expectedError: fmt.Errorf("expected []any value for SetPartial check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.SetPartial([]knownvalue.Check{
				knownvalue.Float64Exact(1.23),
				knownvalue.Float64Exact(4.56),
				knownvalue.Float64Exact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for SetPartial check, got: float64"),
		},
		"equal-empty": {
			self: knownvalue.SetPartial([]knownvalue.Check{
				knownvalue.Float64Exact(1.23),
				knownvalue.Float64Exact(4.56),
				knownvalue.Float64Exact(7.89),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("missing value 1.23 for SetPartial check"),
		},
		"not-equal": {
			self: knownvalue.SetPartial([]knownvalue.Check{
				knownvalue.Float64Exact(1.23),
				knownvalue.Float64Exact(4.56),
				knownvalue.Float64Exact(7.89),
			}),
			other: []any{
				json.Number("1.23"),
				json.Number("4.56"),
				json.Number("6.54"),
				json.Number("5.46"),
			},
			expectedError: fmt.Errorf("missing value 7.89 for SetPartial check"),
		},
		"equal-different-order": {
			self: knownvalue.SetPartial([]knownvalue.Check{
				knownvalue.Float64Exact(1.23),
				knownvalue.Float64Exact(4.56),
				knownvalue.Float64Exact(7.89),
			}),
			other: []any{
				json.Number("1.23"),
				json.Number("0.00"),
				json.Number("7.89"),
				json.Number("4.56"),
			},
		},
		"equal-same-order": {
			self: knownvalue.SetPartial([]knownvalue.Check{
				knownvalue.Float64Exact(1.23),
				knownvalue.Float64Exact(4.56),
				knownvalue.Float64Exact(7.89),
			}),
			other: []any{
				json.Number("1.23"),
				json.Number("0.00"),
				json.Number("4.56"),
				json.Number("7.89"),
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

func TestSetValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.SetPartial([]knownvalue.Check{
		knownvalue.Float64Exact(1.23),
		knownvalue.Float64Exact(4.56),
		knownvalue.Float64Exact(7.89),
	}).String()

	if diff := cmp.Diff(got, "[1.23 4.56 7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
