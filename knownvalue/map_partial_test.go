// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestMapValuePartial_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		other    any
		expected bool
	}{
		"nil": {},
		"wrong-type": {
			other: 1.23,
		},
		"empty": {
			other: map[string]any{},
		},
		"not-equal-different-len": {
			other: map[string]any{
				"one": 1.23,
				"two": 4.56,
			},
		},
		"not-equal-same-len": {
			other: map[string]any{
				"one":   1.23,
				"two":   4.56,
				"three": 6.54,
			},
		},
		"equal": {
			other: map[string]any{
				"one":   1.23,
				"two":   4.56,
				"three": 7.89,
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.MapValuePartialMatch(map[string]knownvalue.KnownValue{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMapValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.MapValuePartialMatch(map[string]knownvalue.KnownValue{
		"one":   knownvalue.Float64ValueExact(1.23),
		"three": knownvalue.Float64ValueExact(7.89),
	}).String()

	if diff := cmp.Diff(got, "map[one:1.23 three:7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
