// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestExpectKnownValue_CheckState_ResourceNotFound(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.two",
						tfjsonpath.New("bool_attribute"),
						knownvalue.Bool(true),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two - Resource not found in state"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_AttributeValueNull(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_nested_block"),
						knownvalue.ListExact([]knownvalue.Check{}),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_attribute"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_nested_block"),
						knownvalue.SetExact([]knownvalue.Check{}),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_AttributeValueNotNull(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_nested_block"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_attribute"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_nested_block"),
						knownvalue.SetSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Bool(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Bool_KnownValueWrongType(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						knownvalue.Float64Exact(1.23),
					),
				},
				ExpectError: regexp.MustCompile(`expected json\.Number value for Float64Exact check, got: bool`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Bool_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						knownvalue.Bool(false),
					),
				},
				ExpectError: regexp.MustCompile("expected value false for Bool check, got: true"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Float64(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						knownvalue.Float64Exact(1.23),
					),
				},
			},
		},
	})
}

// We do not need equivalent tests for Int64 and Number as they all test the same logic.
func TestExpectKnownValue_CheckState_Float64_KnownValueWrongType(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						knownvalue.StringExact("str"),
					),
				},
				ExpectError: regexp.MustCompile(`expected string value for StringExact check, got: json\.Number`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Float64_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						knownvalue.Float64Exact(3.21),
					),
				},
				ExpectError: regexp.MustCompile("expected value 3.21 for Float64Exact check, got: 1.23"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Int64(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						knownvalue.Int64Exact(123),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Int64_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						knownvalue.Int64Exact(321),
					),
				},
				ExpectError: regexp.MustCompile("expected value 321 for Int64Exact check, got: 123"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_List(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_List_KnownValueWrongType(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.MapExact(map[string]knownvalue.Check{}),
					),
				},
				ExpectError: regexp.MustCompile(`expected map\[string\]any value for MapExact check, got: \[\]interface {}`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_List_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("value3"),
							knownvalue.StringExact("value4"),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`list element index 0: expected value value3 for StringExact check, got: value1`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_ListPartial(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.StringExact("value1"),
						}),
					),
				},
			},
		},
	})
}

// No need to check KnownValueWrongType for ListPartial as all lists, and sets are []any in
// tfjson.State.
func TestExpectKnownValue_CheckState_ListPartial_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.StringExact("value3"),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`list element 0: expected value value3 for StringExact check, got: value1`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_ListElements(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.ListSizeExact(2),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_ListElements_WrongNum(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_attribute"),
						knownvalue.ListSizeExact(3),
					),
				},
				ExpectError: regexp.MustCompile("expected 3 elements for ListSizeExact check, got 2 elements"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_ListNestedBlock(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_ListNestedBlockPartial(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_ListNestedBlockElements(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("list_nested_block"),
						knownvalue.ListSizeExact(2),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Map(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_Map_KnownValueWrongType(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.ListExact([]knownvalue.Check{}),
					),
				},
				ExpectError: regexp.MustCompile(`expected \[\]any value for ListExact check, got: map\[string\]interface {}`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Map_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"key3": knownvalue.StringExact("value3"),
							"key4": knownvalue.StringExact("value4"),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`missing element key3 for MapExact check`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_MapPartial(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"key1": knownvalue.StringExact("value1"),
						}),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_MapPartial_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"key3": knownvalue.StringExact("value1"),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`missing element key3 for MapPartial check`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_MapElements(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.MapSizeExact(2),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_MapElements_WrongNum(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("map_attribute"),
						knownvalue.MapSizeExact(3),
					),
				},
				ExpectError: regexp.MustCompile("expected 3 elements for MapSizeExact check, got 2 elements"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Number(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						knownvalue.NumberExact(f),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Number_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						knownvalue.NumberExact(f),
					),
				},
				ExpectError: regexp.MustCompile("expected value 321 for NumberExact check, got: 123"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_Set(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_Set_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_attribute"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("value1"),
							knownvalue.StringExact("value3"),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`missing value value3 for SetExact check`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_SetPartial(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_attribute"),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringExact("value2"),
						}),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_SetPartial_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_attribute"),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringExact("value3"),
						}),
					),
				},
				ExpectError: regexp.MustCompile(`missing value value3 for SetPartial check`),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_SetElements(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_attribute"),
						knownvalue.SetSizeExact(2),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_SetNestedBlock(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_SetNestedBlock_Custom(t *testing.T) {
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
						set_nested_block_attribute = "string"
					}
					set_nested_block {
						set_nested_block_attribute = "girts"
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_nested_block"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.MapExact(map[string]knownvalue.Check{
								"set_nested_block_attribute": StringContains("str"),
							}),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"set_nested_block_attribute": StringContains("rts"),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_SetNestedBlockPartial(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
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
	})
}

