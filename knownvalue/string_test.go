// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestStringValue_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		other    any
		expected bool
	}{
		"nil": {},
		"wrong-type": {
			other: 1.23,
		},
		"not-equal": {
			other: "other",
		},
		"equal": {
			other:    "str",
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := knownvalue.StringValueExact("str").Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStringValue_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.StringValueExact("str").String()

	if diff := cmp.Diff(got, "str"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
