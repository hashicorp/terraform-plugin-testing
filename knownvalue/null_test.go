// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestNullValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self: knownvalue.Null(),
		},
		"not-nil": {
			self:          knownvalue.Null(),
			other:         false,
			expectedError: fmt.Errorf("expected nil value for Null check, got: bool"),
		},
		"equal": {
			self:  knownvalue.Null(),
			other: nil,
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

func TestNullValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.Null().String()

	if diff := cmp.Diff(got, "null"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
