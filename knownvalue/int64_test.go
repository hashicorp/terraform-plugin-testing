// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestInt64Value_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Int64Value
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is int64"),
		},
		"zero-other": {
			other: int64(0), // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.Int64ValueExact(1234),
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is int64"),
		},
		"wrong-type": {
			self:          knownvalue.Int64ValueExact(1234),
			other:         1.234,
			expectedError: fmt.Errorf("wrong type: float64, known value type is int64"),
		},
		"not-equal": {
			self:          knownvalue.Int64ValueExact(1234),
			other:         int64(4321),
			expectedError: fmt.Errorf("value: 4321 does not equal expected value: 1234"),
		},
		"equal": {
			self:  knownvalue.Int64ValueExact(1234),
			other: int64(1234),
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

func TestInt64Value_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Int64ValueExact(1234567890123).String()

	if diff := cmp.Diff(got, "1234567890123"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
