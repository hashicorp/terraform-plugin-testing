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

func Test_ExpectNullOutputValue_StringAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
						return testProvider(), nil
					},
				},
				Config: `resource "test_resource" "test" {
				}

				output "string_attribute" {
					value = test_resource.test.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "string_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_StringAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
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

				output "string_attribute" {
					value = test_resource.test.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "string_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_StringAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
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

				output "string_attribute" {
					value = test_resource.test.string_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "string_attribute",
						}),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValue_ListAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
				}

				output "list_attribute" {
					value = test_resource.test.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_ListAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_attribute = null
				}

				output "list_attribute" {
					value = test_resource.test.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_ListAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_attribute = ["one", "two"]
				}

				output "list_attribute" {
					value = test_resource.test.list_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_attribute",
						}),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValue_SetAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
				}

				output "set_attribute" {
					value = test_resource.test.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "set_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_SetAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					set_attribute = null
				}

				output "set_attribute" {
					value = test_resource.test.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "set_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_SetAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					set_attribute = ["one", "two"]
				}

				output "set_attribute" {
					value = test_resource.test.set_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "set_attribute",
						}),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_MapAttribute_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
				}

				output "map_attribute" {
					value = test_resource.test.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "map_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_MapAttribute_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					map_attribute = null
				}

				output "map_attribute" {
					value = test_resource.test.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "map_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_MapAttribute_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					map_attribute = {
						"one": "str",
						"two": "str"
					}
				}

				output "map_attribute" {
					value = test_resource.test.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "map_attribute",
						}),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectUnknownOutputValue_MapAttribute_PartiallyNullConfig_ExpectError(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					map_attribute = {
						key1 = "value1",
						key2 = null
					}
				}

				output "map_attribute" {
					value = test_resource.test.map_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "map_attribute",
							AttributePath: tfjsonpath.New("key2"),
						}),
					},
				},
				ExpectError: regexp.MustCompile(`path not found: specified key key2 not found in map at key2`),
			},
		},
	})
}

func Test_ExpectNullOutputValue_ListNestedBlock_EmptyConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_nested_block {}
				}

				output "list_nested_block_attribute" {
					value = test_resource.test.list_nested_block.0.list_nested_block_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_nested_block_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_ListNestedBlock_NullConfig(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_nested_block {
						list_nested_block_attribute = null
					}
				}

				output "list_nested_block_attribute" {
					value = test_resource.test.list_nested_block.0.list_nested_block_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_nested_block_attribute",
						}),
					},
				},
			},
		},
	})
}

func Test_ExpectNullOutputValue_ListNestedBlock_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					list_nested_block {
						list_nested_block_attribute = "str"
					}
				}

				output "list_nested_block_attribute" {
					value = test_resource.test.list_nested_block.0.list_nested_block_attribute
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "list_nested_block_attribute",
						}),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}

func Test_ExpectNullOutputValue_SetNestedBlock_NullConfig_ExpectErrorNotNull(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProvider(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "test" {
					set_nested_block {
						set_nested_block_attribute = null
					}
				}

				output "set_nested_block" {
					value = test_resource.test.set_nested_block
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNullOutputValue(plancheck.OutputValueParams{
							OutputAddress: "set_nested_block",
							AttributePath: tfjsonpath.New(0).AtMapKey("set_nested_block_attribute"),
						}),
					},
				},
				ExpectError: regexp.MustCompile(`attribute at path is not null`),
			},
		},
	})
}
