// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"regexp"
	"testing"
)

func Test_RemoveState_Ok(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 2
				}`,
				Check: TestCheckResourceAttr("random_string.one", "length", "2"),
			},
			{
				RemoveState: []string{"random_string.one"},
				Check:       TestCheckNoResourceAttr("random_string.one", "length"),
			},
		},
	})
}

func Test_RemoveState_Error(t *testing.T) {
	t.Parallel()

	Test(t, TestCase{
		ExternalProviders: map[string]ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []TestStep{
			{
				Config: `resource "random_string" "one" {
					length = 2
				}`,
			},
			{
				RemoveState: []string{"resource.other"},
				ExpectError: regexp.MustCompile("Error: Invalid target address"),
			},
		},
	})
}
