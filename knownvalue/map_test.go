// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestMapValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.MapExact(map[string]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected map[string]any value for MapExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.MapExact(map[string]knownvalue.Check{}),
			other: map[string]any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			expectedError: fmt.Errorf("expected map[string]any value for MapExact check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected map[string]any value for MapExact check, got: float64"),
		},
		"empty": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other:         map[string]any{},
			expectedError: fmt.Errorf("expected 3 elements for MapExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one": json.Number("1.23"),
				"two": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("expected 3 elements for MapExact check, got 2 elements"),
		},
		"not-equal": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("6.54"),
			},
			expectedError: fmt.Errorf("three map element: expected value 7.89 for Float64Exact check, got: 6.54"),
		},
		"wrong-order": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("7.89"),
				"three": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("three map element: expected value 7.89 for Float64Exact check, got: 4.56"),
		},
		"key-not-found": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"four": json.Number("1.23"),
				"five": json.Number("7.89"),
				"six":  json.Number("4.56"),
			},
			expectedError: fmt.Errorf("missing element one for MapExact check"),
		},
		"equal": {
			self: knownvalue.MapExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("7.89"),
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

func TestMapValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.MapExact(map[string]knownvalue.Check{
		"one":   knownvalue.Float64Exact(1.23),
		"two":   knownvalue.Float64Exact(4.56),
		"three": knownvalue.Float64Exact(7.89),
	}).String()

	if diff := cmp.Diff(got, "map[one:1.23 three:7.89 two:4.56]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
