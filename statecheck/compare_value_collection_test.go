// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestCompareValueCollection_CheckState_Map_ValuesSame_DifferError(t *testing.T) {
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

				resource "test_resource" "two" {
					map_attribute = {
						"a": "str2",
						"b": "str3",
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("map_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str2 != str\nexpected values to be the same, but they differ: str3 != str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Map_ValuesSame(t *testing.T) {
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

				resource "test_resource" "two" {
					map_attribute = {
						"a": "str2",
						"b": test_resource.one.string_attribute,
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("map_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Map_ValuesDiffer_SameError(t *testing.T) {
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

				resource "test_resource" "two" {
					map_attribute = {
						"a": test_resource.one.string_attribute,
						"b": test_resource.one.string_attribute,
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("map_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: str == str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Map_ValuesDiffer(t *testing.T) {
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

				resource "test_resource" "two" {
					map_attribute = {
						"a": test_resource.one.string_attribute,
						"b": "str2",
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("map_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Set_ValuesSame_DifferError(t *testing.T) {
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

				resource "test_resource" "two" {
					set_attribute = [
						"str2",
						"str3"
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str3 != str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Set_ValuesSame(t *testing.T) {
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

				resource "test_resource" "two" {
					set_attribute = [
						"str2",
						test_resource.one.string_attribute
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Set_ValuesDiffer_SameError(t *testing.T) {
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

				resource "test_resource" "two" {
					set_attribute = [
						test_resource.one.string_attribute,
						test_resource.one.string_attribute,
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: str == str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Set_ValuesDiffer(t *testing.T) {
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

				resource "test_resource" "two" {
					set_attribute = [
						test_resource.one.string_attribute,
						"str2"
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedBlock_String_ValuesSame_DifferError(t *testing.T) {
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

				resource "test_resource" "two" {
					set_nested_block {		
						set_nested_block_attribute = "str1"	
					}
					set_nested_block {		
						set_nested_block_attribute = "str2"	
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_block"),
							tfjsonpath.New("set_nested_block_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str1 != str\nexpected values to be the same, but they differ: str2 != str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedBlock_String_ValuesSame(t *testing.T) {
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

				resource "test_resource" "two" {
					set_nested_block {		
						set_nested_block_attribute = "str1"	
					}
					set_nested_block {		
						set_nested_block_attribute = "str"	
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_block"),
							tfjsonpath.New("set_nested_block_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedBlock_SetNestedBlock_ValuesSame(t *testing.T) {
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
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str1"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str"	
						}
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
							tfjsonpath.New("set_nested_block"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_block"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedBlockBlock_SetNestedBlockBlock_ValuesSame(t *testing.T) {
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
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str1"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str2"	
						}
						set_nested_block {		
							set_nested_block_attribute = "st3"	
						}
					}
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str1"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str"	
						}
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesDiffer(t *testing.T) {
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
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str_x"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_y"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str_c"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_d"	
						}
					}
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str_a"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_b"	
						}
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
							tfjsonpath.New("set_nested_block"),
							tfjsonpath.New("set_nested_block_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block_attribute"),
						compare.ValuesDiffer(),
					),
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
							tfjsonpath.New("set_nested_block"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0),
						compare.ValuesDiffer(),
					),
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block"),
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSame(t *testing.T) {
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
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str_a"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_b"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str_c"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_d"	
						}
					}
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str_a"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_b"	
						}
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
							tfjsonpath.New("set_nested_block"),
							tfjsonpath.New("set_nested_block_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block_attribute"),
						compare.ValuesSame(),
					),
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
							tfjsonpath.New("set_nested_block"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0),
						compare.ValuesSame(),
					),
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("set_nested_nested_block"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}
