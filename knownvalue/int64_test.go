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

func TestInt64Value_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Int64Value
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("expected json.Number value for Int64Value check, got: <nil>"),
		},
		"zero-other": {
			other: json.Number("0"), // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.Int64ValueExact(1234),
			expectedError: fmt.Errorf("expected json.Number value for Int64Value check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.Int64ValueExact(1234),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as int64 value for Int64Value check: strconv.ParseInt: parsing \"str\": invalid syntax"),
		},
		"not-equal": {
			self:          knownvalue.Int64ValueExact(1234),
			other:         json.Number("4321"),
			expectedError: fmt.Errorf("expected value 1234 for Int64Value check, got: 4321"),
		},
		"equal": {
			self:  knownvalue.Int64ValueExact(1234),
			other: json.Number("1234"),
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

func TestInt64Value_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Int64ValueExact(1234567890123).String()

	if diff := cmp.Diff(got, "1234567890123"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
