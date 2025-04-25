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
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func Test_SkipBelow_Lower(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.0.7",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": func() (tfprotov6.ProviderServer, error) { //nolint:unparam // required signature
				return nil, nil
			},
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.1.0"))),
		},
		Steps: []r.TestStep{
			{
				//nullable argument only available in TF v1.1.0+
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}

func Test_SkipBelow_Equal(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.1.0",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.1.0"))),
		},
		Steps: []r.TestStep{
			{
				//nullable argument only available in TF v1.1.0+
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}

func Test_SkipBelow_Higher(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.1.1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.1.0"))),
		},
		Steps: []r.TestStep{
			{
				//nullable argument only available in TF v1.1.0+
				Config: `variable "a" {
  					nullable = true
					default  = "hello"
				}`,
			},
		},
	})
}

func Test_SkipBelow_Prerelease_EqualCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// Pragmatic compromise that 1.8.0-rc1 prerelease is considered to
	// be equivalent to the 1.8.0 core release. This enables developers
	// to assert that prerelease versions are compatible with upcoming
	// core versions.
	//
	// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/303
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipBelow_Prerelease_EqualPrerelease(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0-rc1"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipBelow_Prerelease_HigherCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.7.0-rc1 prerelease should always be considered to be below the
	// 1.8.0 core version. This intentionally verifies that the logic does not
	// ignore the core version of the prerelease version when compared against
	// the core version of the check.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.7.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipBelow_Prerelease_HigherPrerelease(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.7.0-rc1 prerelease should always be considered to be
	// below the 1.7.0-rc2 prerelease.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.7.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.7.0-rc2"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipBelow_Prerelease_LowerCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.7.0 core version.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.7.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_SkipBelow_Prerelease_LowerPrerelease(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.8.0-beta1 prerelease.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}
