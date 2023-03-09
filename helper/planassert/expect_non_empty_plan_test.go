package planassert

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_ExpectNonEmptyPlan(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 12
				}`,
				ConfigPlanAsserts: r.ConfigPlanAsserts{
					PreApply: []r.PlanAssert{
						ExpectNonEmptyPlan(),
					},
				},
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_Error(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 16
				}`,
				ConfigPlanAsserts: r.ConfigPlanAsserts{
					PreApply: []r.PlanAssert{
						ExpectNonEmptyPlan(),
					},
				},
				ExpectError: regexp.MustCompile(`expected a non-empty plan, but got an empty plan`),
			},
		},
	})
}
