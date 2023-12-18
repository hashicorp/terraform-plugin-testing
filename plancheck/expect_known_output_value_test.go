// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestExpectKnownOutputValue_CheckPlan_OutputNotFound(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output bool_output {
					value = test_resource.one.bool_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"bool_not_found",
							knownvalue.NewBoolValue(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("bool_not_found - Output not found in plan OutputChanges"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_AttributeValueNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}
				output bool_output {
					value = test_resource.one.bool_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"bool_output",
							knownvalue.NewBoolValue(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("output value is null"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Bool(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}
		
				output bool_output {
					value = test_resource.one.bool_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"bool_output",
							knownvalue.NewBoolValue(true),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Bool_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output bool_output {
					value = test_resource.one.bool_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"bool_output",
							knownvalue.NewFloat64Value(1.23),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: output value is bool, known value type is knownvalue.Float64Value"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Bool_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					bool_attribute = true
				}

				output bool_output {
					value = test_resource.one.bool_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"bool_output",
							knownvalue.NewBoolValue(false),
						),
					},
				},
				ExpectError: regexp.MustCompile("value: true does not equal expected value: false"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Float64(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}

				output float64_output {
					value = test_resource.one.float_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"float64_output",
							knownvalue.NewFloat64Value(1.23),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Float64_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}

				output float64_output {
					value = test_resource.one.float_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"float64_output",
							knownvalue.NewStringValue("str"),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: output value is float64 or int64, known value type is knownvalue.StringValue"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Float64_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					float_attribute = 1.23
				}

				output float64_output {
					value = test_resource.one.float_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"float64_output",
							knownvalue.NewFloat64Value(3.21),
						),
					},
				},
				ExpectError: regexp.MustCompile("output value: 1.23 does not equal expected value: 3.21"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Int64(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output int64_output {
					value = test_resource.one.int_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"int64_output",
							knownvalue.NewInt64Value(123),
						),
					},
				},
			},
		},
	})
}

