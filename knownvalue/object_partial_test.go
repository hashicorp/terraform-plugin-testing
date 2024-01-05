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

func TestObjectValuePartial_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.ObjectValuePartial
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("expected map[string]any value for ObjectValuePartial check, got: <nil>"),
		},
		"zero-other": {
			other: map[string]any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			expectedError: fmt.Errorf("expected map[string]any value for ObjectValuePartial check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected map[string]any value for ObjectValuePartial check, got: float64"),
		},
		"empty": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other:         map[string]any{},
			expectedError: fmt.Errorf("missing attribute one for ObjectValuePartial check"),
		},
		"wrong-length": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one": json.Number("1.23"),
				"two": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("missing attribute three for ObjectValuePartial check"),
		},
		"not-equal": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("6.54"),
			},
			expectedError: fmt.Errorf("three object attribute: expected value 7.89 for Float64Value check, got: 6.54"),
		},
		"wrong-order": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("7.89"),
				"three": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("three object attribute: expected value 7.89 for Float64Value check, got: 4.56"),
		},
		"key-not-found": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"two":   knownvalue.Float64ValueExact(4.56),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"four": json.Number("1.23"),
				"five": json.Number("7.89"),
				"six":  json.Number("4.56"),
			},
			expectedError: fmt.Errorf("missing attribute one for ObjectValuePartial check"),
		},
		"equal": {
			self: knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
				"one":   knownvalue.Float64ValueExact(1.23),
				"three": knownvalue.Float64ValueExact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("7.89"),
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

func TestObjectValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.ObjectValuePartialMatch(map[string]knownvalue.Check{
		"one":   knownvalue.Float64ValueExact(1.23),
		"three": knownvalue.Float64ValueExact(7.89),
	}).String()

	if diff := cmp.Diff(got, "map[one:1.23 three:7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
