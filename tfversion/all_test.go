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
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_All_RunTest(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.1.0",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.Any(
				tfversion.All(
					tfversion.RequireNot(version.Must(version.NewVersion("0.15.0"))),  //returns nil
					tfversion.SkipIf(version.Must(version.NewVersion("1.2.0"))),       //returns nil
					tfversion.RequireBelow(version.Must(version.NewVersion("1.2.0"))), //returns nil
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
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.0.7",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.Any(
				tfversion.All(
					tfversion.RequireNot(version.Must(version.NewVersion("0.15.0"))),  //returns nil
					tfversion.SkipBelow(version.Must(version.NewVersion("1.2.0"))),    //returns skip
					tfversion.SkipIf(version.Must(version.NewVersion("1.0.7"))),       //returns skip
					tfversion.RequireBelow(version.Must(version.NewVersion("1.2.0"))), //returns nil
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
	t.Parallel()

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			TFExactVersion: "1.0.7",
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
					return nil, nil
				},
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.Any(
					tfversion.All(
						tfversion.RequireNot(version.Must(version.NewVersion("1.1.0"))),   //returns error
						tfversion.SkipIf(version.Must(version.NewVersion("1.1.0"))),       //returns skip
						tfversion.RequireAbove(version.Must(version.NewVersion("1.2.0"))), //returns nil
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
