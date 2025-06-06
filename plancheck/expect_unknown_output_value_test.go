// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_ExpectUnknownOutputValue_StringAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"data": {
				Source: "terraform.io/builtin/terraform",
			},
		},
		// The terraform_data resource is not available prior to Terraform v1.4.0
		// Reference: https://github.com/hashicorp/terraform/blob/v1.4/CHANGELOG.md#140-march-08-2023
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.4.0"))),
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
		// The terraform_data resource is not available prior to Terraform v1.4.0
		// Reference: https://github.com/hashicorp/terraform/blob/v1.4/CHANGELOG.md#140-march-08-2023
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.4.0"))),
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
		// The terraform_data resource is not available prior to Terraform v1.4.0
		// Reference: https://github.com/hashicorp/terraform/blob/v1.4/CHANGELOG.md#140-march-08-2023
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.4.0"))),
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
		// The terraform_data resource is not available prior to Terraform v1.4.0
		// Reference: https://github.com/hashicorp/terraform/blob/v1.4/CHANGELOG.md#140-march-08-2023
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.4.0"))),
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
		// The terraform_data resource is not available prior to Terraform v1.4.0
		// Reference: https://github.com/hashicorp/terraform/blob/v1.4/CHANGELOG.md#140-march-08-2023
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.4.0"))),
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

func Test_ExpectUnknownOutputValue_ExpectError_KnownValue_ListAttribute(t *testing.T) {
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
					list_attribute = ["value1"]
				}

				output "list_attribute" {
					value = test_resource.one.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("list_attribute"),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at output "list_attribute", but found known value: "\[value1\]"`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ExpectError_KnownValue_StringAttribute(t *testing.T) {
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
					string_attribute = "hello world!"
				}

				output "string_attribute" {
					value = test_resource.one.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("string_attribute"),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at output "string_attribute", but found known value: "hello world!"`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ExpectError_KnownValue_BoolAttribute(t *testing.T) {
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
					bool_attribute = true
				}

				output "bool_attribute" {
					value = test_resource.one.bool_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("bool_attribute"),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at output "bool_attribute", but found known value: "true"`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ExpectError_KnownValue_FloatAttribute(t *testing.T) {
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
					float_attribute = 1.234
				}

				output "float_attribute" {
					value = test_resource.one.float_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("float_attribute"),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at output "float_attribute", but found known value: "1.234"`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ExpectError_KnownValue_ListNestedBlock(t *testing.T) {
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
					list_nested_block {
						list_nested_block_attribute = "value 1"
					}
				}

				output "list_nested_block" {
					value = test_resource.one.list_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue("list_nested_block"),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at output "list_nested_block", but found known value: "\[map\[list_nested_block_attribute:value 1\]\]"`),
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
