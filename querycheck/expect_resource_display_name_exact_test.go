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
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectResourceDisplayNameExact(t *testing.T) {
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
			{
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
					querycheck.ExpectResourceDisplayNameExact("examplecloud_containerette.test", queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("ananas"),
						"resource_group_name": knownvalue.StringExact("foo"),
					}), "ananas"),
					querycheck.ExpectResourceDisplayNameExact("examplecloud_containerette.test", queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("banane"),
						"resource_group_name": knownvalue.StringExact("foo"),
					}), "banane"),
				},
			},
		},
	})
}

func TestExpectResourceDisplayNameExact_TooManyResults(t *testing.T) {
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
			{
				Query: true,
				Config: `
				provider "examplecloud" {} 

				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceDisplayNameExact("examplecloud_containerette.test", nil, "ananas"),
				},
				ExpectError: regexp.MustCompile("examplecloud_containerette.test - more than 1 query result found after filtering"),
			},
		},
	})
}

func TestExpectResourceDisplayNameExact_NoResults(t *testing.T) {
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
			{
				Query: true,
				Config: `
				provider "examplecloud" {} 

				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceDisplayNameExact("examplecloud_containerette.test", queryfilter.ByResourceIdentity(map[string]knownvalue.Check{}),
						"ananas"),
				},
				ExpectError: regexp.MustCompile("examplecloud_containerette.test - no query results found after filtering"),
			},
		},
	})
}

func TestExpectResourceDisplayNameExact_InvalidDisplayName(t *testing.T) {
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
			{
				Query: true,
				Config: `
				provider "examplecloud" {} 

				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceDisplayNameExact("examplecloud_containerette.test", queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
						"name":                knownvalue.StringExact("ananas"),
						"resource_group_name": knownvalue.StringExact("foo"),
					}), "invalid"),
				},
				ExpectError: regexp.MustCompile("expected to find resource with display name \"invalid\" in results but resource was not found"),
			},
		},
	})
}
