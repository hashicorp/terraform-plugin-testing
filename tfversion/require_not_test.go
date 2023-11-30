// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	testinginterface "github.com/mitchellh/go-testing-interface"
)

func Test_RequireNot(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.4.3")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireNot(version.Must(version.NewVersion("1.1.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_RequireNot_Error(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
					return nil, nil
				},
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireNot(version.Must(version.NewVersion("1.1.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}
