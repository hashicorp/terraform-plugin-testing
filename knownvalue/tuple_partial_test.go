// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestTuplePartial_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          knownvalue.TuplePartial(map[int]knownvalue.Check{}),
			expectedError: fmt.Errorf("expected []any value for TuplePartial check, got: <nil>"),
		},
		"zero-other": {
			self:  knownvalue.TuplePartial(map[int]knownvalue.Check{}),
			other: []any{}, // checking against the underlying value field zero-value
		},
		"nil": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			expectedError: fmt.Errorf("expected []any value for TuplePartial check, got: <nil>"),
		},
		"wrong-type": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			other:         1.234,
			expectedError: fmt.Errorf("expected []any value for TuplePartial check, got: float64"),
		},
		"empty": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			other:         []any{},
			expectedError: fmt.Errorf("missing element index 0 for TuplePartial check"),
		},
		"wrong-length": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			other: []any{
				json.Number("1.23"),
				"hello",
			},
			expectedError: fmt.Errorf("missing element index 2 for TuplePartial check"),
		},
		"not-equal": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			other: []any{
				json.Number("1.23"),
				"world",
				"hello",
			},
			expectedError: fmt.Errorf("tuple element 2: expected value world for StringExact check, got: hello"),
		},
		"wrong-order": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			other: []any{
				json.Number("1.23"),
				"world",
				true,
			},
			expectedError: fmt.Errorf("tuple element 2: expected string value for StringExact check, got: bool"),
		},
		"equal": {
			self: knownvalue.TuplePartial(map[int]knownvalue.Check{
				0: knownvalue.Float64Exact(1.23),
				2: knownvalue.StringExact("world"),
				3: knownvalue.Bool(true),
			}),
			other: []any{
				json.Number("1.23"),
				"hello",
				"world",
				true,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.self.CheckValue(testCase.other)

			if diff := cmp.Diff(got, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTuplePartial_String(t *testing.T) {
	t.Parallel()

	got := knownvalue.TuplePartial(map[int]knownvalue.Check{
		0: knownvalue.Float64Exact(1.23),
		2: knownvalue.StringExact("world"),
		3: knownvalue.Bool(true),
	}).String()

	if diff := cmp.Diff(got, "[0:1.23 2:world 3:true]"); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}
