// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestNullValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.NullValue
		other         any
		expectedError error
	}{
		"zero-nil": {},
		"zero-other": {
			other: nil, // checking against the underlying value field zero-value
		},
		"not-nil": {
			self:          knownvalue.NullValueExact(),
			other:         false,
			expectedError: fmt.Errorf("expected value nil for NullValue check, got: bool"),
		},
		"equal": {
			self:  knownvalue.NullValueExact(),
			other: nil,
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

func TestNullValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NullValueExact().String()

	if diff := cmp.Diff(got, "nil"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
