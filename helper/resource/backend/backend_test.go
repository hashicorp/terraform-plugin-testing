// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backend_test

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper equivalent: TestBackendConfig
// https://github.com/hashicorp/terraform/blob/643266dc90523ae794b586ad42e8b30864b61aaa/internal/backend/testing.go#L23-L25
//
// 1. Validates configures the backend (i.e. terraform init)
func TestBackend_s3_no_region_error(t *testing.T) {
	t.Setenv("AWS_REGION", "")

	r.UnitTest(t, r.TestCase{
		// TODO: Plugin testing won't let you run a test without any provider defintitions, so this is temporary since
		// we're just testing Terraform core itself.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				BackendSmokeTest: true,
				Config: `
					terraform {
					  backend "s3" {
					    bucket = "test-pss-backend"
					    key    = "av-terraform-state"
					  }
					}
				`,
				ExpectError: regexp.MustCompile(`The \"region\" attribute or the \"AWS_REGION\" or \"AWS_DEFAULT_REGION\"`),
			},
		},
	})
}

// Helper equivalent: TestBackendStates (TODO: doesn't work for backends that don't support workspaces ATM, i.e. just http)
// https://github.com/hashicorp/terraform/blob/643266dc90523ae794b586ad42e8b30864b61aaa/internal/backend/testing.go#L74-L78
func TestBackend_s3_state(t *testing.T) {
	t.Parallel()
	// TODO: currently I'm not actually setting this test up and clearing it properly because I'm too lazy to write 10 lines of code :)

	r.UnitTest(t, r.TestCase{
		// TODO: Plugin testing won't let you run a test without any provider defintitions, so this is temporary since
		// we're just testing Terraform core itself.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				BackendSmokeTest: true,
				Config: `
					terraform {
					  backend "s3" {
					    bucket = "test-pss-backend"
					    key    = "av-terraform-state"
						region = "us-east-1"
					  }
					}
				`,
			},
		},
	})
}

// Helper equivalent: TestBackendStates (TODO: doesn't work for backends that don't support workspaces ATM, i.e. just http)
// https://github.com/hashicorp/terraform/blob/643266dc90523ae794b586ad42e8b30864b61aaa/internal/backend/testing.go#L74-L78
func TestBackend_s3_lock(t *testing.T) {
	t.Parallel()
	// TODO: currently I'm not actually setting this test up and clearing it properly because I'm too lazy to write 10 lines of code :)

	r.UnitTest(t, r.TestCase{
		// TODO: Plugin testing won't let you run a test without any provider defintitions, so this is temporary since
		// we're just testing Terraform core itself.
		ExternalProviders: map[string]r.ExternalProvider{
			"terraform": {Source: "terraform.io/builtin/terraform"},
		},
		Steps: []r.TestStep{
			{
				BackendLockTest: true,
				Config: `
					terraform {
					  backend "s3" {
					    bucket = "test-pss-backend"
					    key    = "av-terraform-state"
						region = "us-east-1"
						use_lockfile = true
					  }
					}
				`,
			},
		},
	})
}
