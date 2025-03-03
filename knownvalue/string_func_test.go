// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestStringFunc_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"nil": {
			self:          knownvalue.StringFunc(func(string) error { return nil }),
			expectedError: fmt.Errorf("expected string value for StringFunc check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.StringFunc(func(string) error { return nil }),
			other:         1.234,
			expectedError: fmt.Errorf("expected string value for StringFunc check, got: float64"),
		},
		"failure": {
			self: knownvalue.StringFunc(func(s string) error {
				if s != "foo" {
					return fmt.Errorf("%s was not foo", s)
				}
				return nil
			}),
			other:         "bar",
			expectedError: fmt.Errorf("bar was not foo"),
		},
		"success": {
			self: knownvalue.StringFunc(func(s string) error {
				if s != "foo" {
					return fmt.Errorf("%s was not foo", s)
				}
				return nil
			}),
			other: "foo",
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

func TestStringFunc_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.StringFunc(func(string) error { return nil }).String()

	if diff := cmp.Diff(got, "StringFunc"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
