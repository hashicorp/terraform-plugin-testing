// Copyright IBM Corp. 2014, 2025
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

func TestNumberValue_Equal(t *testing.T) {
	t.Parallel()

	bigFloat, _, err := big.ParseFloat("1.797693134862315797693134862315797693134862315", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.NumberExact(nil),
			expectedError: fmt.Errorf("value in NumberExact check is nil"),
		},
		"zero-other": {
			self:          knownvalue.NumberExact(nil),
			other:         json.Number("1.797693134862315797693134862315797693134862314"), // checking against the underlying value field zero-value
			expectedError: fmt.Errorf("value in NumberExact check is nil"),
		},
		"nil": {
			self:          knownvalue.NumberExact(bigFloat),
			expectedError: fmt.Errorf("expected json.Number value for NumberExact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.NumberExact(bigFloat),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as big.Float value for NumberExact check: number has no digits"),
		},
		"not-equal": {
			self:          knownvalue.NumberExact(bigFloat),
			other:         json.Number("1.797693134862315797693134862315797693134862314"),
			expectedError: fmt.Errorf("expected value 1.797693134862315797693134862315797693134862315 for NumberExact check, got: 1.797693134862315797693134862315797693134862314"),
		},
		"equal": {
			self:  knownvalue.NumberExact(bigFloat),
			other: json.Number("1.797693134862315797693134862315797693134862315"),
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

func TestNumberValue_String(t *testing.T) {
	t.Parallel()

	bigFloat, _, err := big.ParseFloat("1.797693134862315797693134862315797693134862315", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

	got := knownvalue.NumberExact(bigFloat).String()

	if diff := cmp.Diff(got, "1.797693134862315797693134862315797693134862315"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
