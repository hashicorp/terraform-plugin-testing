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

func TestListValuePartial_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.ListValuePartial(map[int]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected []any value for ListValuePartial check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.ListValuePartial(map[int]knownvalue.Check{}),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}),
			expectedError: fmt.Errorf("expected []any value for ListValuePartial check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for ListValuePartial check, got: float64"),
		},
		"empty": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("missing element index 0 for ListValuePartial check"),
		},
		"wrong-length": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}),
			other: []any{
				json.Number("1.23"),
				json.Number("4.56"),
			},
			expectedError: fmt.Errorf("missing element index 2 for ListValuePartial check"),
		},
		"not-equal": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}),
			other: []any{
				json.Number("1.23"),
				json.Number("4.56"),
				json.Number("6.54"),
				json.Number("5.46"),
			},
			expectedError: fmt.Errorf("list element 2: expected value 4.56 for Float64ValueExact check, got: 6.54"),
		},
		"wrong-order": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}),
			other: []any{
				json.Number("1.23"),
				json.Number("0.00"),
				json.Number("7.89"),
				json.Number("4.56"),
			},
			expectedError: fmt.Errorf("list element 2: expected value 4.56 for Float64ValueExact check, got: 7.89"),
		},
		"equal": {
			self: knownvalue.ListValuePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
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

func TestListValuePartialPartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.ListValuePartial(map[int]knownvalue.Check{
		0: knownvalue.Float64ValueExact(1.23),
		2: knownvalue.Float64ValueExact(4.56),
		3: knownvalue.Float64ValueExact(7.89),
	}).String()

	if diff := cmp.Diff(got, "[0:1.23 2:4.56 3:7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
