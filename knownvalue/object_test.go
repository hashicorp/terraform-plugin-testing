// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestObjectValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.ObjectExact(map[string]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected map[string]any value for ObjectExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.ObjectExact(map[string]knownvalue.Check{}),
			other: map[string]any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			expectedError: fmt.Errorf("expected map[string]any value for ObjectExact check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected map[string]any value for ObjectExact check, got: float64"),
		},
		"empty": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other:         map[string]any{},
			expectedError: fmt.Errorf("expected 3 attribute(s) for ObjectExact check, got 0 attribute(s): actual value is missing attribute(s): \"one\", \"three\", \"two\""),
		},
		"missing-one-attribute": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one": json.Number("1.23"),
				"two": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("expected 3 attribute(s) for ObjectExact check, got 2 attribute(s): actual value is missing attribute(s): \"three\""),
		},
		"missing-multiple-attributes": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
				"four":  knownvalue.Float64Exact(0.12),
				"five":  knownvalue.Float64Exact(3.45),
			}),
			other: map[string]any{
				"one": json.Number("1.23"),
				"two": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("expected 5 attribute(s) for ObjectExact check, got 2 attribute(s): actual value is missing attribute(s): \"five\", \"four\", \"three\""),
		},
		"extra-one-attribute": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one": knownvalue.Float64Exact(1.23),
				"two": knownvalue.Float64Exact(4.56),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("7.89"),
			},
			expectedError: fmt.Errorf("expected 2 attribute(s) for ObjectExact check, got 3 attribute(s): actual value has extra attribute(s): \"three\""),
		},
		"extra-multiple-attributes": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one": knownvalue.Float64Exact(1.23),
				"two": knownvalue.Float64Exact(4.56),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("7.89"),
				"four":  json.Number("0.12"),
				"five":  json.Number("3.45"),
			},
			expectedError: fmt.Errorf("expected 2 attribute(s) for ObjectExact check, got 5 attribute(s): actual value has extra attribute(s): \"five\", \"four\", \"three\""),
		},
		"not-equal": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("4.56"),
				"three": json.Number("6.54"),
			},
			expectedError: fmt.Errorf("three object attribute: expected value 7.89 for Float64Exact check, got: 6.54"),
		},
		"wrong-order": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"one":   json.Number("1.23"),
				"two":   json.Number("7.89"),
				"three": json.Number("4.56"),
			},
			expectedError: fmt.Errorf("three object attribute: expected value 7.89 for Float64Exact check, got: 4.56"),
		},
		"key-not-found": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
				"one":   knownvalue.Float64Exact(1.23),
				"two":   knownvalue.Float64Exact(4.56),
				"three": knownvalue.Float64Exact(7.89),
			}),
			other: map[string]any{
				"four": json.Number("1.23"),
				"five": json.Number("7.89"),
				"six":  json.Number("4.56"),
			},
			expectedError: fmt.Errorf("missing attribute one for ObjectExact check"),
		},
		"equal": {
			self: knownvalue.ObjectExact(map[string]knownvalue.Check{
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

func TestObjectValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.ObjectExact(map[string]knownvalue.Check{
		"one":   knownvalue.Float64Exact(1.23),
		"two":   knownvalue.Float64Exact(4.56),
		"three": knownvalue.Float64Exact(7.89),
	}).String()

	if diff := cmp.Diff(got, "map[one:1.23 three:7.89 two:4.56]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
