// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package tfjsonpath

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Traverse_StringValue(t *testing.T) {
	t.Parallel()

	path := New("StringValue")

	actual, err := Traverse(createTestObject(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}
	expected := "example"

	if expected != actual {
		t.Errorf("Output %v not equal to expected %v", actual, expected)
	}
}

func Test_Traverse_Array_StringValue(t *testing.T) {
	t.Parallel()

	path := New(0).AtMapKey("StringValue")

	actual, err := Traverse(createTestArray(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}
	expected := "example"

	if expected != actual {
		t.Errorf("Output %v not equal to expected %v", actual, expected)
	}
}

func Test_Traverse_NumberValue(t *testing.T) {
	t.Parallel()

	path := New("NumberValue")

	actual, err := Traverse(createTestObject(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}
	expected := 0.0

	if expected != actual {
		t.Errorf("Output %v not equal to expected %v", actual, expected)
	}
}

func Test_Traverse_Array_NumberValue(t *testing.T) {
	t.Parallel()

	path := New(0).AtMapKey("NumberValue")

	actual, err := Traverse(createTestArray(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}
	expected := 0.0

	if expected != actual {
		t.Errorf("Output %v not equal to expected %v", actual, expected)
	}
}

func Test_Traverse_BooleanValue(t *testing.T) {
	t.Parallel()

	path := New("BooleanValue")

	actual, err := Traverse(createTestObject(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}
	expected := false

	if expected != actual {
		t.Errorf("Output %v not equal to expected %v", actual, expected)
	}
}

func Test_Traverse_Array_BooleanValue(t *testing.T) {
	t.Parallel()

	path := New(0).AtMapKey("BooleanValue")

	actual, err := Traverse(createTestArray(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}
	expected := false

	if expected != actual {
		t.Errorf("Output %v not equal to expected %v", actual, expected)
	}
}

func Test_Traverse_NullValue(t *testing.T) {
	t.Parallel()

	path := New("NullValue")

	actual, err := Traverse(createTestObject(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}

	if actual != nil {
		t.Errorf("Output %v not equal to expected %v", actual, nil)
	}
}

func Test_Traverse_Array_NullValue(t *testing.T) {
	t.Parallel()

	path := New(0).AtMapKey("NullValue")

	actual, err := Traverse(createTestArray(), path)
	if err != nil {
		t.Errorf("Error traversing JSON object %s", err)
	}

	if actual != nil {
		t.Errorf("Output %v not equal to expected %v", actual, nil)
	}
}

func Test_Traverse_Array(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path     Path
		expected any
	}{
		{
			path:     New("Array").AtSliceIndex(0),
			expected: 10.0,
		},
		{
			path:     New("Array").AtSliceIndex(1),
			expected: 15.2,
		},
		{
			path:     New("Array").AtSliceIndex(2),
			expected: "example2",
		},
		{
			path:     New("Array").AtSliceIndex(3),
			expected: nil,
		},
		{
			path:     New("Array").AtSliceIndex(4),
			expected: true,
		},
		{
			path:     New("Array").AtSliceIndex(5).AtMapKey("NestedStringValue"),
			expected: "example3",
		},
		{
			path:     New("Array").AtSliceIndex(6).AtSliceIndex(0),
			expected: true,
		},
	}

	for _, tc := range testCases {
		actual, err := Traverse(createTestObject(), tc.path)
		if err != nil {
			t.Errorf("Error traversing JSON object %s", err)
		}
		expected := tc.expected

		if expected != actual {
			t.Errorf("Output %v not equal to expected %v", actual, expected)
		}
	}
}

func Test_Traverse_Array_Array(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path     Path
		expected any
	}{
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(0),
			expected: 10.0,
		},
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(1),
			expected: 15.2,
		},
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(2),
			expected: "example2",
		},
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(3),
			expected: nil,
		},
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(4),
			expected: true,
		},
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(5).AtMapKey("NestedStringValue"),
			expected: "example3",
		},
		{
			path:     New(0).AtMapKey("Array").AtSliceIndex(6).AtSliceIndex(0),
			expected: true,
		},
	}

	for _, tc := range testCases {
		actual, err := Traverse(createTestArray(), tc.path)
		if err != nil {
			t.Errorf("Error traversing JSON object %s", err)
		}
		expected := tc.expected

		if expected != actual {
			t.Errorf("Output %v not equal to expected %v", actual, expected)
		}
	}
}

