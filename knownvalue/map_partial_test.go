// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestMapValuePartial_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.MapValuePartial
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("expected map[string]any value for MapValuePartial check, got: <nil>"),
		},
		"zero-other": {
			other: map[string]any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			expectedError: fmt.Errorf("expected map[string]any value for MapValuePartial check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected map[string]any value for MapValuePartial check, got: float64"),
		},
		"empty": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other:         map[string]any{},
			expectedError: fmt.Errorf("missing element one for MapValuePartial check"),
		},
		"wrong-length": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one": 1.23,
				"two": 4.56,
			},
			expectedError: fmt.Errorf("missing element three for MapValuePartial check"),
		},
		"not-equal": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one":   1.23,
				"two":   4.56,
				"three": 6.54,
			},
			expectedError: fmt.Errorf("three map element: expected value 7.89 for Float64Value check, got: 6.54"),
		},
		"wrong-order": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one":   1.23,
				"two":   7.89,
				"three": 4.56,
			},
			expectedError: fmt.Errorf("three map element: expected value 7.89 for Float64Value check, got: 4.56"),
		},
		"key-not-found": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"two":   knownvalue.Float64ValueExact(4.56),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"four": 1.23,
				"five": 7.89,
				"six":  4.56,
			},
			expectedError: fmt.Errorf("missing element one for MapValuePartial check"),
		},
		"equal": {
			self: knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one":   1.23,
				"two":   4.56,
				"three": 7.89,
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

func TestMapValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.MapValuePartialMatch(map[string]knownvalue.Check{
		"one":   knownvalue.Float64ValueExact(1.23),
		"three": knownvalue.Float64ValueExact(7.89),
	}).String()

	if diff := cmp.Diff(got, "map[one:1.23 three:7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
