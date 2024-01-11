// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestStringValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.StringValueExact(""),
			expectedError: fmt.Errorf("expected string value for StringValueExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.StringValueExact(""),
			other: "", // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.StringValueExact("str"),
			expectedError: fmt.Errorf("expected string value for StringValueExact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.StringValueExact("str"),
			other:         1.234,
			expectedError: fmt.Errorf("expected string value for StringValueExact check, got: float64"),
		},
		"not-equal": {
			self:          knownvalue.StringValueExact("str"),
			other:         "rts",
			expectedError: fmt.Errorf("expected value str for StringValueExact check, got: rts"),
		},
		"equal": {
			self:  knownvalue.StringValueExact("str"),
			other: "str",
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

func TestStringValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.StringValueExact("str").String()

	if diff := cmp.Diff(got, "str"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
