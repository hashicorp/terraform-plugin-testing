// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestFloat64Value_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Float64Value
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is float64"),
		},
		"zero-other": {
			other: 0.0, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.Float64ValueExact(1.234),
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is float64"),
		},
		"wrong-type": {
			self:          knownvalue.Float64ValueExact(1.234),
			other:         int64(1234),
			expectedError: fmt.Errorf("wrong type: int64, known value type is float64"),
		},
		"not-equal": {
			self:          knownvalue.Float64ValueExact(1.234),
			other:         4.321,
			expectedError: fmt.Errorf("value: 4.321 does not equal expected value: 1.234"),
		},
		"equal": {
			self:  knownvalue.Float64ValueExact(1.234),
			other: 1.234,
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

func TestFloat64Value_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Float64ValueExact(1.234567890123e+09).String()

	if diff := cmp.Diff(got, "1234567890.123"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