func Test_Traverse_Object(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path     Path
		expected any
	}{
		{
			path:     New("Object").AtMapKey("StringValue"),
			expected: "example",
		},
		{
			path:     New("Object").AtMapKey("NumberValue"),
			expected: 0.0,
		},
		{
			path:     New("Object").AtMapKey("BooleanValue"),
			expected: false,
		},
		{
			path:     New("Object").AtMapKey("ArrayValue").AtSliceIndex(0),
			expected: 10.0,
		},
		{
			path:     New("Object").AtMapKey("ArrayValue").AtSliceIndex(1),
			expected: 15.2,
		},
		{
			path:     New("Object").AtMapKey("ArrayValue").AtSliceIndex(2),
			expected: "example2",
		},
		{
			path:     New("Object").AtMapKey("ArrayValue").AtSliceIndex(3),
			expected: nil,
		},
		{
			path:     New("Object").AtMapKey("ArrayValue").AtSliceIndex(4),
			expected: true,
		},
		{
			path:     New("Object").AtMapKey("ObjectValue").AtMapKey("NestedStringValue"),
			expected: "example3",
		},
	}

	for _, tc := range testCases {
		actual, err := Traverse(createTestObject(), tc.path)
		if err != nil {
			t.Errorf("Error traversing JSON object %s", err)
		}
		expected := tc.expected

		if expected != actual {
			t.Errorf("Output %v not equal to expected %v", actual, expected)
		}
	}
}

func Test_Traverse_Array_Object(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path     Path
		expected any
	}{
		{
			path:     New(0).AtMapKey("Object").AtMapKey("StringValue"),
			expected: "example",
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("NumberValue"),
			expected: 0.0,
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("BooleanValue"),
			expected: false,
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("ArrayValue").AtSliceIndex(0),
			expected: 10.0,
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("ArrayValue").AtSliceIndex(1),
			expected: 15.2,
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("ArrayValue").AtSliceIndex(2),
			expected: "example2",
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("ArrayValue").AtSliceIndex(3),
			expected: nil,
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("ArrayValue").AtSliceIndex(4),
			expected: true,
		},
		{
			path:     New(0).AtMapKey("Object").AtMapKey("ObjectValue").AtMapKey("NestedStringValue"),
			expected: "example3",
		},
	}

	for _, tc := range testCases {
		actual, err := Traverse(createTestArray(), tc.path)
		if err != nil {
			t.Errorf("Error traversing JSON object %s", err)
		}
		expected := tc.expected

		if expected != actual {
			t.Errorf("Output %v not equal to expected %v", actual, expected)
		}
	}
}

