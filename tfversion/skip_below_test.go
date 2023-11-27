// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_SkipBelow_SkipTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.0.7")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.1.0"))),
		},
		Steps: []r.TestStep{
			{
				//nullable argument only available in TF v1.1.0+
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}

func Test_SkipBelow_RunTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.1.0"))),
		},
		Steps: []r.TestStep{
			{
				//nullable argument only available in TF v1.1.0+
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}
