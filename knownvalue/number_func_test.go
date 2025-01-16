// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestNumberFunc_CheckValue(t *testing.T) {
	t.Parallel()

	expected, _, err := big.ParseFloat("1.797693134862315797693134862315797693134862315", 10, 512, big.ToNearestEven)
	if err != nil {
		t.Errorf("%s", err)
	}

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"nil": {
			self:          knownvalue.NumberFunc(func(*big.Float) error { return nil }),
			expectedError: fmt.Errorf("expected json.Number value for NumberFunc check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.NumberFunc(func(*big.Float) error { return nil }),
			other:         "wrongtype",
			expectedError: fmt.Errorf("expected json.Number value for NumberFunc check, got: string"),
		},
		"no-digits": {
			self:          knownvalue.NumberFunc(func(*big.Float) error { return nil }),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as big.Float value for NumberFunc check: number has no digits"),
		},
		"failure": {
			self: knownvalue.NumberFunc(func(i *big.Float) error {
				if i.Cmp(expected) != 0 {
					return fmt.Errorf("%s was not %s", i.Text('f', -1), expected.Text('f', -1))
				}
				return nil
			}),
			other:         json.Number("1.667114241575161769818551140818851511176942075"),
			expectedError: fmt.Errorf("1.667114241575161769818551140818851511176942075 was not 1.797693134862315797693134862315797693134862315"),
		},
		"success": {
			self: knownvalue.NumberFunc(func(i *big.Float) error {
				if i.Cmp(expected) != 0 {
					return fmt.Errorf("%s was not %s", i.Text('f', -1), expected.Text('f', -1))
				}
				return nil
			}),
			other: json.Number("1.797693134862315797693134862315797693134862315"),
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

func TestNumberFunc_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NumberFunc(func(*big.Float) error { return nil }).String()

	if diff := cmp.Diff(got, "NumberFunc"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
