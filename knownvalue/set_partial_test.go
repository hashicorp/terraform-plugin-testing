// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestSetValuePartial_Equal(t *testing.T) {
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
			other: []any{1.23, 4.56},
		},
		"not-equal": {
			other: []any{1.23, 4.56, 6.54, 5.46},
		},
		"equal-different-order": {
			other:    []any{1.23, 0.00, 7.89, 4.56},
			expected: true,
		},
		"equal-same-order": {
			other:    []any{1.23, 0.00, 4.56, 7.89},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
				knownvalue.NewFloat64Value(1.23),
				knownvalue.NewFloat64Value(4.56),
				knownvalue.NewFloat64Value(7.89),
			}).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestSetValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
		knownvalue.NewFloat64Value(1.23),
		knownvalue.NewFloat64Value(4.56),
		knownvalue.NewFloat64Value(7.89),
	}).String()

	if diff := cmp.Diff(got, "[1.23 4.56 7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
