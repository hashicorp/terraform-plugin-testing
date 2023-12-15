// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestObjectValue_Equal(t *testing.T) {
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
		"different-len": {
			other: map[string]any{
				"one": 1.23,
				"two": 4.56,
			},
		},
		"not-equal": {
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

			got := knownvalue.NewObjectValue(map[string]knownvalue.KnownValue{
				"one":   knownvalue.NewFloat64Value(1.23),
				"two":   knownvalue.NewFloat64Value(4.56),
				"three": knownvalue.NewFloat64Value(7.89),
			}).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestObjectValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NewObjectValue(map[string]knownvalue.KnownValue{
		"one":   knownvalue.NewFloat64Value(1.23),
		"two":   knownvalue.NewFloat64Value(4.56),
		"three": knownvalue.NewFloat64Value(7.89),
	}).String()

	if diff := cmp.Diff(got, "map[one:1.23 three:7.89 two:4.56]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
