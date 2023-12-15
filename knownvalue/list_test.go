// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestListValue_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		other    any
		expected bool
	}{
		"nil": {},
		"wrong-type": {
			other: 1.23,
		},
		"empty": {
			other: []any{},
		},
		"different-len": {
			other: []any{
				int64(123),
				int64(456),
			},
		},
		"not-equal": {
			other: []any{
				int64(123),
				int64(456),
				int64(654),
			},
		},
		"wrong-order": {
			other: []any{
				int64(789),
				int64(456),
				int64(123),
			},
		},
		"equal": {
			other: []any{
				int64(123),
				int64(456),
				int64(789),
			},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.NewListValue([]knownvalue.KnownValue{
				knownvalue.NewInt64Value(123),
				knownvalue.NewInt64Value(456),
				knownvalue.NewInt64Value(789),
			}).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NewListValue([]knownvalue.KnownValue{
		knownvalue.NewInt64Value(123),
		knownvalue.NewInt64Value(456),
		knownvalue.NewInt64Value(789),
	}).String()

	if diff := cmp.Diff(got, "[123 456 789]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
