// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestImportBlock_InConfigFile(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/1/examplecloud_container.tf`),
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithID,
				ConfigFile:      config.StaticFile(`testdata/writeable-config-file/examplecloud_container.tf`),
			},
		},
	})
}

func TestImportBlock_WithResourceIdentity_InConfigFile(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0), // ImportBlockWithResourceIdentity requires Terraform 1.12.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/1/examplecloud_container.tf`),
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithResourceIdentity,
				ConfigFile:      config.StaticFile(`testdata/writeable-config-file/examplecloud_container.tf`),
			},
		},
	})
}

func TestImportBlock_InConfigFile_ConfigExactTrue(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				ConfigFile: config.StaticFile(`testdata/1/examplecloud_container.tf`),
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithResourceIdentity,
				ConfigFile:      config.StaticFile(`testdata/examplecloud_container_import_with_identity.tf`),
				ConfigExact:     true,
			},
		},
	})
}