func Test_Traverse_ExpectError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path          Path
		expectedError func(err error) bool
	}{
		// specified key not found
		{
			path: New("ObjectA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: specified key ObjectA not found in map at ObjectA`)
			},
		},
		{
			path: New("Object").AtMapKey("MapValueA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: specified key MapValueA not found in map at Object.MapValueA`)
			},
		},

		// cannot convert object
		{
			path: New("StringValue").AtSliceIndex(0),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at SliceStep StringValue.0 to []any`)
			},
		},
		{
			path: New("StringValue").AtMapKey("MapKeyA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at MapStep StringValue.MapKeyA to map[string]any`)
			},
		},
		{
			path: New("Array").AtSliceIndex(0).AtMapKey("MapValueA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at MapStep Array.0.MapValueA to map[string]any`)
			},
		},

		// index out of bounds
		{
			path: New("Array").AtSliceIndex(10),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: SliceStep index Array.10 is out of range with slice length 7`)
			},
		},
		{
			path: New("Array").AtSliceIndex(7),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: SliceStep index Array.7 is out of range with slice length 7`)
			},
		},
	}

	for _, tc := range testCases {
		_, err := Traverse(createTestObject(), tc.path)
		if err == nil {
			t.Fatalf("Expected error but got none")
		}

		if !tc.expectedError(err) {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func Test_Traverse_Array_ExpectError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		path          Path
		expectedError func(err error) bool
	}{
		// specified index not found
		"unknown_index": {
			path: New(1),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: SliceStep index 1 is out of range with slice length 1`)
			},
		},
		"unknown_nested_index": {
			path: New(0).AtSliceIndex(0),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at SliceStep 0.0 to []any`)
			},
		},

		// cannot convert object
		"unknown_map_index": {
			path: New(0).AtMapKey("StringValue").AtSliceIndex(0),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at SliceStep 0.StringValue.0 to []any`)
			},
		},
		"unknown_map_key": {
			path: New(0).AtMapKey("StringValue").AtMapKey("MapKeyA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at MapStep 0.StringValue.MapKeyA to map[string]any`)
			},
		},
		"unknown_slice_map_key": {
			path: New(0).AtMapKey("Array").AtSliceIndex(0).AtMapKey("MapValueA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at MapStep 0.Array.0.MapValueA to map[string]any`)
			},
		},

		// index out of bounds
		"out_of_bounds": {
			path: New(0).AtMapKey("Array").AtSliceIndex(10),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: SliceStep index 0.Array.10 is out of range with slice length 7`)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := Traverse(createTestArray(), tc.path)

			if err == nil {
				t.Fatalf("Expected error but got none")
			}

			if !tc.expectedError(err) {
				t.Errorf("Unexpected error: %s", err)
			}
		})
	}
}

func TestPath_String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		path     Path
		expected string
	}{
		"slice_step": {
			path:     New(1),
			expected: "1",
		},
		"map_step": {
			path:     New("attr"),
			expected: "attr",
		},
		"slice_step_map_step": {
			path:     New(0).AtMapKey("attr"),
			expected: "0.attr",
		},
		"map_step_slice_step": {
			path:     New("attr").AtSliceIndex(0),
			expected: "attr.0",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tc.path.String()

			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func createTestObject() any {
	var jsonObject any
	jsonstring :=
		`{
		"StringValue": "example",
		"NumberValue": 0,
		"BooleanValue": false,
		"NullValue": null,
		"Array": [10, 15.2, "example2", null, true, {"NestedStringValue": "example3"}, [true]],
		"Object":{
			"StringValue": "example",
			"NumberValue": 0,
			"BooleanValue": false,
			"ArrayValue": [10, 15.2, "example2", null, true],
			"ObjectValue": {
				"NestedStringValue": "example3"
			}
		}
	}`
	err := json.Unmarshal([]byte(jsonstring), &jsonObject)
	if err != nil {
		return nil
	}

	return jsonObject
}

func createTestArray() any {
	var jsonObject any
	jsonstring :=
		`[{
		"StringValue": "example",
		"NumberValue": 0,
		"BooleanValue": false,
		"NullValue": null,
		"Array": [10, 15.2, "example2", null, true, {"NestedStringValue": "example3"}, [true]],
		"Object":{
			"StringValue": "example",
			"NumberValue": 0,
			"BooleanValue": false,
			"ArrayValue": [10, 15.2, "example2", null, true],
			"ObjectValue": {
				"NestedStringValue": "example3"
			}
		}
	}]`
	err := json.Unmarshal([]byte(jsonstring), &jsonObject)
	if err != nil {
		return nil
	}

	return jsonObject
}
