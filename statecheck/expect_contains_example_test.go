// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func ExampleExpectContains() {
	t := &testing.T{}
	t.Parallel()

	resource.Test(t, resource.TestCase{
		// Provider definition omitted.
		Steps: []resource.TestStep{
			{
				Config: `resource "test_resource" "one" {
		          string_attribute = "value1"
		        }

		        resource "test_resource" "two" {
		          set_attribute = [
		            test_resource.one.string_attribute,
		            "value2"
		          ]
		        }`,
				ConfigStateChecks: resource.ConfigStateChecks{
					statecheck.ExpectContains(
						"test_resource.two",
						tfjsonpath.New("set_attribute"),
						"test_resource.one",
						tfjsonpath.New("string_attribute"),
					),
				},
			},
		},
	})
}
