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
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.4.3",
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
	t.Parallel()

	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			TFExactVersion: "1.1.0",
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

func Test_RequireNot_Prerelease_EqualCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// Pragmatic compromise that 1.8.0-rc1 prerelease is considered to
	// be equivalent to the 1.8.0 core version. This enables developers
	// to assert that prerelease versions are not ran with upcoming
	// core versions.
	//
	// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/303
	plugintest.TestExpectTFatal(t, func() {
		r.UnitTest(&testinginterface.RuntimeT{}, r.TestCase{
			TFExactVersion: "1.8.0-rc1",
			ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"test": providerserver.NewProviderServer(testprovider.Provider{}),
			},
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.RequireNot(version.Must(version.NewVersion("1.8.0"))),
			},
			Steps: []r.TestStep{
				{
					Config: `//non-empty config`,
				},
			},
		})
	})
}

func Test_RequireNot_Prerelease_HigherCoreVersion(t *testing.T) { //nolint:paralleltest
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
			tfversion.RequireNot(version.Must(version.NewVersion("1.8.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_RequireNot_Prerelease_HigherPrerelease(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.7.0-rc1 prerelease should always be considered to be
	// below the 1.7.0-rc2 prerelease.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.7.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireNot(version.Must(version.NewVersion("1.7.0-rc2"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_RequireNot_Prerelease_LowerCoreVersion(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.7.0 core version.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireNot(version.Must(version.NewVersion("1.7.0"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}

func Test_RequireNot_Prerelease_LowerPrerelease(t *testing.T) { //nolint:paralleltest
	t.Parallel()

	// The 1.8.0-rc1 prerelease should always be considered to be
	// above the 1.8.0-beta1 prerelease.
	r.UnitTest(t, r.TestCase{
		TFExactVersion: "1.8.0-rc1",
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"test": providerserver.NewProviderServer(testprovider.Provider{}),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireNot(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		Steps: []r.TestStep{
			{
				Config: `//non-empty config`,
			},
		},
	})
}
