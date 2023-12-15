// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestFloat64Value_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		other    any
		expected bool
	}{
		"nil": {},
		"wrong-type": {
			other: 123,
		},
		"not-equal": {
			other: false,
		},
		"equal": {
			other:    1.23,
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.NewFloat64Value(1.23).Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestFloat64Value_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.NewFloat64Value(1.234567890123e+09).String()

	if diff := cmp.Diff(got, "1234567890.123"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
