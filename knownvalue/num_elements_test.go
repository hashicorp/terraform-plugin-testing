// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestNumElements_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		other    any
		expected bool
	}{
		"nil": {},
		"wrong-type": {
			other: 1.23,
		},
		"empty-map": {
			other: map[string]any{},
		},
		"empty-slice": {
			other: []any{},
		},
		"map-different-len": {
			other: map[string]any{
				"one": 1.23,
				"two": 4.56,
			},
		},
		"slice-different-len": {
			other: []any{1.23, 4.56},
		},
		"equal-map": {
			other: map[string]any{
				"one":   1.23,
				"two":   4.56,
				"three": 7.89,
			},
			expected: true,
		},
		"equal-slice": {
			other:    []any{1.23, 4.56, 7.89},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.NumElementsExact(3).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestNumElements_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NumElementsExact(2).String()

	if diff := cmp.Diff(got, "2"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
