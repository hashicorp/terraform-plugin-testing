// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func ExampleExpectIdentity() {
	// A typical test would accept *testing.T as a function parameter, for instance `func TestSomething(t *testing.T) { ... }`.
	t := &testing.T{}
	t.Parallel()

	resource.Test(t, resource.TestCase{
		// Resource identity support is only available in Terraform v1.12+
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		// Provider definition omitted. Assuming "test_resource" has an identity schema with "id" and "name" string attributes
		Steps: []resource.TestStep{
			{
				Config: `resource "test_resource" "one" {}`,
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectIdentity(
						"test_resource.one",
						map[string]knownvalue.Check{
							"id":   knownvalue.StringExact("id-123"),
							"name": knownvalue.StringExact("John Doe"),
						},
					),
				},
			},
		},
	})
}
