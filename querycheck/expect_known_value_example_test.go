// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func ExampleExpectKnownValue() {
	// A typical test would accept *testing.T as a function parameter, for instance `func TestSomething(t *testing.T) { ... }`.
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
				ConfigQueryChecks: []querycheck.QueryCheck{
					querycheck.ExpectKnownValue(
						"test_resource.one",
						tfjsonpath.New("bool_attribute"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}
