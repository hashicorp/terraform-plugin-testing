// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestStringRegularExpression_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.StringRegularExpression(regexp.MustCompile("")),
			expectedError: fmt.Errorf("expected string value for StringRegularExpression check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.StringRegularExpression(regexp.MustCompile("")),
			other: "", // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.StringRegularExpression(regexp.MustCompile("str")),
			expectedError: fmt.Errorf("expected string value for StringRegularExpression check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.StringRegularExpression(regexp.MustCompile("str")),
			other:         1.234,
			expectedError: fmt.Errorf("expected string value for StringRegularExpression check, got: float64"),
		},
		"not-equal": {
			self:          knownvalue.StringRegularExpression(regexp.MustCompile("str")),
			other:         "rts",
			expectedError: fmt.Errorf("expected regex match str for StringRegularExpression check, got: rts"),
		},
		"equal": {
			self:  knownvalue.StringRegularExpression(regexp.MustCompile("str")),
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

func TestStringRegularExpression_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.StringRegularExpression(regexp.MustCompile("^str[0-9a-z]")).String()

	if diff := cmp.Diff(got, "^str[0-9a-z]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
