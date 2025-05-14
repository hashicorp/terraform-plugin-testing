// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func ExampleExpectIdentityValueMatchesStateAtPath() {
	// A typical test would accept *testing.T as a function parameter, for instance `func TestSomething(t *testing.T) { ... }`.
	t := &testing.T{}
	t.Parallel()

	resource.Test(t, resource.TestCase{
		// Resource identity support is only available in Terraform v1.12+
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		// Provider definition omitted. Assuming "test_resource":
		//  - Has an identity schema with an "identity_id" string attribute
		//  - Has a resource schema with an "state_id" string attribute
		Steps: []resource.TestStep{
			{
				Config: `resource "test_resource" "one" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					// The identity attribute at "identity_id" and state attribute at "state_id" must match
					statecheck.ExpectIdentityValueMatchesStateAtPath(
						"test_resource.one",
						tfjsonpath.New("identity_id"),
						tfjsonpath.New("state_id"),
					),
				},
			},
		},
	})
}
