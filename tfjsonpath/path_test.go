// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfjsonpath

import (
	"encoding/json"
	"strings"
	"testing"
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
				return strings.Contains(err.Error(), `path not found: specified key ObjectA not found in map`)
			},
		},
		{
			path: New("Object").AtMapKey("MapValueA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: specified key MapValueA not found in map`)
			},
		},

		// cannot convert object
		{
			path: New("StringValue").AtSliceIndex(0),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at SliceStep`)
			},
		},
		{
			path: New("StringValue").AtMapKey("MapKeyA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at MapStep`)
			},
		},
		{
			path: New("Array").AtSliceIndex(0).AtMapKey("MapValueA"),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: cannot convert object at MapStep`)
			},
		},

		// index out of bounds
		{
			path: New("Array").AtSliceIndex(10),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: SliceStep index 10 is out of range with slice length 7`)
			},
		},
		{
			path: New("Array").AtSliceIndex(7),
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), `path not found: SliceStep index 7 is out of range with slice length 7`)
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

func createTestObject() any {
	var jsonObject any
	jsonstring :=
		`{
		"StringValue": "example",
		"NumberValue": 0,
		"BooleanValue": false,
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
