// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestImportBlock_WithResourceIdentity(t *testing.T) {
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
				Config: `
				resource "examplecloud_container" "test" {
					location = "westeurope"
					name     = "somevalue"
				}`,
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func TestImportBlock_WithResourceIdentity_NullAttribute(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0), // ImportBlockWithResourceIdentity requires Terraform 1.12.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResourceWithNullIdentityAttr(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "examplecloud_container" "test" {
					location = "westeurope"
					name     = "somevalue"
				}`,
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func TestImportBlock_WithResourceIdentity_WithEveryType(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0), // ImportBlockWithResourceIdentity requires Terraform 1.12.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResourceWithEveryIdentitySchemaType(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: `
				resource "examplecloud_container" "test" {
					cabinet = "A1"
					unit    = 14
					tags    = ["storage", "fast"]
					active  = true
				}`,
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func TestImportBlock_WithResourceIdentity_RequiresVersion1_12_0(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),  // ProtoV6ProviderFactories
			tfversion.SkipAbove(tfversion.Version1_11_0), // ImportBlockWithResourceIdentity requires Terraform 1.12.0 or later
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
				Config: `
				resource "examplecloud_container" "test" {
					location = "westeurope"
					name     = "somevalue"
				}`,
			},
			{
				ResourceName:    "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithResourceIdentity,
				ExpectError:     regexp.MustCompile(`Terraform 1.12.0\S* or later`),
			},
		},
	})
}
