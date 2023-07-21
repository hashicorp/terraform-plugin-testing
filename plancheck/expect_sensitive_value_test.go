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

func Test_ExpectSensitiveValue_SensitiveStringAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					sensitive_string_attribute = "test"
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.one",
							tfjsonpath.New("sensitive_string_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectSensitiveValue_SensitiveListAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					sensitive_list_attribute = ["value1"]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.one",
							tfjsonpath.New("sensitive_list_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectSensitiveValue_SensitiveSetAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					sensitive_set_attribute = ["value1"]
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.one",
							tfjsonpath.New("sensitive_set_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectSensitiveValue_SensitiveMapAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					sensitive_map_attribute = {
						key1 = "value1",
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.one",
							tfjsonpath.New("sensitive_map_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectSensitiveValue_ListNestedBlock_SensitiveAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					list_nested_block_sensitive_attribute {
						sensitive_list_nested_block_attribute = "sensitive-test"
						list_nested_block_attribute = "test"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.one",
							tfjsonpath.New("list_nested_block_sensitive_attribute").AtSliceIndex(0).
								AtMapKey("sensitive_list_nested_block_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectSensitiveValue_SetNestedBlock_SensitiveAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {
					set_nested_block_sensitive_attribute {
						sensitive_set_nested_block_attribute = "sensitive-test"
						set_nested_block_attribute = "test"
					}
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.one",
							tfjsonpath.New("set_nested_block_sensitive_attribute")),
					},
				},
			},
		},
	})
}

func Test_ExpectSensitiveValue_ExpectError_ResourceNotFound(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderSensitive(), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "test_resource" "one" {}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectSensitiveValue("test_resource.two", tfjsonpath.New("set_attribute")),
					},
				},
				ExpectError: regexp.MustCompile(`test_resource.two - Resource not found in plan ResourceChanges`),
			},
		},
	})
}

func testProviderSensitive() *schema.Provider {
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
					"sensitive_string_attribute": {
						Sensitive: true,
						Optional:  true,
						Type:      schema.TypeString,
					},
					"sensitive_list_attribute": {
						Sensitive: true,
						Type:      schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					"sensitive_set_attribute": {
						Sensitive: true,
						Type:      schema.TypeSet,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					"sensitive_map_attribute": {
						Sensitive: true,
						Type:      schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					"list_nested_block_sensitive_attribute": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"list_nested_block_attribute": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"sensitive_list_nested_block_attribute": {
									Sensitive: true,
									Type:      schema.TypeString,
									Optional:  true,
								},
							},
						},
					},
					"set_nested_block_sensitive_attribute": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"set_nested_block_attribute": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"sensitive_set_nested_block_attribute": {
									Sensitive: true,
									Type:      schema.TypeString,
									Optional:  true,
								},
							},
						},
					},
				},
			},
		},
	}
}
