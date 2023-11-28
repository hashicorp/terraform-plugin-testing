// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func Test_ExpectUnknownOutputValueAtPath_StringAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				output "resource" {
					value = terraform_data.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("output")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_ListAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				resource "test_resource" "two" {
					list_attribute = ["value1", terraform_data.one.output]
				}

				output "resource" {
					value = test_resource.two
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("list_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_SetAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
						input = "string"
					}

					resource "test_resource" "two" {
						set_attribute = ["value1", terraform_data.one.output]
					}

					output "resource" {
						value = test_resource.two
					}
					`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("set_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_MapAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				resource "test_resource" "two" {
					map_attribute = {
						key1 = "value1",
						key2 = terraform_data.one.output
					}
				}

				output "resource" {
					value = test_resource.two
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("map_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_ListNestedBlock_Resource(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = terraform_data.one.output
					}
				}

				output "resource" {
					value = test_resource.two
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_ListNestedBlock_ResourceBlocks(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = terraform_data.one.output
					}
				}

				output "resource_blocks" {
					value = test_resource.two.list_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource_blocks", tfjsonpath.New(0).AtMapKey("list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_ListNestedBlock_ObjectBlockIndex(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = terraform_data.one.output
					}
				}

				output "resource_blocks_index" {
					value = test_resource.two.list_nested_block.0
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource_blocks_index", tfjsonpath.New("list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_SetNestedBlock_Object(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					input = "string"
				}

				resource "test_resource" "two" {
					set_nested_block {
						set_nested_block_attribute = terraform_data.one.output
					}
				}

				output "resource" {
					value = test_resource.two
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("set_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_ExpectError_KnownValue(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					set_attribute = ["value1"]
				}

				output "resource" {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("resource", tfjsonpath.New("set_attribute").AtSliceIndex(0)),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is known`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValueAtPath_ExpectError_OutputNotFound(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					set_attribute = ["value1"]
				}

				output "output_one" {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValueAtPath("output_two", tfjsonpath.New("set_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`output_two - Output not found in plan OutputChanges`),
			},
		},
	})
}
