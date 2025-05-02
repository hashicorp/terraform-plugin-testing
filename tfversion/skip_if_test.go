// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !testingtesting

package tfversion_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	testinginterface "github.com/mitchellh/go-testing-interface"
)

func Test_SkipIf_SkipTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.4.3")

	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
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
}

func Test_SkipIf_RunTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.1.0")

	r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
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
}

func Test_SkipIf_Prerelease_EqualCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// Pragmatic compromise that 1.8.0-rc1 prerelease is considered to
	// be equivalent to the 1.8.0 core version. This enables developers
	// to assert that prerelease versions are skipped with upcoming
	// core versions.
	//
	// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/303
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIf(version.Must(version.NewVersion("1.8.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIf_Prerelease_HigherCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.7.0-rc1")

	// The 1.7.0-rc1 prerelease should always be considered to be below the
	// 1.8.0 core version. This intentionally verifies that the logic does not
	// ignore the core version of the prerelease version when compared against
	// the core version of the check.
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIf(version.Must(version.NewVersion("1.8.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIf_Prerelease_HigherPrerelease(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.7.0-rc1")

	// The 1.7.0-rc1 prerelease should always be considered to be
	// below the 1.7.0-rc2 prerelease.
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIf(version.Must(version.NewVersion("1.7.0-rc2"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIf_Prerelease_LowerCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.7.0 core version.
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIf(version.Must(version.NewVersion("1.7.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipIf_Prerelease_LowerPrerelease(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.8.0-beta1 prerelease.
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipIf(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}
