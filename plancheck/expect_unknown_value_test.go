// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func Test_ExpectUnknownValue_StringAttribute(t *testing.T) {
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
					string_attribute = time_static.one.rfc3339
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two", tfjsonpath.New("string_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownValue_ListAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two", tfjsonpath.New("list_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownValue_SetAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two", tfjsonpath.New("set_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownValue_MapAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two", tfjsonpath.New("map_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownValue_ListNestedBlock(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two",
							tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownValue_SetNestedBlock(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two",
							tfjsonpath.New("set_nested_block").AtSliceIndex(0).AtMapKey("set_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_KnownValue_PathNotFound(t *testing.T) {
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
				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = "value 1"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two",
							tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("not_correct_attr")),
					},
				},
				ExpectError: regexp.MustCompile(`path not found: specified key not_correct_attr not found in map at list_nested_block.0.not_correct_attr`),
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_KnownValue_ListAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.one", tfjsonpath.New("list_attribute").AtSliceIndex(0)),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at "list_attribute.0", but found known value: "value1"`),
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_KnownValue_StringAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.one", tfjsonpath.New("string_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at "string_attribute", but found known value: "hello world!"`),
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_KnownValue_BoolAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.one", tfjsonpath.New("bool_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at "bool_attribute", but found known value: "true"`),
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_KnownValue_FloatAttribute(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.one", tfjsonpath.New("float_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at "float_attribute", but found known value: "1.234"`),
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_KnownValue_ListNestedBlock(t *testing.T) {
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
				resource "test_resource" "two" {
					list_nested_block {
						list_nested_block_attribute = "value 1"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two",
							tfjsonpath.New("list_nested_block").AtSliceIndex(0).AtMapKey("list_nested_block_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`Expected unknown value at "list_nested_block.0.list_nested_block_attribute", but found known value: "value 1"`),
			},
		},
	})
}

func Test_ExpectUnknownValue_ExpectError_ResourceNotFound(t *testing.T) {
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
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("test_resource.two", tfjsonpath.New("set_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`test_resource.two - Resource not found in plan ResourceChanges`),
			},
		},
	})
}

func testProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test_resource": {
				CreateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
					d.SetId("test")
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
					"string_attribute": {
						Optional: true,
						Type:     schema.TypeString,
					},
				},
			},
		},
	}
}
