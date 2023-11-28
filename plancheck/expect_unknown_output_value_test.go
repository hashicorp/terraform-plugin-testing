// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func Test_ExpectUnknownOutputValue_StringAttribute(t *testing.T) {
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

				output "string_attribute" {
					value = terraform_data.one.output
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("string_attribute"),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ListAttribute(t *testing.T) {
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

				output "list_attribute" {
					value = test_resource.two.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("list_attribute"),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_SetAttribute(t *testing.T) {
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

					output "set_attribute" {
						value = test_resource.two.set_attribute
					}
					`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("set_attribute"),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_MapAttribute(t *testing.T) {
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

				output "map_attribute" {
					value = test_resource.two.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("map_attribute"),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ListNestedBlock(t *testing.T) {
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

				output "list_nested_block_attribute" {
					value = test_resource.two.list_nested_block.0.list_nested_block_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("list_nested_block_attribute"),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ExpectError_KnownValue(t *testing.T) {
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

				output "set_attribute" {
					value = test_resource.one.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("set_attribute"),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is known`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ExpectError_OutputNotFound(t *testing.T) {
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
				resource "test_resource" "one" {}

				output "output_one" {
					value = test_resource.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("output_two"),
					},
				},
				ExpectError: regexp.MustCompile(`output_two - Output not found in plan OutputChanges`),
			},
		},
	})
}
