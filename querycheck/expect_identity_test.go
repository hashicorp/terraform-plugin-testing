// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectIdentity(t *testing.T) {
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
				// this step should provision all the resources that the query is support to list
				// for simplicity we're only "provisioning" one here
				Config: `
				resource "examplecloud_containerette" "primary" {
					name                = "banana"
					resource_group_name = "foo"
					location  			= "westeurope"
			
					instances = 5
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
						resource_group_name = "foo"
 					}
				}

				list "examplecloud_containerette" "test2" {
					provider = examplecloud

					config {
						resource_group_name = "bar"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("banane"),
						"resource_group_name": knownvalue.StringExact("foo"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("ananas"),
						"resource_group_name": knownvalue.StringExact("foo"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("kiwi"),
						"resource_group_name": knownvalue.StringExact("foo"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test2", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("papaya"),
						"resource_group_name": knownvalue.StringExact("bar"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test2", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("birne"),
						"resource_group_name": knownvalue.StringExact("bar"),
					}),
					querycheck.ExpectIdentity("examplecloud_containerette.test2", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("kirsche"),
						"resource_group_name": knownvalue.StringExact("bar"),
					}),
				},
			},
		},
	})
}

func TestExpectIdentity_NotFound(t *testing.T) {
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
				// this step should provision all the resources that the query is support to list
				// for simplicity we're only "provisioning" one here
				Config: `
				resource "examplecloud_containerette" "primary" {
					name                = "banana"
					resource_group_name = "foo"
					location  			= "westeurope"
			
					instances = 5
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
						resource_group_name = "foo"
 					}
				}

				list "examplecloud_containerette" "test2" {
					provider = examplecloud

					config {
						resource_group_name = "bar"
					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("examplecloud_containerette.test", map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("owo"),
						"resource_group_name": knownvalue.StringExact("uwu"),
					}),
				},
				ExpectError: regexp.MustCompile("an identity with the following attributes was not found\nattribute \"name\": owo\nattribute \"resource_group_name\": uwu\naddress: examplecloud_containerette.test"),
			},
		},
	})
}
