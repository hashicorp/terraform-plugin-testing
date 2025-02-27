// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestTupleSizeExact_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.TupleSizeExact(0),
			expectedError: fmt.Errorf("expected []any value for TupleSizeExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.TupleSizeExact(0),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.TupleSizeExact(3),
			expectedError: fmt.Errorf("expected []any value for TupleSizeExact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.TupleSizeExact(3),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for TupleSizeExact check, got: float64"),
		},
		"empty": {
			self:          knownvalue.TupleSizeExact(3),
			other:         []any{},
			expectedError: fmt.Errorf("expected 3 elements for TupleSizeExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.TupleSizeExact(4),
			other: []any{
				int64(123),
				"hello",
				true,
			},
			expectedError: fmt.Errorf("expected 4 elements for TupleSizeExact check, got 3 elements"),
		},
		"equal": {
			self: knownvalue.TupleSizeExact(3),
			other: []any{
				int64(123),
				"hello",
				true,
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

func TestTupleSizeExact_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.TupleSizeExact(2).String()

	if diff := cmp.Diff(got, "2"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
