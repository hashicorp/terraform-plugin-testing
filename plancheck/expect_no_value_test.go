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

func Test_ExpectNoValue_ValueNotExists(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) { //nolint:unparam // required signature
				return testProviderExpectNoValue(t), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_string"),
						),
					},
				},
			},
		},
	})
}

func TestExpectNoValue_ValueExists_String(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) {
				return testProviderExpectNoValue(t), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					attribute_string = "dummy"
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_string"),
						),
					},
				},
				ExpectError: regexp.MustCompile(
					`expected .* to have no value`,
				),
			},
		},
	})
}

func TestExpectNoValue_ValueExists_List(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) {
				return testProviderExpectNoValue(t), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					attribute_list = ["dummy"]
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_list"),
						),
					},
				},
				ExpectError: regexp.MustCompile(
					`expected .* to have no value`,
				),
			},
		},
	})
}

func TestExpectNoValue_ValueExists_Map(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) {
				return testProviderExpectNoValue(t), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					attribute_map = { "key" = "value" }
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_map"),
						),
					},
				},
				ExpectError: regexp.MustCompile(
					`expected .* to have no value`,
				),
			},
		},
	})
}

func TestExpectNoValue_ValueExists_Set(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) {
				return testProviderExpectNoValue(t), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					attribute_set = [ "value" ]
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_set"),
						),
					},
				},
				ExpectError: regexp.MustCompile(
					`expected .* to have no value`,
				),
			},
		},
	})
}

func TestExpectNoValue_ValueExists_UnknownValues(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) {
				return testProviderExpectNoValue(t), nil
			},
		},

		Steps: []r.TestStep{
			{
				Config: `
				resource "time_static" "this" {}

				resource "test_resource" "one" {
					attribute_string = time_static.this.rfc3339
				}
				`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_string"),
						),
					},
				},
				ExpectError: regexp.MustCompile(
					`expected .* to have no value`,
				),
			},
		},
	})
}

func TestExpectNoValue_ValueExists_NestedPath(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test": func() (*schema.Provider, error) {
				return testProviderExpectNoValue(t), nil
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "test_resource" "one" {
					attribute_set_nested {
						attribute_nested_string = "dummy"
					}
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNoValue(
							"test_resource.one",
							tfjsonpath.New("attribute_set_nested"),
						),
					},
				},
				ExpectError: regexp.MustCompile(
					`expected .* to have no value`,
				),
			},
		},
	})
}

func testProviderExpectNoValue(_ *testing.T) *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test_resource": {
				Schema: map[string]*schema.Schema{
					"attribute_string": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"attribute_list": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"attribute_map": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"attribute_set": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"attribute_set_nested": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"attribute_nested_string": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
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
			},
		},
	}
}
