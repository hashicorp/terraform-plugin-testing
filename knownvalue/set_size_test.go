// Copyright IBM Corp. 2014, 2025
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
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.SetSizeExact(0),
			expectedError: fmt.Errorf("expected []any value for SetElementExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.SetSizeExact(0),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.SetSizeExact(3),
			expectedError: fmt.Errorf("expected []any value for SetElementExact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.SetSizeExact(3),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for SetElementExact check, got: float64"),
		},
		"empty": {
			self:          knownvalue.SetSizeExact(3),
			other:         []any{},
			expectedError: fmt.Errorf("expected 3 elements for SetElementExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.SetSizeExact(3),
			other: []any{
				int64(123),
				int64(456),
			},
			expectedError: fmt.Errorf("expected 3 elements for SetElementExact check, got 2 elements"),
		},
		"equal": {
			self: knownvalue.SetSizeExact(3),
			other: []any{
				int64(123),
				int64(456),
				int64(789),
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

func TestSetElements_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.SetSizeExact(2).String()

	if diff := cmp.Diff(got, "2"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
