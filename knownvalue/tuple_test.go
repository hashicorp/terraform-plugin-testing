// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestTupleExact_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.TupleExact([]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected []any value for TupleExact check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.TupleExact([]knownvalue.Check{}),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			expectedError: fmt.Errorf("expected []any value for TupleExact check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for TupleExact check, got: float64"),
		},
		"empty": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("expected 3 elements for TupleExact check, got 0 elements"),
		},
		"wrong-length": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			other: []any{
				json.Number("123"),
				true,
			},
			expectedError: fmt.Errorf("expected 3 elements for TupleExact check, got 2 elements"),
		},
		"not-equal": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			other: []any{
				json.Number("123"),
				true,
				"goodbye",
			},
			expectedError: fmt.Errorf("tuple element index 2: expected value hello for StringExact check, got: goodbye"),
		},
		"wrong-order": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			other: []any{
				json.Number("123"),
				"hello",
				true,
			},
			expectedError: fmt.Errorf("tuple element index 1: expected bool value for Bool check, got: string"),
		},
		"equal": {
			self: knownvalue.TupleExact([]knownvalue.Check{
				knownvalue.Int64Exact(123),
				knownvalue.Bool(true),
				knownvalue.StringExact("hello"),
			}),
			other: []any{
				json.Number("123"),
				true,
				"hello",
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

func TestTupleExact_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.TupleExact([]knownvalue.Check{
		knownvalue.Int64Exact(123),
		knownvalue.Bool(true),
		knownvalue.StringExact("hello"),
	}).String()

	if diff := cmp.Diff(got, "[123 true hello]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
