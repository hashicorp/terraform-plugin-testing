// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestExpectKnownValue_CheckPlan_ResourceNotFound(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.two",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("test_resource.two - Resource not found in plan ResourceChanges"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_AttributeValueNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute value is null"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Bool(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(true),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Bool_KnownValueWrongType(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("bool_attribute"),
							knownvalue.Float64ValueExact(1.23),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: attribute value is bool, known value type is knownvalue.Float64Value"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Bool_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("bool_attribute"),
							knownvalue.BoolValueExact(false),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute value: true does not equal expected value: false"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Float64(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("float_attribute"),
							knownvalue.Float64ValueExact(1.23),
						),
					},
				},
			},
		},
	})
}

// We do not need equivalent tests for Int64 and Number as they all test the same logic.
func TestExpectKnownValue_CheckPlan_Float64_KnownValueWrongType(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("float_attribute"),
							knownvalue.StringValueExact("str"),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: attribute value is number, known value type is knownvalue.StringValue"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Float64_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("float_attribute"),
							knownvalue.Float64ValueExact(3.21),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute value: 1.23 does not equal expected value: 3.21"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Int64(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.Int64ValueExact(123),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Int64_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.Int64ValueExact(321),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute value: 123 does not equal expected value: 321"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_List(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.StringValueExact("value1"),
								knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_List_KnownValueWrongType(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.MapValueExact(map[string]knownvalue.KnownValue{}),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: attribute value is list, or set, known value type is knownvalue.MapValue"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_List_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.StringValueExact("value3"),
								knownvalue.StringValueExact("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`attribute value: \[value1 value2\] does not equal expected value: \[value3 value4\]`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListPartial(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								0: knownvalue.StringValueExact("value1"),
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
func TestExpectKnownValue_CheckPlan_ListPartial_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								0: knownvalue.StringValueExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`attribute value: \[0:value1 1:value2\] does not contain elements at the specified indices: \[0:value3\]`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNumElements(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NumElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNumElements_WrongNum(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NumElementsExact(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute contains 2 elements, expected 3"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlock(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.NewListValue([]knownvalue.KnownValue{
								knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.StringValueExact("str"),
								}),
								knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlockPartial(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.NewListValuePartial(map[int]knownvalue.KnownValue{
								1: knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
									"list_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlockNumElements(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.NumElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Map(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
								"key1": knownvalue.StringValueExact("value1"),
								"key2": knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Map_KnownValueWrongType(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NewListValue([]knownvalue.KnownValue{}),
						),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: attribute value is map, or object, known value type is knownvalue.ListValue"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Map_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
								"key3": knownvalue.StringValueExact("value3"),
								"key4": knownvalue.StringValueExact("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`attribute value: map\[key1:value1 key2:value2\] does not equal expected value: map\[key3:value3 key4:value4\]`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapPartial(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValuePartialMatch(map[string]knownvalue.KnownValue{
								"key1": knownvalue.StringValueExact("value1"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapPartial_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.MapValuePartialMatch(map[string]knownvalue.KnownValue{
								"key3": knownvalue.StringValueExact("value1"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`attribute value: map\[key1:value1 key2:value2\] does not contain: map\[key3:value1\]`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapNumElements(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NumElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapNumElements_WrongNum(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NumElementsExact(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute contains 2 elements, expected 3"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Number(t *testing.T) {
	t.Parallel()

	f, _, err := big.ParseFloat("123", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.NumberValueExact(f),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Number_KnownValueWrongValue(t *testing.T) {
	t.Parallel()

	f, _, err := big.ParseFloat("321", 10, 512, big.ToNearestEven)

	if err != nil {
		t.Errorf("%s", err)
	}

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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.NumberValueExact(f),
						),
					},
				},
				ExpectError: regexp.MustCompile("attribute value: 123 does not equal expected value: 321"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Set(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValueExact([]knownvalue.KnownValue{
								knownvalue.StringValueExact("value1"),
								knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_Set_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValueExact([]knownvalue.KnownValue{
								knownvalue.StringValueExact("value1"),
								knownvalue.StringValueExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`attribute value: \[value1 value2\] does not equal expected value: \[value1 value3\]`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetPartial(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValuePartialMatch([]knownvalue.KnownValue{
								knownvalue.StringValueExact("value2"),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetPartial_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.SetValuePartialMatch([]knownvalue.KnownValue{
								knownvalue.StringValueExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`attribute value: \[value1 value2\] does not contain: \[value3\]`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNumElements(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.NumElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlock(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetValueExact([]knownvalue.KnownValue{
								knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.StringValueExact("str"),
								}),
								knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlockPartial(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetValuePartialMatch([]knownvalue.KnownValue{
								knownvalue.MapValueExact(map[string]knownvalue.KnownValue{
									"set_nested_block_attribute": knownvalue.StringValueExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlockNumElements(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.NumElementsExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_String(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("string_attribute"),
							knownvalue.StringValueExact("str")),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_String_KnownValueWrongType(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("string_attribute"),
							knownvalue.BoolValueExact(true)),
					},
				},
				ExpectError: regexp.MustCompile("wrong type: attribute value is string, known value type is knownvalue.BoolValue"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_String_KnownValueWrongValue(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("string_attribute"),
							knownvalue.StringValueExact("rts")),
					},
				},
				ExpectError: regexp.MustCompile("attribute value: str does not equal expected value: rts"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_UnknownAttributeType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		knownValue  knownvalue.KnownValue
		req         plancheck.CheckPlanRequest
		expectedErr error
	}{
		"unrecognised-type": {
			knownValue: knownvalue.Int64ValueExact(123),
			req: plancheck.CheckPlanRequest{
				Plan: &tfjson.Plan{
					ResourceChanges: []*tfjson.ResourceChange{
						{
							Address: "example_resource.test",
							Change: &tfjson.Change{
								After: map[string]any{
									"attribute": float32(123),
								},
							},
						},
					},
				},
			},
			expectedErr: fmt.Errorf("unrecognised attribute type: float32, known value type is knownvalue.Int64Value\n\nThis is an error in plancheck.ExpectKnownValue.\nPlease report this to the maintainers."),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := plancheck.ExpectKnownValue("example_resource.test", tfjsonpath.New("attribute"), testCase.knownValue)

			resp := plancheck.CheckPlanResponse{}

			e.CheckPlan(context.Background(), testCase.req, &resp)

			if diff := cmp.Diff(resp.Error, testCase.expectedErr, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

var equateErrorMessage = cmp.Comparer(func(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}

	return x.Error() == y.Error()
})
