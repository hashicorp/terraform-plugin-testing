// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		variable      config.Variable
		expected      []byte
		expectedError string
	}{
		"bool": {
			variable: config.BoolVariable(true),
			expected: []byte(`true`),
		},
		"float": {
			variable: config.FloatVariable(1.2),
			expected: []byte(`1.2`),
		},
		"integer": {
			variable: config.IntegerVariable(12),
			expected: []byte(`12`),
		},
		"list_bool": {
			variable: config.ListVariable(
				config.BoolVariable(false),
				config.BoolVariable(false),
				config.BoolVariable(true),
			),
			expected: []byte(`[false,false,true]`),
		},
		"list_list": {
			variable: config.ListVariable(
				config.ListVariable(
					config.BoolVariable(false),
					config.BoolVariable(false),
					config.BoolVariable(true),
				),
				config.ListVariable(
					config.BoolVariable(true),
					config.BoolVariable(true),
					config.BoolVariable(false),
				),
			),
			expected: []byte(`[[false,false,true],[true,true,false]]`),
		},
		"list_mixed_types": {
			variable: config.ListVariable(
				config.BoolVariable(false),
				config.StringVariable("str"),
			),
			expectedError: "lists must contain the same type",
		},
		"list_list_mixed_types": {
			variable: config.ListVariable(
				config.ListVariable(
					config.BoolVariable(false),
					config.StringVariable("str"),
				),
			),
			expectedError: "lists must contain the same type",
		},
		"list_list_mixed_types_multiple_lists": {
			variable: config.ListVariable(
				config.ListVariable(
					config.BoolVariable(false),
					config.BoolVariable(false),
				),
				config.ListVariable(
					config.StringVariable("str"),
					config.BoolVariable(false),
				),
			),
			expectedError: "lists must contain the same type",
		},
		"map_bool": {
			variable: config.MapVariable(
				map[string]config.Variable{
					"one":   config.BoolVariable(false),
					"two":   config.BoolVariable(false),
					"three": config.BoolVariable(true),
				},
			),
			expected: []byte(`{"one":false,"three":true,"two":false}`),
		},
		"map_map": {
			variable: config.ListVariable(
				config.MapVariable(
					map[string]config.Variable{
						"one":   config.BoolVariable(false),
						"two":   config.BoolVariable(false),
						"three": config.BoolVariable(true),
					},
				),
				config.MapVariable(
					map[string]config.Variable{
						"one":   config.BoolVariable(true),
						"two":   config.BoolVariable(true),
						"three": config.BoolVariable(false),
					},
				),
			),
			expected: []byte(`[{"one":false,"three":true,"two":false},{"one":true,"three":false,"two":true}]`),
		},
		"map_mixed_types": {
			variable: config.MapVariable(
				map[string]config.Variable{
					"one": config.BoolVariable(false),
					"two": config.StringVariable("str"),
				},
			),
			expectedError: "maps must contain the same type",
		},
		"map_map_mixed_types": {
			variable: config.MapVariable(
				map[string]config.Variable{
					"mapA": config.MapVariable(
						map[string]config.Variable{
							"one": config.BoolVariable(false),
							"two": config.StringVariable("str"),
						},
					),
				},
			),
			expectedError: "maps must contain the same type",
		},
		"map_map_mixed_types_multiple_maps": {
			variable: config.MapVariable(
				map[string]config.Variable{
					"mapA": config.MapVariable(
						map[string]config.Variable{
							"one": config.BoolVariable(false),
							"two": config.BoolVariable(true),
						},
					),
					"mapB": config.MapVariable(
						map[string]config.Variable{
							"one": config.BoolVariable(false),
							"two": config.StringVariable("str"),
						},
					),
				},
			),
			expectedError: "maps must contain the same type",
		},
		"object": {
			variable: config.ObjectVariable(
				map[string]config.Variable{
					"bool": config.BoolVariable(true),
					"list": config.ListVariable(
						config.BoolVariable(false),
						config.BoolVariable(true),
					),
					"map": config.MapVariable(
						map[string]config.Variable{
							"one": config.StringVariable("str_one"),
							"two": config.StringVariable("str_two"),
						},
					),
				},
			),
			expected: []byte(`{"bool":true,"list":[false,true],"map":{"one":"str_one","two":"str_two"}}`),
		},
		"object_map_mixed_types": {
			variable: config.ObjectVariable(
				map[string]config.Variable{
					"bool": config.BoolVariable(true),
					"list": config.ListVariable(
						config.BoolVariable(false),
						config.BoolVariable(true),
					),
					"map": config.MapVariable(
						map[string]config.Variable{
							"one": config.BoolVariable(false),
							"two": config.StringVariable("str_two"),
						},
					),
				},
			),
			expectedError: "maps must contain the same type",
		},
		"set_bool": {
			variable: config.SetVariable(
				config.BoolVariable(false),
				config.BoolVariable(false),
				config.BoolVariable(true),
			),
			expected: []byte(`[false,false,true]`),
		},
		"set_set": {
			variable: config.SetVariable(
				config.SetVariable(
					config.BoolVariable(false),
					config.BoolVariable(false),
					config.BoolVariable(true),
				),
				config.SetVariable(
					config.BoolVariable(true),
					config.BoolVariable(true),
					config.BoolVariable(false),
				),
			),
			expected: []byte(`[[false,false,true],[true,true,false]]`),
		},
		"set_mixed_types": {
			variable: config.SetVariable(
				config.BoolVariable(false),
				config.StringVariable("str"),
			),
			expectedError: "sets must contain the same type",
		},
		"set_set_mixed_types": {
			variable: config.SetVariable(
				config.SetVariable(
					config.BoolVariable(false),
					config.StringVariable("str"),
				),
			),
			expectedError: "sets must contain the same type",
		},
		"set_set_mixed_types_multiple_sets": {
			variable: config.SetVariable(
				config.SetVariable(
					config.BoolVariable(false),
					config.BoolVariable(false),
				),
				config.SetVariable(
					config.StringVariable("str"),
					config.BoolVariable(false),
				),
			),
			expectedError: "sets must contain the same type",
		},
		"set_non_unique": {
			variable: config.SetVariable(
				config.SetVariable(
					config.BoolVariable(false),
					config.BoolVariable(false),
				),
				config.SetVariable(
					config.BoolVariable(false),
					config.BoolVariable(false),
				),
			),
			expectedError: "sets must contain unique elements",
		},
		"string": {
			variable: config.StringVariable("str"),
			expected: []byte(`"str"`),
		},
		"tuple": {
			variable: config.TupleVariable(
				config.BoolVariable(true),
				config.FloatVariable(1.2),
				config.StringVariable("str"),
			),
			expected: []byte(`[true,1.2,"str"]`),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.variable.MarshalJSON()

			if testCase.expectedError == "" && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != "" && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != "" && err != nil {
				if diff := cmp.Diff(err.Error(), testCase.expectedError); diff != "" {
					t.Errorf("expected error %s, got error %s", testCase.expectedError, err)
				}
			}

			if !bytes.Equal(testCase.expected, got) {
				t.Errorf("expected %s, got %s", testCase.expected, got)
			}
		})
	}
}

