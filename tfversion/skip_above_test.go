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

func Test_SkipAbove_SkipTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.3.0")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipAbove(version.Must(version.NewVersion("1.2.9"))),
		},
		Steps: []r.TestStep{
			{
				//module_variable_optional_attrs experiment is deprecated in TF v1.3.0
				Config: `
					terraform {
  						experiments = [module_variable_optional_attrs]
					}
				`,
			},
		},
	})
}

func Test_SkipAbove_RunTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.2.9")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipAbove(version.Must(version.NewVersion("1.2.9"))),
		},
		Steps: []r.TestStep{
			{
				//module_variable_optional_attrs experiment is deprecated in TF v1.3.0
				Config: `
					terraform {
  						experiments = [module_variable_optional_attrs]
					}
				`,
			},
		},
	})
}
