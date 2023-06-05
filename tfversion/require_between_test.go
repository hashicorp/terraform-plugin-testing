// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	testinginterface "github.com/mitchellh/go-testing-interface"
)

func Test_RequireBetween(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.2.0")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireBetween(version.Must(version.NewVersion("1.2.0")), version.Must(version.NewVersion("1.3.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `
					terraform {
  						experiments = [module_variable_optional_attrs]
					}

					locals {
						ex_var = "hello"
					}

					output "example" {
						value = "output"
						precondition {
							condition = local.ex_var != "hi"
							error_message = "precondition_error"
						}
					}`,
			},
		},
	})
}

func Test_RequireBetween_Error_BelowMin(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
					return nil, nil
				},
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireBetween(version.Must(version.NewVersion("1.2.0")), version.Must(version.NewVersion("1.3.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_RequireBetween_Error_EqToMax(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.3.0")

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
					return nil, nil
				},
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireBetween(version.Must(version.NewVersion("1.2.0")), version.Must(version.NewVersion("1.3.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}
