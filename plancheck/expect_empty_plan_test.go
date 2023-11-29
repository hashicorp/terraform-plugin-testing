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

func Test_ExpectEmptyPlan_OutputChanges_None(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				Config: `output "test" { value = "original" }`,
			},
			{
				Config: `output "test" { value = "original" }`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func Test_ExpectEmptyPlan_OutputChanges_Error(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		Steps: []r.TestStep{
			{
				Config: `output "test" { value = "original" }`,
			},
			{
				Config: `output "test" { value = "new" }`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ExpectError: regexp.MustCompile(`output \"test\" has planned action\(s\): \[update\]`),
			},
		},
	})
}

func Test_ExpectEmptyPlan_ResourceChanges_None(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "test" {}`,
			},
			{
				Config: `resource "terraform_data" "test" {}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func Test_ExpectEmptyPlan_ResourceChanges_Error(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		Steps: []r.TestStep{
			{
				Config: `resource "terraform_data" "one" {
					triggers_replace = ["original"]
				}
				resource "terraform_data" "two" {
					triggers_replace = ["original"]
				}
				resource "terraform_data" "three" {
					triggers_replace = ["original"]
				}`,
			},
			{
				Config: `resource "terraform_data" "one" {
					triggers_replace = ["new"]
				}
				resource "terraform_data" "two" {
					triggers_replace = ["original"]
				}
				resource "terraform_data" "three" {
					triggers_replace = ["new"]
				}`,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ExpectError: regexp.MustCompile(`.*?(terraform_data.one has planned action\(s\): \[delete create\])\n.*?(terraform_data.three has planned action\(s\): \[delete create\])`),
			},
		},
	})
}
