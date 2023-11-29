// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck_test

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_ExpectNonEmptyPlan_OutputChanges_None(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_14_0), // outputs before 0.14 always show as created
		},
		// Avoid our own validation that requires at least one provider config.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				Config: `output "test" { value = "original" }`,
			},
			{
				Config: `output "test" { value = "new" }`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_OutputChanges_Error(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_14_0), // outputs before 0.14 always show as created
		},
		// Avoid our own validation that requires at least one provider config.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				Config: `output "test" { value = "original" }`,
			},
			{
				Config: `output "test" { value = "original" }`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				ExpectError: regexp.MustCompile(`expected a non-empty plan, but got an empty plan`),
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_ResourceChanges_None(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					triggers_replace = ["original"]
				}`,
			},
			{
				Config: `resource "terraform_data" "one" {
					triggers_replace = ["new"]
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
		},
	})
}

func Test_ExpectNonEmptyPlan_ResourceChanges_Error(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					triggers_replace = ["original"]
				}`,
			},
			{
				Config: `resource "terraform_data" "one" {
					triggers_replace = ["original"]
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				ExpectError: regexp.MustCompile(`expected a non-empty plan, but got an empty plan`),
			},
		},
	})
}
