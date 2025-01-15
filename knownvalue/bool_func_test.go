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

func TestBoolFunc_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"nil": {
			self:          knownvalue.BoolFunc(func(bool) error { return nil }),
			expectedError: fmt.Errorf("expected bool value for BoolFunc check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.BoolFunc(func(bool) error { return nil }),
			other:         json.Number("1.234"),
			expectedError: fmt.Errorf("expected bool value for BoolFunc check, got: float64"),
		},
		"failure": {
			self: knownvalue.BoolFunc(func(b bool) error {
				if b != true {
					return fmt.Errorf("%t was not true", b)
				}
				return nil
			}),
			other:         false,
			expectedError: fmt.Errorf("%t was not true", false),
		},
		"success": {
			self: knownvalue.BoolFunc(func(b bool) error {
				if b != true {
					return fmt.Errorf("%t was not foo", b)
				}
				return nil
			}),
			other: true,
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

func TestBoolFunc_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.BoolFunc(func(bool) error { return nil }).String()

	if diff := cmp.Diff(got, "BoolFunc"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
