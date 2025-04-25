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

	testinginterface "github.com/mitchellh/go-testing-interface"
)

func Test_SkipIfNotPrerelease_SkipTest_Stable(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotPrerelease(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIfNotPrerelease_RunTest_Alpha(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
		TFExactVersion: "1.9.0-alpha20240501",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotPrerelease(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIfNotPrerelease_RunTest_Beta1(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
		TFExactVersion: "1.8.0-beta1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotPrerelease(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}
func Test_SkipIfNotPrerelease_RunTest_RC(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
		TFExactVersion: "1.8.0-rc2",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIfNotPrerelease(),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}
