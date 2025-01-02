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

func Test_SkipBetween_SkipTest(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.2.0")

	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
					return nil, nil
				},
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.2.0")), version.Must(version.NewVersion("1.3.0"))),
			},
			Steps: []r.TestStep{
				{
					//module_variable_optional_attrs experiment is deprecated in TF v1.3.0
					//precondition block is only available in TF v1.2.0+
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
	})
}

func Test_SkipBetween_RunTest_AboveMax(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.3.0")

	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.2.0")), version.Must(version.NewVersion("1.3.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_RunTest_EqToMin(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.2.0")

	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.2.0")), version.Must(version.NewVersion("1.3.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MaxEqualCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// Pragmatic compromise that 1.8.0-rc1 prerelease is considered to
	// be equivalent to the 1.8.0 core version. This enables developers
	// to assert that prerelease versions are ran with upcoming
	// core versions.
	//
	// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/303
	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.7.0")), version.Must(version.NewVersion("1.8.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MinEqualCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// Pragmatic compromise that 1.8.0-rc1 prerelease is considered to
	// be equivalent to the 1.8.0 core version. This enables developers
	// to assert that prerelease versions are ran with upcoming
	// core versions.
	//
	// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/303
	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.8.0")), version.Must(version.NewVersion("1.9.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MaxHigherCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.7.0-rc1")

	// The 1.7.0-rc1 prerelease should always be considered to be below the
	// 1.8.0 core version. This intentionally verifies that the logic does not
	// ignore the core version of the prerelease version when compared against
	// the core version of the check.
	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.6.0")), version.Must(version.NewVersion("1.8.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MinHigherCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.7.0-rc1")

	// The 1.7.0-rc1 prerelease should always be considered to be below the
	// 1.8.0 core version. This intentionally verifies that the logic does not
	// ignore the core version of the prerelease version when compared against
	// the core version of the check.
	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.8.0")), version.Must(version.NewVersion("1.9.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MaxHigherPrerelease(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.7.0-rc1")

	// The 1.7.0-rc1 prerelease should always be considered to be
	// below the 1.7.0-rc2 prerelease.
	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.6.0")), version.Must(version.NewVersion("1.7.0-rc2"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MinHigherPrerelease(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.7.0-rc1")

	// The 1.7.0-rc1 prerelease should always be considered to be
	// below the 1.7.0-rc2 prerelease.
	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.7.0-rc2")), version.Must(version.NewVersion("1.8.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MaxLowerCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.7.0 core version.
	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.6.0")), version.Must(version.NewVersion("1.7.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MinLowerCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.7.0 core version.
	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.7.0")), version.Must(version.NewVersion("1.9.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MaxLowerPrerelease(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.8.0-beta1 prerelease.
	testingiface.ExpectPass(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.7.0")), version.Must(version.NewVersion("1.8.0-beta1"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_SkipBetween_Prerelease_MinLowerPrerelease(t *testing.T) { //nolint:paralleltest
	t.Setenv("TF_ACC_TERRAFORM_PATH", "")
	t.Setenv("TF_ACC_TERRAFORM_VERSION", "1.8.0-rc1")

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.8.0-beta1 prerelease.
	testingiface.ExpectSkip(t, func(mockT *testingiface.MockT) {
		r.UnitTestT(mockT, r.TestCase{
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBetween(version.Must(version.NewVersion("1.8.0-beta1")), version.Must(version.NewVersion("1.9.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}
