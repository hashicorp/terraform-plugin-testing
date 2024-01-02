// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestListValuePartial_Equal(t *testing.T) {
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
		"wrong-order": {
			other: []any{1.23, 0.00, 7.89, 4.56},
		},
		"equal": {
			other:    []any{1.23, 0.00, 4.56, 7.89},
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
				0: knownvalue.Float64ValueExact(1.23),
				2: knownvalue.Float64ValueExact(4.56),
				3: knownvalue.Float64ValueExact(7.89),
			}).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListValuePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
		0: knownvalue.Float64ValueExact(1.23),
		2: knownvalue.Float64ValueExact(4.56),
		3: knownvalue.Float64ValueExact(7.89),
	}).String()

	if diff := cmp.Diff(got, "[0:1.23 2:4.56 3:7.89]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
