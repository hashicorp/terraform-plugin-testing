// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestBoolValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.BoolValue
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is bool"),
		},
		"zero-other": {
			other: false, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.BoolValueExact(false),
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is bool"),
		},
		"wrong-type": {
			self:          knownvalue.BoolValueExact(true),
			other:         1.23,
			expectedError: fmt.Errorf("wrong type: float64, known value type is bool"),
		},
		"not-equal": {
			self:          knownvalue.BoolValueExact(true),
			other:         false,
			expectedError: fmt.Errorf("value: false does not equal expected value: true"),
		},
		"equal": {
			self:  knownvalue.BoolValueExact(true),
			other: true,
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

func TestBoolValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.BoolValueExact(true).String()

	if diff := cmp.Diff(got, "true"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
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
