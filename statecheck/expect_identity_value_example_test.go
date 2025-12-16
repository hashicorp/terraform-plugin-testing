// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func ExampleExpectIdentityValue() {
	// A typical test would accept *testing.T as a function parameter, for instance `func TestSomething(t *testing.T) { ... }`.
	t := &testing.T{}
	t.Parallel()

	resource.Test(t, resource.TestCase{
		// Resource identity support is only available in Terraform v1.12+
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		// Provider definition omitted. Assuming "test_resource" has an identity schema with an "id" string attribute
		Steps: []resource.TestStep{
			{
				Config: `resource "test_resource" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue(
						"test_resource.one",
						tfjsonpath.New("id"),
						knownvalue.StringExact("id-123"),
					),
				},
			},
		},
	})
}