func TestExpectKnownValue_CheckState_SetNestedBlockElements(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("set_nested_block"),
						knownvalue.SetSizeExact(2),
					),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_String(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						knownvalue.StringExact("str")),
				},
			},
		},
	})
}

func TestExpectKnownValue_CheckState_String_Custom(t *testing.T) {
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
					string_attribute = "string"
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						StringContains("tri")),
				},
			},
		},
	})
}

var _ knownvalue.Check = stringContains{}

type stringContains struct {
	value string
}

func (v stringContains) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringContains check, got: %T", other)
	}

	if !strings.Contains(otherVal, v.value) {
		return fmt.Errorf("expected string %q to contain %q for StringContains check", otherVal, v.value)
	}

	return nil
}

func (v stringContains) String() string {
	return v.value
}

func StringContains(value string) stringContains {
	return stringContains{
		value: value,
	}
}

func TestExpectKnownValue_CheckState_String_KnownValueWrongType(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						knownvalue.Bool(true)),
				},
				ExpectError: regexp.MustCompile("expected bool value for Bool check, got: string"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_String_KnownValueWrongValue(t *testing.T) {
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						knownvalue.StringExact("rts")),
				},
				ExpectError: regexp.MustCompile("expected value rts for StringExact check, got: str"),
			},
		},
	})
}

func TestExpectKnownValue_CheckState_UnknownAttributeType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		knownValue  knownvalue.Check
		req         statecheck.CheckStateRequest
		expectedErr error
	}{
		"unrecognised-type": {
			knownValue: knownvalue.Int64Exact(123),
			req: statecheck.CheckStateRequest{
				State: &tfjson.State{
					Values: &tfjson.StateValues{
						RootModule: &tfjson.StateModule{
							Resources: []*tfjson.StateResource{
								{
									Address: "example_resource.test",
									AttributeValues: map[string]any{
										"attribute": float32(123),
									},
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := statecheck.ExpectKnownValue("example_resource.test", tfjsonpath.New("attribute"), testCase.knownValue)

			resp := statecheck.CheckStateResponse{}

			e.CheckState(context.Background(), testCase.req, &resp)

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

func testProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test_resource": {
				CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
					d.SetId("test")

					err := d.Set("string_computed_attribute", "computed")
					if err != nil {
						return diag.Errorf("error setting string_computed_attribute: %s", err)
					}

					return nil
				},
				UpdateContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
					return nil
				},
				DeleteContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
					return nil
				},
				ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
					return nil
				},
				Schema: map[string]*schema.Schema{
					"bool_attribute": {
						Optional: true,
						Type:     schema.TypeBool,
					},
					"float_attribute": {
						Optional: true,
						Type:     schema.TypeFloat,
					},
					"int_attribute": {
						Optional: true,
						Type:     schema.TypeInt,
					},
					"list_attribute": {
						Type: schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					"list_nested_block": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"list_nested_block_attribute": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"map_attribute": {
						Type: schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					"set_attribute": {
						Type: schema.TypeSet,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					"set_nested_block": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"set_nested_block_attribute": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"set_nested_nested_block": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"set_nested_block": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"set_nested_block_attribute": {
												Type:     schema.TypeString,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"string_attribute": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"string_computed_attribute": {
						Computed: true,
						Type:     schema.TypeString,
					},
				},
			},
		},
	}
}