// TestExpectKnownOutputValue_CheckPlan_Int64_KnownValueWrongType highlights a limitation of tfjson.Plan in that all numerical
// values are represented as float64.
func TestExpectKnownOutputValue_CheckPlan_Int64_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output int64_output {
					value = test_resource.one.int_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"int64_output",
							knownvalue.NewStringValue("str"),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: output value is float64 or int64, known value type is knownvalue.StringValue"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Int64_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					int_attribute = 123
				}

				output int64_output {
					value = test_resource.one.int_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"int64_output",
							knownvalue.NewInt64Value(321),
						),
					},
				},
				ExpectError: regexp.MustCompile("output value: 123 does not equal expected value: 321"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_List(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value1"),
								knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_List_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}
				
				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewMapValue(map[string]knownvalue.KnownValue{}),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: output type is list, or set, known value type is knownvalue.MapValue"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_List_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value3"),
								knownvalue.NewStringValue("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`output value: \[value1 value2\] does not equal expected value: \[value3 value4\]`),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_ListPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								0: knownvalue.NewStringValue("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

// No need to check KnownValueWrongType for ListPartial as all lists, and sets are []any in
// tfjson.Plan.
func TestExpectKnownOutputValue_CheckPlan_ListPartial_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								0: knownvalue.NewStringValue("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`output value: \[0:value1 1:value2\] does not contain elements at the specified indices: \[0:value3\]`),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_ListNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_ListNumElements_WrongNum(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_attribute = [
						"value1",
						"value2"
					]
				}

				output list_output {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_output",
							knownvalue.NewNumElements(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("output contains 2 elements, expected 3"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_ListNestedBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
					list_nested_block {
						list_nested_block_attribute = "rts"
					}
				}

				output list_nested_block_output {
					value = test_resource.one.list_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_nested_block_output",
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.NewStringValue("str"),
								}),
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_ListNestedBlockPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
					list_nested_block {
						list_nested_block_attribute = "rts"
					}
				}

				output list_nested_block_output {
					value = test_resource.one.list_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_nested_block_output",
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								1: knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_ListNestedBlockNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
					list_nested_block {
						list_nested_block_attribute = "rts"
					}
				}

				output list_nested_block_output {
					value = test_resource.one.list_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"list_nested_block_output",
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Map(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
								"key1": knownvalue.NewStringValue("value1"),
								"key2": knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Map_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewListValue([]knownvalue.KnownValue{}),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: output type is map, or object, known value type is knownvalue.ListValue"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Map_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
								"key3": knownvalue.NewStringValue("value3"),
								"key4": knownvalue.NewStringValue("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`output value: map\[key1:value1 key2:value2\] does not equal expected value: map\[key3:value3 key4:value4\]`),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_MapPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewMapValuePartial(map[string]knownvalue.KnownValue{
								"key1": knownvalue.NewStringValue("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_MapPartial_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewMapValuePartial(map[string]knownvalue.KnownValue{
								"key3": knownvalue.NewStringValue("value1"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`output value: map\[key1:value1 key2:value2\] does not contain: map\[key3:value1\]`),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_MapNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_MapNumElements_WrongNum(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					map_attribute = {
						key1 = "value1"
						key2 = "value2"
					}
				}

				output map_output {
					value = test_resource.one.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"map_output",
							knownvalue.NewNumElements(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("output contains 2 elements, expected 3"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Set(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output set_output {
					value = test_resource.one.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_output",
							knownvalue.NewSetValue([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value1"),
								knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_Set_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output set_output {
					value = test_resource.one.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_output",
							knownvalue.NewSetValue([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value1"),
								knownvalue.NewStringValue("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`output value: \[value1 value2\] does not equal expected value: \[value1 value3\]`),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_SetPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output set_output {
					value = test_resource.one.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_output",
							knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_SetPartial_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output set_output {
					value = test_resource.one.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_output",
							knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
								knownvalue.NewStringValue("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`output value: \[value1 value2\] does not contain: \[value3\]`),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_SetNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_attribute = [
						"value1",
						"value2"
					]
				}

				output set_output {
					value = test_resource.one.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_output",
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_SetNestedBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_nested_block {
						set_nested_block_attribute = "str"
					}
					set_nested_block {
						set_nested_block_attribute = "rts"
					}
				}

				output set_nested_block_output {
					value = test_resource.one.set_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_nested_block_output",
							knownvalue.NewSetValue([]knownvalue.KnownValue{
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.NewStringValue("str"),
								}),
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_SetNestedBlockPartial(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_nested_block {
						set_nested_block_attribute = "str"
					}
					set_nested_block {
						set_nested_block_attribute = "rts"
					}
				}

				output set_nested_block_output {
					value = test_resource.one.set_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_nested_block_output",
							knownvalue.NewSetValuePartial([]knownvalue.KnownValue{
								knownvalue.NewMapValue(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.NewStringValue("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_SetNestedBlockNumElements(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					set_nested_block {
						set_nested_block_attribute = "str"
					}
					set_nested_block {
						set_nested_block_attribute = "rts"
					}
				}

				output set_nested_block_output {
					value = test_resource.one.set_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"set_nested_block_output",
							knownvalue.NewNumElements(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_String(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}

				output string_output {
					value = test_resource.one.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"string_output",
							knownvalue.NewStringValue("str")),
					},
				},
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_String_KnownValueWrongType(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}

				output string_output {
					value = test_resource.one.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"string_output",
							knownvalue.NewBoolValue(true)),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: output value is string, known value type is knownvalue.BoolValue"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_String_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					string_attribute = "str"
				}

				output string_output {
					value = test_resource.one.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"string_output",
							knownvalue.NewStringValue("rts")),
					},
				},
				ExpectError: regexp.MustCompile("output value: str does not equal expected value: rts"),
			},
		},
	})
}

func TestExpectKnownOutputValue_CheckPlan_UnknownAttributeType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		knownValue  knownvalue.KnownValue
		req         plancheck.CheckPlanRequest
		expectedErr error
	}{
		"unrecognised-type": {
			knownValue: knownvalue.NewInt64Value(123),
			req: plancheck.CheckPlanRequest{
				Plan: &tfjson.Plan{
					OutputChanges: map[string]*tfjson.Change{
						"float32_output": {
							After: float32(123),
						},
					},
				},
			},
			expectedErr: fmt.Errorf("unrecognised output type: float32, known value type is knownvalue.Int64Value\n\nThis is an error in plancheck.ExpectKnownOutputValue.\nPlease report this to the maintainers."),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := plancheck.ExpectKnownOutputValue("float32_output", testCase.knownValue)

			resp := plancheck.CheckPlanResponse{}

			e.CheckPlan(context.Background(), testCase.req, &resp)

			if diff := cmp.Diff(resp.Error, testCase.expectedErr, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
