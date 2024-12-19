// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-plugin-testing/internal/testingiface"
)

func Test_SkipIfNotAlpha_SkipTest_Stable(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotAlpha(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIfNotAlpha_SkipTest_Beta1(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-beta1")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotAlpha(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}
func Test_SkipIfNotAlpha_SkipTest_RC(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc2")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotAlpha(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIfNotAlpha_RunTest_Alpha(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.9.0-alpha20240501")

	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTest(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipIfNotAlpha(),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}
