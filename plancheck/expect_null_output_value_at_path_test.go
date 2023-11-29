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
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_ExpectNullOutputValueAtPath_StringAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return testProvider(), nil
					},
				},
				Config: `resource "test_resource" "test" {
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("string_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_StringAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return testProvider(), nil
					},
				},
				Config: `resource "test_resource" "test" {
					string_attribute = null
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("string_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_StringAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return testProvider(), nil
					},
				},
				Config: `resource "test_resource" "test" {
					string_attribute = "str"
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("string_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_ListAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("list_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_ListAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_attribute = null
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("list_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_ListAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_attribute = ["one", "two"]
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("list_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_SetAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("set_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_SetAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					set_attribute = null
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("set_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_SetAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					set_attribute = ["one", "two"]
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("set_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_MapAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("map_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_MapAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					map_attribute = null
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("map_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_MapAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					map_attribute = {
						"one": "str",
						"two": "str"
					}
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("map_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_MapAttribute_PartiallyNullConfig_ExpectError(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					map_attribute = {
						key1 = "value1",
						key2 = null
					}
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("map_attribute").AtMapKey("key2")),
					},
				},
				ExpectError: regexp.MustCompile(`path not found: specified key key2 not found in map at map_attribute.key2`),
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_ListNestedBlock_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_nested_block {}
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_ListNestedBlock_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_nested_block {
						list_nested_block_attribute = null
					}
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_ListNestedBlock_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValueAtPath_SetNestedBlock_NullConfig_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		// Prior to Terraform v1.3.0 a planned output is marked as fully unknown
		// if any attribute is unknown. The id attribute within the test provider
		// is unknown.
		// Reference: https://github.com/hashicorp/terraform/blob/v1.3/CHANGELOG.md#130-september-21-2022
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_3_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					set_nested_block {
						set_nested_block_attribute = null
					}
				}

				output "resource" {
					value = test_resource.test
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValueAtPath("resource", tfjsonpath.New("set_nested_block")),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}
