// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestMapElements_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.MapSizeExact(0),
			expectedError: fmt.Errorf("expected map[string]any value for MapSizeExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.MapSizeExact(0),
			other: map[string]any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.MapSizeExact(3),
			expectedError: fmt.Errorf("expected map[string]any value for MapSizeExact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.MapSizeExact(3),
			other:         1.234,
			expectedError: fmt.Errorf("expected map[string]any value for MapSizeExact check, got: float64"),
		},
		"empty": {
			self:          knownvalue.MapSizeExact(3),
			other:         map[string]any{},
			expectedError: fmt.Errorf("expected 3 elements for MapSizeExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.MapSizeExact(3),
			other: map[string]any{
				"one": int64(123),
				"two": int64(456),
			},
			expectedError: fmt.Errorf("expected 3 elements for MapSizeExact check, got 2 elements"),
		},
		"equal": {
			self: knownvalue.MapSizeExact(3),
			other: map[string]any{
				"one":   int64(123),
				"two":   int64(456),
				"three": int64(789),
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

func TestMapElements_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.MapSizeExact(2).String()

	if diff := cmp.Diff(got, "2"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
