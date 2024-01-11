// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestObjectAttributes_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.ObjectAttributesExact(0),
			expectedError: fmt.Errorf("expected map[string]any value for ObjectAttributesExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.ObjectAttributesExact(0),
			other: map[string]any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self:          knownvalue.ObjectAttributesExact(3),
			expectedError: fmt.Errorf("expected map[string]any value for ObjectAttributesExact check, got: <nil>"),
		},
		"wrong-type": {
			self:          knownvalue.ObjectAttributesExact(3),
			other:         1.234,
			expectedError: fmt.Errorf("expected map[string]any value for ObjectAttributesExact check, got: float64"),
		},
		"empty": {
			self:          knownvalue.ObjectAttributesExact(3),
			other:         map[string]any{},
			expectedError: fmt.Errorf("expected 3 attributes for ObjectAttributesExact check, got 0 attributes"),
		},
		"wrong-length": {
			self: knownvalue.ObjectAttributesExact(3),
			other: map[string]any{
				"one": int64(123),
				"two": int64(456),
			},
			expectedError: fmt.Errorf("expected 3 attributes for ObjectAttributesExact check, got 2 attributes"),
		},
		"equal": {
			self: knownvalue.ObjectAttributesExact(3),
			other: map[string]any{
				"one":   int64(123),
				"two":   int64(456),
				"three": int64(789),
			},
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

func TestObjectAttributes_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.ObjectAttributesExact(2).String()

	if diff := cmp.Diff(got, "2"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
