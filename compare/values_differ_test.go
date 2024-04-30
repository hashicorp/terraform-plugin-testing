// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

func TestValuesDiffer_CompareValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in            []any
		expectedError error
	}{
		"nil": {},
		"single-value": {
			in: []any{"str"},
		},
		"non-matching-sequential-values": {
			in: []any{"str", "other_str", "str"},
		},
		"matching-values": {
			in:            []any{"str", "other_str", "other_str"},
			expectedError: fmt.Errorf("expected values to differ, but they are the same: other_str == other_str"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := compare.ValuesDiffer{}.CompareValues(testCase.in...)

			if diff := cmp.Diff(err, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
