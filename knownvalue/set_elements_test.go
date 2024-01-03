// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestSetElements_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.SetElements
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("wrong type: <nil>, expected []any"),
		},
		"zero-other": {
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.SetElementsExact(3),
			expectedError: fmt.Errorf("wrong type: <nil>, expected []any"),
		},
		"wrong-type": {
			self:          knownvalue.SetElementsExact(3),
			other:         1.234,
			expectedError: fmt.Errorf("wrong type: float64, expected []any"),
		},
		"empty": {
			self:          knownvalue.SetElementsExact(3),
			other:         []any{},
			expectedError: fmt.Errorf("wrong length: 0, expected 3"),
		},
		"wrong-length": {
			self: knownvalue.SetElementsExact(3),
			other: []any{
				int64(123),
				int64(456),
			},
			expectedError: fmt.Errorf("wrong length: 2, expected 3"),
		},
		"equal": {
			self: knownvalue.SetElementsExact(3),
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

func TestSetElements_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.SetElementsExact(2).String()

	if diff := cmp.Diff(got, "2"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}