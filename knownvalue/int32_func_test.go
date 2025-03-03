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

func TestInt32Func_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"nil": {
			self:          knownvalue.Int32Func(func(int32) error { return nil }),
			expectedError: fmt.Errorf("expected json.Number value for Int32Func check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.Int32Func(func(int32) error { return nil }),
			other:         "wrongtype",
			expectedError: fmt.Errorf("expected json.Number value for Int32Func check, got: string"),
		},
		"no-digits": {
			self:          knownvalue.Int32Func(func(int32) error { return nil }),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as int32 value for Int32Func check: strconv.ParseInt: parsing \"str\": invalid syntax"),
		},
		"failure": {
			self: knownvalue.Int32Func(func(i int32) error {
				if i != 1 {
					return fmt.Errorf("%d was not 1", i)
				}
				return nil
			}),
			other:         json.Number("2"),
			expectedError: fmt.Errorf("%d was not 1", 2),
		},
		"success": {
			self: knownvalue.Int32Func(func(i int32) error {
				if i != 1 {
					return fmt.Errorf("%d was not 1", i)
				}
				return nil
			}),
			other: json.Number("1"),
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

func TestInt32Func_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Int32Func(func(int32) error { return nil }).String()

	if diff := cmp.Diff(got, "Int32Func"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
