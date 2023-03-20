package plancheck_test

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func Test_ExpectEmptyPlan(t *testing.T) {
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
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func Test_ExpectEmptyPlan_Error(t *testing.T) {
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
				}
				resource "random_string" "two" {
					length = 16
				}
				resource "random_string" "three" {
					length = 16
				}`,
			},
			{
				Config: `resource "random_string" "one" {
					length = 12
				}
				resource "random_string" "two" {
					length = 16
				}
				resource "random_string" "three" {
					length = 12
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ExpectError: regexp.MustCompile(`.*?(random_string.one has planned action\(s\): \[delete create\])\n.*?(random_string.three has planned action\(s\): \[delete create\])`),
			},
		},
	})
}
