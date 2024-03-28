// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testingiface"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_SkipIf_SkipTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.4.3")

	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTest(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipIf(version.Must(version.NewVersion("1.4.3"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipIf_RunTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTest(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipIf(version.Must(version.NewVersion("1.2.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}
