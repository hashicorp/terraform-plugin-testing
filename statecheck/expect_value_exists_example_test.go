// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func ExampleExpectValueExists() {
	t := &testing.T{}
	t.Parallel()

	resource.Test(t, resource.TestCase{
		// Provider definition omitted.
		Steps: []resource.TestStep{
			{
				Config: `resource "test_resource" "one" {
		          bool_attribute = true
		        }
		        `,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectValueExists(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
					),
				},
			},
		},
	})
}
