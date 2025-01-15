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
		"failure": {
			self: knownvalue.NumberFunc(func(i *big.Float) error {
				if i != big.NewFloat(1.667114241575161769818551140818851511176942075) {
					return fmt.Errorf("%f was not 1.667114241575161769818551140818851511176942075", i)
				}
				return nil
			}),
			other:         json.Number("1.797693134862315797693134862315797693134862315"),
			expectedError: fmt.Errorf("%f was not 1.667114241575161769818551140818851511176942075", 1.797693134862315797693134862315797693134862315),
		},
		"success": {
			self: knownvalue.NumberFunc(func(i *big.Float) error {
				if i != big.NewFloat(1.667114241575161769818551140818851511176942075) {
					return fmt.Errorf("%f was not 1.667114241575161769818551140818851511176942075", i)
				}
				return nil
			}),
			other: json.Number("1.667114241575161769818551140818851511176942075"),
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
