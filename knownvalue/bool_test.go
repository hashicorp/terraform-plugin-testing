// Copyright IBM Corp. 2014, 2025
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
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.Bool(false),
			expectedError: fmt.Errorf("expected bool value for Bool check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.Bool(false),
			other: false, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.Bool(false),
			expectedError: fmt.Errorf("expected bool value for Bool check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.Bool(true),
			other:         1.23,
			expectedError: fmt.Errorf("expected bool value for Bool check, got: float64"),
		},
		"not-equal": {
			self:          knownvalue.Bool(true),
			other:         false,
			expectedError: fmt.Errorf("expected value true for Bool check, got: false"),
		},
		"equal": {
			self:  knownvalue.Bool(true),
			other: true,
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

func TestBoolValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Bool(true).String()

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
