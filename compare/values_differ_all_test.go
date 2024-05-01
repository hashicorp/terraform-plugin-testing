// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

func TestValuesDifferAll_CompareValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            []any
		expectedError error
	}{
		"nil": {},
		"single-value": {
			in: []any{"str"},
		},
		"non-matching-values": {
			in: []any{"str", "other_str", "another_str"},
		},
		"matching-values": {
			in:            []any{"str", "str", "str"},
			expectedError: fmt.Errorf("expected values to differ, but value is duplicated: str"),
		},
		"non-sequential-matching-values": {
			in:            []any{"str", "other_str", "str"},
			expectedError: fmt.Errorf("expected values to differ, but value is duplicated: str"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := compare.ValuesDifferAll().CompareValues(testCase.in...)

			if diff := cmp.Diff(err, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
