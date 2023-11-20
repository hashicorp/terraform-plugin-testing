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

func Test_ExpectUnknownOutputValue_StringAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				output "string_attribute" {
					value = time_static.one
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "string_attribute",
							AttributePath: tfjsonpath.New("rfc3339"),
						}),
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
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					list_attribute = ["value1", time_static.one.rfc3339]
				}

				output "list_attribute" {
					value = test_resource.two.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_attribute",
						}),
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
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

					resource "test_resource" "two" {
						set_attribute = ["value1", time_static.one.rfc3339]
					}

					output "set_attribute" {
						value = test_resource.two.set_attribute
					}
					`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "set_attribute",
						}),
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
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					map_attribute = {
						key1 = "value1",
						key2 = time_static.one.rfc3339
					}
				}

				output "map_attribute" {
					value = test_resource.two.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "map_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ListNestedBlock_Object(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = time_static.one.rfc3339
					}
				}

				output "object" {
					value = test_resource.two
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "object",
							AttributePath: tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute"),
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ListNestedBlock_ObjectBlock(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = time_static.one.rfc3339
					}
				}

				output "object_block" {
					value = test_resource.two.list_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "object_block",
							AttributePath: tfjsonpath.New(0).AtMapKey("list_nested_block_attribute"),
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ListNestedBlock_ObjectBlockIndex(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = time_static.one.rfc3339
					}
				}

				output "object_block_index" {
					value = test_resource.two.list_nested_block.0
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "object_block_index",
							AttributePath: tfjsonpath.New("list_nested_block_attribute"),
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_ListNestedBlock_ObjectBlockIndexAttribute(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = time_static.one.rfc3339
					}
				}

				output "object_block_index_attribute" {
					value = test_resource.two.list_nested_block.0.list_nested_block_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "object_block_index_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_SetNestedBlock_Object(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "time_static" "one" {}

				resource "test_resource" "two" {
					set_nested_block {
						set_nested_block_attribute = time_static.one.rfc3339
					}
				}

				output "object" {
					value = test_resource.two
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "object",
							AttributePath: tfjsonpath.New("set_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block_attribute"),
						}),
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
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "set_attribute",
						}),
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
						plancheck.ExpectUnknownOutputValue(plancheck.OutputValueParams{
							OutputAddress: "output_two",
						}),
					},
				},
				ExpectError: regexp.MustCompile(`output_two - Output not found in plan OutputChanges`),
			},
		},
	})
}