func TestVariablesWrite(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	testCases := map[string]struct {
		variables     config.Variables
		expected      []byte
		expectedError string
	}{
		"write": {
			variables: map[string]config.Variable{
				"bool":   config.BoolVariable(true),
				"string": config.StringVariable("str"),
			},
			expected: []byte(`{"bool": true,"string": "str"}`),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.variables.Write(tempDir)

			if testCase.expectedError == "" && err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if testCase.expectedError != "" && err == nil {
				t.Errorf("expected error but got none")
			}

			if testCase.expectedError != "" && err != nil {
				if diff := cmp.Diff(err.Error(), testCase.expectedError); diff != "" {
					t.Errorf("expected error %s, got error %s", testCase.expectedError, err)
				}
			}

			b, err := os.ReadFile(filepath.Join(tempDir, "terraform-plugin-testing.auto.tfvars.json"))

			if err != nil {
				t.Errorf("error reading tfvars file: %s", err)
			}

			var expectedUnmarshalled map[string]any

			err = json.Unmarshal(testCase.expected, &expectedUnmarshalled)

			if err != nil {
				t.Errorf("error unmarshalling expected: %s", err)
			}

			var gotUnmarshalled map[string]any

			err = json.Unmarshal(b, &gotUnmarshalled)

			if err != nil {
				t.Errorf("error unmarshalling got: %s", err)
			}

			if diff := cmp.Diff(expectedUnmarshalled, gotUnmarshalled); diff != "" {
				t.Errorf("expected %s, got %s", expectedUnmarshalled, gotUnmarshalled)
			}
		})
	}
}
