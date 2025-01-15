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

func TestInt64Func_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"nil": {
			self:          knownvalue.Int64Func(func(int64) error { return nil }),
			expectedError: fmt.Errorf("expected json.Number value for Int32Func check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.Int64Func(func(int64) error { return nil }),
			other:         "wrongtype",
			expectedError: fmt.Errorf("expected json.Number value for Int32Func check, got: string"),
		},
		"failure": {
			self: knownvalue.Int64Func(func(i int64) error {
				if i != 1 {
					return fmt.Errorf("%d was not 1", i)
				}
				return nil
			}),
			other:         json.Number("2"),
			expectedError: fmt.Errorf("%d was not 1", 2),
		},
		"success": {
			self: knownvalue.Int64Func(func(i int64) error {
				if i != 1 {
					return fmt.Errorf("%d was not 1", i)
				}
				return nil
			}),
			other: json.Number("1"),
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

func TestInt64Func_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Int64Func(func(int64) error { return nil }).String()

	if diff := cmp.Diff(got, "Int64Func"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
