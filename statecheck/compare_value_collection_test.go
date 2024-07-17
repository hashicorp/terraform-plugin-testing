// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestCompareValueCollection_CheckState_Bool_Error_NotCollection(t *testing.T) {
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

				resource "test_resource" "two" {
					bool_attribute = true
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("bool_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two.bool_attribute is not a collection type: bool"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Float_Error_NotCollection(t *testing.T) {
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
					float_attribute = 1.234
				}

				resource "test_resource" "two" {
					float_attribute = 1.234
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("float_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("float_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two.float_attribute is not a collection type: json.Number"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_Int_Error_NotCollection(t *testing.T) {
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
					int_attribute = 1234
				}

				resource "test_resource" "two" {
					int_attribute = 1234
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("int_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("int_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two.int_attribute is not a collection type: json.Number"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_List_ValuesSame_ErrorDiffer(t *testing.T) {
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
					list_attribute = [
						"str2",
						"str3",
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_attribute"),
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

func TestCompareValueCollection_CheckState_EmptyCollectionPath(t *testing.T) {
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
					list_attribute = [
						"str2",
						"str",
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						// Empty path is invalid
						[]tfjsonpath.Path{},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two - No collection path was provided"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_List_ValuesSame(t *testing.T) {
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
					list_attribute = [
						"str2",
						"str",
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_attribute"),
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

func TestCompareValueCollection_CheckState_List_ValuesDiffer_ErrorSame(t *testing.T) {
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
					list_attribute = [
						"str",
						"str",
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_attribute"),
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

func TestCompareValueCollection_CheckState_List_ValuesDiffer(t *testing.T) {
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
					list_attribute = [
						"str",
						"str2",
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_attribute"),
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

func TestCompareValueCollection_CheckState_ListNestedBlock_ValuesSame_ErrorDiffer(t *testing.T) {
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
					list_nested_block {		
						list_nested_block_attribute = "str2"	
					}
					list_nested_block {		
						list_nested_block_attribute = "str3"	
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_nested_block"),
							tfjsonpath.New("list_nested_block_attribute"),
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

func TestCompareValueCollection_CheckState_ListNestedBlock_ValuesSame(t *testing.T) {
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
					list_nested_block {		
						list_nested_block_attribute = "str2"	
					}
					list_nested_block {		
						list_nested_block_attribute = "str"	
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_nested_block"),
							tfjsonpath.New("list_nested_block_attribute"),
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

func TestCompareValueCollection_CheckState_ListNestedBlock_ValuesDiffer_ErrorSame(t *testing.T) {
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
					list_nested_block {		
						list_nested_block_attribute = "str"	
					}
					list_nested_block {		
						list_nested_block_attribute = "str"	
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_nested_block"),
							tfjsonpath.New("list_nested_block_attribute"),
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

func TestCompareValueCollection_CheckState_ListNestedBlock_ValuesDiffer(t *testing.T) {
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
					list_nested_block {		
						list_nested_block_attribute = "str2"	
					}
					list_nested_block {		
						list_nested_block_attribute = "str3"	
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("list_nested_block"),
							tfjsonpath.New("list_nested_block_attribute"),
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

func TestCompareValueCollection_CheckState_Map_ValuesSame_ErrorDiffer(t *testing.T) {
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
						"b": "str",
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

func TestCompareValueCollection_CheckState_Map_ValuesDiffer_ErrorSame(t *testing.T) {
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
						"a": "str",
						"b": "str",
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
						"a": "str",
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

func TestCompareValueCollection_CheckState_Set_ValuesSame_ErrorDiffer(t *testing.T) {
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
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str2 != str\nexpected values to be the same, but they differ: str3 != str"),
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
						"str"
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

func TestCompareValueCollection_CheckState_Set_ValuesDiffer_ErrorSame(t *testing.T) {
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
						"str",
						"str",
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
						"str",
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

func TestCompareValueCollection_CheckState_SetNestedBlock_ValuesSame_ErrorDiffer(t *testing.T) {
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
						set_nested_block_attribute = "str2"	
					}
					set_nested_block {		
						set_nested_block_attribute = "str3"	
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
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str2 != str\nexpected values to be the same, but they differ: str3 != str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedBlock_ValuesSame(t *testing.T) {
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
						set_nested_block_attribute = "str2"	
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

func TestCompareValueCollection_CheckState_SetNestedBlock_ValuesDiffer_ErrorSame(t *testing.T) {
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
						set_nested_block_attribute = "str"	
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
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: str == str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedBlock_ValuesDiffer(t *testing.T) {
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
						set_nested_block_attribute = "str2"	
					}
					set_nested_block {		
						set_nested_block_attribute = "str3"	
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
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedNestedBlock_ValuesDiffer_ErrorSameAttribute(t *testing.T) {
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
							set_nested_block_attribute = "str"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
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
							tfjsonpath.New("set_nested_block_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block_attribute"),
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to differ, but they are the same: str == str"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedNestedBlock_ValuesDiffer_ErrorSameNestedBlock(t *testing.T) {
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
							set_nested_block_attribute = "str"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
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
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0),
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile(`expected values to differ, but they are the same: map\[set_nested_block_attribute:str\] == map\[set_nested_block_attribute:str\]`),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedNestedBlock_ValuesDiffer_ErrorSameNestedNestedBlock(t *testing.T) {
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
							set_nested_block_attribute = "str"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
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
						compare.ValuesDiffer(),
					),
				},
				ExpectError: regexp.MustCompile(`expected values to differ, but they are the same: \[map\[set_nested_block:\[map\[set_nested_block_attribute:str\]\]\]\] == \[map\[set_nested_block:\[map\[set_nested_block_attribute:str\]\]\]\]`),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedNestedBlock_ValuesDifferAttribute(t *testing.T) {
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
							set_nested_block_attribute = "str2"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str3"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str4"	
						}
					}
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str5"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str6"	
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
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedNestedBlock_ValuesDifferNestedBlock(t *testing.T) {
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
							set_nested_block_attribute = "str2"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str3"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str4"	
						}
					}
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str5"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str6"	
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
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0),
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNestedNestedBlock_ValuesDifferNestedNestedBlock(t *testing.T) {
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
							set_nested_block_attribute = "str2"	
						}
					}
				}

				resource "test_resource" "two" {
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str3"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str4"	
						}
					}
					set_nested_nested_block {		
						set_nested_block {		
							set_nested_block_attribute = "str5"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str6"	
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
						compare.ValuesDiffer(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSame_ErrorAttribute(t *testing.T) {
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
							set_nested_block_attribute = "str_e"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_f"	
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
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str_c != str_a\nexpected values to be the same, but they differ: str_d != str_a\nexpected values to be the same, but they differ: str_e != str_a\nexpected values to be the same, but they differ: str_f != str_a"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSame_ErrorNestedBlock(t *testing.T) {
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
							set_nested_block_attribute = "str_e"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_f"	
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
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile(`expected values to be the same, but they differ: map\[set_nested_block_attribute:str_c\] != map\[set_nested_block_attribute:str_a\]\nexpected values to be the same, but they differ: map\[set_nested_block_attribute:str_d\] != map\[set_nested_block_attribute:str_a\]\nexpected values to be the same, but they differ: map\[set_nested_block_attribute:str_e\] != map\[set_nested_block_attribute:str_a\]\nexpected values to be the same, but they differ: map\[set_nested_block_attribute:str_f\] != map\[set_nested_block_attribute:str_a\]`),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSame_ErrorNestedNestedBlock(t *testing.T) {
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
							set_nested_block_attribute = "str_e"	
						}
						set_nested_block {		
							set_nested_block_attribute = "str_f"	
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
				ExpectError: regexp.MustCompile(`expected values to be the same, but they differ: \[map\[set_nested_block:\[map\[set_nested_block_attribute:str_c\] map\[set_nested_block_attribute:str_d\]\]\]\] != \[map\[set_nested_block:\[map\[set_nested_block_attribute:str_a\] map\[set_nested_block_attribute:str_b\]\]\]\]\nexpected values to be the same, but they differ: \[map\[set_nested_block:\[map\[set_nested_block_attribute:str_e\] map\[set_nested_block_attribute:str_f\]\]\]\] != \[map\[set_nested_block:\[map\[set_nested_block_attribute:str_a\] map\[set_nested_block_attribute:str_b\]\]\]\]`),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSameAttribute(t *testing.T) {
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
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSameNestedBlock(t *testing.T) {
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
						},
						"test_resource.one",
						tfjsonpath.New("set_nested_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block").AtSliceIndex(0),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_SetNested_ValuesSameNestedNestedBlock(t *testing.T) {
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

func TestCompareValueCollection_CheckState_String_Error_NotCollection(t *testing.T) {
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
					string_attribute = "str"
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("string_attribute"),
						},
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("test_resource.two.string_attribute is not a collection type: string"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_ListNestedAttribute_ValuesSame(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "str_attr",
											Type:     tftypes.String,
											Optional: true,
										},
										{
											Name: "nested_attr",
											NestedType: &tfprotov6.SchemaObject{
												Attributes: []*tfprotov6.SchemaAttribute{
													{
														Name:     "str_attr",
														Type:     tftypes.String,
														Optional: true,
													},
												},
												Nesting: tfprotov6.SchemaObjectNestingModeList,
											},
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					str_attr = "str2"
				}
				resource "test_resource" "two" {
					nested_attr = [
						{
							str_attr = "str1"
						},
						{
							str_attr = "str2"
						}
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("nested_attr"),
							tfjsonpath.New("str_attr"),
						},
						"test_resource.one",
						tfjsonpath.New("str_attr"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_ListNestedAttribute_ValuesSame_ErrorDiff(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "str_attr",
											Type:     tftypes.String,
											Optional: true,
										},
										{
											Name: "nested_attr",
											NestedType: &tfprotov6.SchemaObject{
												Attributes: []*tfprotov6.SchemaAttribute{
													{
														Name:     "str_attr",
														Type:     tftypes.String,
														Optional: true,
													},
												},
												Nesting: tfprotov6.SchemaObjectNestingModeList,
											},
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					str_attr = "str1"
				}
				resource "test_resource" "two" {
					nested_attr = [
						{
							str_attr = "str2"
						},
						{
							str_attr = "str3"
						}
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("nested_attr"),
							tfjsonpath.New("str_attr"),
						},
						"test_resource.one",
						tfjsonpath.New("str_attr"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str2 != str1\nexpected values to be the same, but they differ: str3 != str1"),
			},
		},
	})
}

func TestCompareValueCollection_CheckState_DoubleListNestedAttribute_ValuesSame(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "str_attr",
											Type:     tftypes.String,
											Optional: true,
										},
										{
											Name: "nested_attr",
											NestedType: &tfprotov6.SchemaObject{
												Attributes: []*tfprotov6.SchemaAttribute{
													{
														Name: "double_nested_attr",
														NestedType: &tfprotov6.SchemaObject{
															Attributes: []*tfprotov6.SchemaAttribute{
																{
																	Name:     "str_attr",
																	Type:     tftypes.String,
																	Optional: true,
																},
															},
															Nesting: tfprotov6.SchemaObjectNestingModeSingle,
														},
														Optional: true,
													},
												},
												Nesting: tfprotov6.SchemaObjectNestingModeList,
											},
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					str_attr = "str2"
				}
				resource "test_resource" "two" {
					nested_attr = [
						{
							double_nested_attr = {
								str_attr = "str1"
							}
						},
						{
							double_nested_attr = {
								str_attr = "str2"
							}
						}
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("nested_attr"),
							tfjsonpath.New("double_nested_attr"),
							tfjsonpath.New("str_attr"),
						},
						"test_resource.one",
						tfjsonpath.New("str_attr"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestCompareValueCollection_CheckState_DoubleListNestedAttribute_ValuesSame_ErrorDiff(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"test_resource": {
						SchemaResponse: &resource.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										{
											Name:     "str_attr",
											Type:     tftypes.String,
											Optional: true,
										},
										{
											Name: "nested_attr",
											NestedType: &tfprotov6.SchemaObject{
												Attributes: []*tfprotov6.SchemaAttribute{
													{
														Name: "double_nested_attr",
														NestedType: &tfprotov6.SchemaObject{
															Attributes: []*tfprotov6.SchemaAttribute{
																{
																	Name:     "str_attr",
																	Type:     tftypes.String,
																	Optional: true,
																},
															},
															Nesting: tfprotov6.SchemaObjectNestingModeSingle,
														},
														Optional: true,
													},
												},
												Nesting: tfprotov6.SchemaObjectNestingModeList,
											},
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					str_attr = "str1"
				}
				resource "test_resource" "two" {
					nested_attr = [
						{
							double_nested_attr = {
								str_attr = "str2"
							}
						},
						{
							double_nested_attr = {
								str_attr = "str3"
							}
						}
					]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValueCollection(
						"test_resource.two",
						[]tfjsonpath.Path{
							tfjsonpath.New("nested_attr"),
							tfjsonpath.New("double_nested_attr"),
							tfjsonpath.New("str_attr"),
						},
						"test_resource.one",
						tfjsonpath.New("str_attr"),
						compare.ValuesSame(),
					),
				},
				ExpectError: regexp.MustCompile("expected values to be the same, but they differ: str2 != str1\nexpected values to be the same, but they differ: str3 != str1"),
			},
		},
	})
}
