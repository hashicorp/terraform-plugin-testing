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

func TestFloat64Value_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.Float64Exact(0),
			expectedError: fmt.Errorf("expected json.Number value for Float64Exact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.Float64Exact(0),
			other: json.Number("0.0"), // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.Float64Exact(1.234),
			expectedError: fmt.Errorf("expected json.Number value for Float64Exact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.Float64Exact(1.234),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as float64 value for Float64Exact check: strconv.ParseFloat: parsing \"str\": invalid syntax"),
		},
		"not-equal": {
			self:          knownvalue.Float64Exact(1.234),
			other:         json.Number("4.321"),
			expectedError: fmt.Errorf("expected value 1.234 for Float64Exact check, got: 4.321"),
		},
		"equal": {
			self:  knownvalue.Float64Exact(1.234),
			other: json.Number("1.234"),
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

func TestFloat64Value_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Float64Exact(1.234567890123e+09).String()

	if diff := cmp.Diff(got, "1234567890.123"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
