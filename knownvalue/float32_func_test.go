// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestFloat32Func_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"nil": {
			self:          knownvalue.Float32Func(func(float32) error { return nil }),
			expectedError: fmt.Errorf("expected json.Number value for Float32Func check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.Float32Func(func(float32) error { return nil }),
			other:         "wrongtype",
			expectedError: fmt.Errorf("expected json.Number value for Float32Func check, got: string"),
		},
		"no-digits": {
			self:          knownvalue.Float32Func(func(float32) error { return nil }),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as float32 value for Float32Func check: strconv.ParseFloat: parsing \"str\": invalid syntax"),
		},
		"failure": {
			self: knownvalue.Float32Func(func(f float32) error {
				if f != 1.1 {
					return fmt.Errorf("%f was not 1.1", f)
				}
				return nil
			}),
			other:         json.Number("1.2"),
			expectedError: fmt.Errorf("%f was not 1.1", 1.2),
		},
		"success": {
			self: knownvalue.Float32Func(func(f float32) error {
				if f != 1.1 {
					return fmt.Errorf("%f was not 1.1", f)
				}
				return nil
			}),
			other: json.Number("1.1"),
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

func TestFloat32Func_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Float32Func(func(float32) error { return nil }).String()

	if diff := cmp.Diff(got, "Float32Func"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
