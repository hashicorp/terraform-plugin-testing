// Copyright IBM Corp. 2014, 2025
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
							knownvalue.Bool(true),
						),
					},
				},
				ExpectError: regexp.MustCompile("test_resource.two - Resource not found in plan"),
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
							knownvalue.Null(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("float_attribute"),
							knownvalue.Null(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.Null(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.Null(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.ListExact([]knownvalue.Check{}),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.Null(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.Null(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetExact([]knownvalue.Check{}),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("string_attribute"),
							knownvalue.Null(),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_AttributeValueNotNull(t *testing.T) {
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
					float_attribute = 1.23
					int_attribute = 123
					list_attribute = ["value1", "value2"]
					list_nested_block {
						list_nested_block_attribute = "str"	
					}
					map_attribute = {
						key1 = "value1"		
					}	
					set_attribute = ["value1", "value2"]		
					set_nested_block {		
						set_nested_block_attribute = "str"	
					}
					string_attribute = "str"
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("bool_attribute"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("float_attribute"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("int_attribute"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_attribute"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("list_nested_block"),
							knownvalue.ListSizeExact(1),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("map_attribute"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_attribute"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("set_nested_block"),
							knownvalue.SetSizeExact(1),
						),
						plancheck.ExpectKnownValue(
							"test_resource.one",
							tfjsonpath.New("string_attribute"),
							knownvalue.NotNull(),
						),
					},
				},
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
							knownvalue.Bool(true),
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
							knownvalue.Float64Exact(1.23),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.bool_attribute, err: expected json\.Number value for Float64Exact check, got: bool`),
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
							knownvalue.Bool(false),
						),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.bool_attribute, err: expected value false for Bool check, got: true"),
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
							knownvalue.Float64Exact(1.23),
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
							knownvalue.StringExact("str"),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.float_attribute, err: expected string value for StringExact check, got: json\.Number`),
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
							knownvalue.Float64Exact(3.21),
						),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.float_attribute, err: expected value 3.21 for Float64Exact check, got: 1.23"),
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
							knownvalue.Int64Exact(123),
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
							knownvalue.Int64Exact(321),
						),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.int_attribute, err: expected value 321 for Int64Exact check, got: 123"),
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
							knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("value1"),
								knownvalue.StringExact("value2"),
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
							knownvalue.MapExact(map[string]knownvalue.Check{}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.list_attribute, err: expected map\[string\]any value for MapExact check, got: \[\]interface {}`),
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
							knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("value3"),
								knownvalue.StringExact("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.list_attribute, err: list element index 0: expected value value3 for StringExact check, got: value1`),
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
							knownvalue.ListPartial(map[int]knownvalue.Check{
								0: knownvalue.StringExact("value1"),
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
							knownvalue.ListPartial(map[int]knownvalue.Check{
								0: knownvalue.StringExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.list_attribute, err: list element 0: expected value value3 for StringExact check, got: value1`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListElements(t *testing.T) {
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
							knownvalue.ListSizeExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListElements_WrongNum(t *testing.T) {
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
							knownvalue.ListSizeExact(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.list_attribute, err: expected 3 elements for ListSizeExact check, got 2 elements"),
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
							knownvalue.ListExact([]knownvalue.Check{
								knownvalue.MapExact(map[string]knownvalue.Check{
									"list_nested_block_attribute": knownvalue.StringExact("str"),
								}),
								knownvalue.MapExact(map[string]knownvalue.Check{
									"list_nested_block_attribute": knownvalue.StringExact("rts"),
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
							knownvalue.ListPartial(map[int]knownvalue.Check{
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"list_nested_block_attribute": knownvalue.StringExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_ListNestedBlockElements(t *testing.T) {
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
							knownvalue.ListSizeExact(2),
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
							knownvalue.MapExact(map[string]knownvalue.Check{
								"key1": knownvalue.StringExact("value1"),
								"key2": knownvalue.StringExact("value2"),
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
							knownvalue.ListExact([]knownvalue.Check{}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.map_attribute, err: expected \[\]any value for ListExact check, got: map\[string\]interface {}`),
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
							knownvalue.MapExact(map[string]knownvalue.Check{
								"key3": knownvalue.StringExact("value3"),
								"key4": knownvalue.StringExact("value4"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.map_attribute, err: missing element key3 for MapExact check`),
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
							knownvalue.MapPartial(map[string]knownvalue.Check{
								"key1": knownvalue.StringExact("value1"),
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
							knownvalue.MapPartial(map[string]knownvalue.Check{
								"key3": knownvalue.StringExact("value1"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.map_attribute, err: missing element key3 for MapPartial check`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapElements(t *testing.T) {
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
							knownvalue.MapSizeExact(2),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_MapElements_WrongNum(t *testing.T) {
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
							knownvalue.MapSizeExact(3),
						),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.map_attribute, err: expected 3 elements for MapSizeExact check, got 2 elements"),
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
							knownvalue.NumberExact(f),
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
							knownvalue.NumberExact(f),
						),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.int_attribute, err: expected value 321 for NumberExact check, got: 123"),
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
							knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("value2"),
								knownvalue.StringExact("value1"),
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
							knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("value1"),
								knownvalue.StringExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.set_attribute, err: missing value value3 for SetExact check`),
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
							knownvalue.SetPartial([]knownvalue.Check{
								knownvalue.StringExact("value2"),
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
							knownvalue.SetPartial([]knownvalue.Check{
								knownvalue.StringExact("value3"),
							}),
						),
					},
				},
				ExpectError: regexp.MustCompile(`error checking value for attribute at path: test_resource.one.set_attribute, err: missing value value3 for SetPartial check`),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetElements(t *testing.T) {
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
							knownvalue.SetSizeExact(2),
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
							knownvalue.SetExact([]knownvalue.Check{
								knownvalue.MapExact(map[string]knownvalue.Check{
									"set_nested_block_attribute": knownvalue.StringExact("str"),
								}),
								knownvalue.MapExact(map[string]knownvalue.Check{
									"set_nested_block_attribute": knownvalue.StringExact("rts"),
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
							knownvalue.SetPartial([]knownvalue.Check{
								knownvalue.MapExact(map[string]knownvalue.Check{
									"set_nested_block_attribute": knownvalue.StringExact("rts"),
								}),
							}),
						),
					},
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_SetNestedBlockElements(t *testing.T) {
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
							knownvalue.SetSizeExact(2),
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
							knownvalue.StringExact("str")),
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
							knownvalue.Bool(true)),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.string_attribute, err: expected bool value for Bool check, got: string"),
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
							knownvalue.StringExact("rts")),
					},
				},
				ExpectError: regexp.MustCompile("error checking value for attribute at path: test_resource.one.string_attribute, err: expected value rts for StringExact check, got: str"),
			},
		},
	})
}

func TestExpectKnownValue_CheckPlan_UnknownAttributeType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		knownValue  knownvalue.Check
		req         plancheck.CheckPlanRequest
		expectedErr error
	}{
		"unrecognised-type": {
			knownValue: knownvalue.Int64Exact(123),
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
			expectedErr: fmt.Errorf("error checking value for attribute at path: example_resource.test.attribute, err: expected json.Number value for Int64Exact check, got: float32"),
		},
	}

	for name, testCase := range testCases {
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
