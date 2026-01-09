// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package statestore_test

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// MAINTAINER NOTE: While the StateStore mode is designed to test state store implementations, it can
// also be used to test existing Terraform core backends, which we do in this test file just for
// additional verification of the test mode.

func TestTerraformBackend_local(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		// MAINTAINER NOTE: Test steps won't run without a provider definition, so this is just
		// needed to pass validation, as we're just testing Terraform core itself.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				StateStore: true,
				Config: `
					terraform {
					  backend "local" {}
					}
				`,
			},
		},
	})
}

func TestTerraformBackend_local_empty_path_validation_error(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		// MAINTAINER NOTE: Test steps won't run without a provider definition, so this is just
		// needed to pass validation, as we're just testing Terraform core itself.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				StateStore: true,
				Config: `
					terraform {
					  backend "local" {
					    path = ""
					  }
					}
				`,
				ExpectError: regexp.MustCompile(`The "path" attribute value must not be empty.`),
			},
		},
	})
}
