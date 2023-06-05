// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testinginterface "github.com/mitchellh/go-testing-interface"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_All_RunTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.Any(
				tfversion.All(
					tfversion.RequireNot(version.Must(version.NewVersion("0.15.0"))),
					tfversion.SkipIf(version.Must(version.NewVersion("1.2.0"))),
					tfversion.RequireBelow(version.Must(version.NewVersion("1.2.0"))),
				),
			),
		},
		Steps: []r.TestStep{
			{
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}

func Test_All_SkipTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.0.7")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.Any(
				tfversion.All(
					tfversion.RequireNot(version.Must(version.NewVersion("0.15.0"))),
					tfversion.SkipBelow(version.Must(version.NewVersion("1.2.0"))),
					tfversion.SkipIf(version.Must(version.NewVersion("1.0.7"))),
					tfversion.RequireBelow(version.Must(version.NewVersion("1.2.0"))),
				),
			),
		},
		Steps: []r.TestStep{
			{
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}

func Test_All_Error(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
					return nil, nil
				},
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.Any(
					tfversion.All(
						tfversion.RequireNot(version.Must(version.NewVersion("1.1.0"))),
						tfversion.SkipIf(version.Must(version.NewVersion("1.1.0"))),
						tfversion.RequireAbove(version.Must(version.NewVersion("1.2.0"))),
					),
				),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}
