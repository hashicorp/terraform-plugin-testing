// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				ListResources: map[string]testprovider.ListResource{
					"examplecloud_containerette": examplecloudListResource(),
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_containerette": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{ // config mode step 1 needs tf file with terraform providers block
				Config: `
				resource "examplecloud_containerette" "primary" {
					id = "westeurope/somevalue"
					location = "westeurope"
					name = "somevalue"
				}`,
			},
			{ // Query mode step 2, operates on .tfquery.hcl files (needs tf file with terraform providers block)
				// ```provider "examplecloud" {}``` has a slightly different syntax for a .tfquery.hcl file
				// provider bock simulates a real providers workflow
				// "config" in this case means configuration of the list resource/filters

				Query: true,
				Config: `
				provider "examplecloud" {} 
				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						id = "westeurope/somevalue"
 					}
				}
				list "examplecloud_containerette" "test2" {
					provider = examplecloud

					config {
						id = "foo"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"id":       knownvalue.StringExact("westeurope/somevalue1"),
						"location": knownvalue.StringExact("westeurope"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"id":       knownvalue.StringExact("westeurope/somevalue2"),
						"location": knownvalue.StringExact("westeurope2"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"id":       knownvalue.StringExact("westeurope/somevalue3"),
						"location": knownvalue.StringExact("westeurope3"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test2", map[string]knownvalue.Check{
						"id":       knownvalue.StringExact("westeurope/somevalue1"),
						"location": knownvalue.StringExact("westeurope"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test2", map[string]knownvalue.Check{
						"id":       knownvalue.StringExact("westeurope/somevalue2"),
						"location": knownvalue.StringExact("westeurope2"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test2", map[string]knownvalue.Check{
						"id":       knownvalue.StringExact("westeurope/somevalue3"),
						"location": knownvalue.StringExact("westeurope3"),
					}),
				},
			},
			{
				Query: true,
				Config: `
				provider "examplecloud" {} 
				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						id = "westeurope/somevalue"
 					}
				}
				list "examplecloud_containerette" "test2" {
					provider = examplecloud

					config {
						id = "foo"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("examplecloud_containerette.test", 3),
					querycheck.ExpectLength("examplecloud_containerette.test2", 3),
				},
			},
			{
				Query: true,
				Config: `
				provider "examplecloud" {} 
				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						id = "westeurope/somevalue"
 					}
				}
				list "examplecloud_containerette" "test2" {
					provider = examplecloud

					config {
						id = "foo"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("examplecloud_containerette.test", 2),
					querycheck.ExpectLengthAtLeast("examplecloud_containerette.test2", 1),
				},
			},
		},
	})
}
