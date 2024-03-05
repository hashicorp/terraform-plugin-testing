// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestNotNullValue_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.NotNull(),
			expectedError: fmt.Errorf("expected non-nil value for NotNull check, got: <nil>"),
		},
		"not-nil": {
			self:          knownvalue.NotNull(),
			other:         nil,
			expectedError: fmt.Errorf("expected non-nil value for NotNull check, got: <nil>"),
		},
		"equal": {
			self:  knownvalue.NotNull(),
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

func TestNotNullValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NotNull().String()

	if diff := cmp.Diff(got, "not-null"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
