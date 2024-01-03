// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestNumberValue_Equal(t *testing.T) {
	t.Parallel()

	bigFloat, _, err := big.ParseFloat("1.797693134862315797693134862315797693134862315", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	otherBigFloat, _, err := big.ParseFloat("1.797693134862315797693134862315797693134862314", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	testCases := map[string]struct {
		self          knownvalue.NumberValue
		other         any
		expectedError error
	}{
		"zero-nil": {
			expectedError: fmt.Errorf("known value type is nil"),
		},
		"zero-other": {
			other:         otherBigFloat, // checking against the underlying value field zero-value
			expectedError: fmt.Errorf("known value type is nil"),
		},
		"nil": {
			self:          knownvalue.NumberValueExact(bigFloat),
			expectedError: fmt.Errorf("wrong type: <nil>, known value type is *big.Float"),
		},
		"wrong-type": {
			self:          knownvalue.NumberValueExact(bigFloat),
			other:         1.234,
			expectedError: fmt.Errorf("wrong type: float64, known value type is *big.Float"),
		},
		"not-equal": {
			self:          knownvalue.NumberValueExact(bigFloat),
			other:         otherBigFloat,
			expectedError: fmt.Errorf("value: 1.797693134862315797693134862315797693134862314 does not equal expected value: 1.797693134862315797693134862315797693134862315"),
		},
		"equal": {
			self:  knownvalue.NumberValueExact(bigFloat),
			other: bigFloat,
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

func TestNumberValue_String(t *testing.T) {
	t.Parallel()

	bigFloat, _, err := big.ParseFloat("1.797693134862315797693134862315797693134862315", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	got := knownvalue.NumberValueExact(bigFloat).String()

	if diff := cmp.Diff(got, "1.797693134862315797693134862315797693134862315"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
