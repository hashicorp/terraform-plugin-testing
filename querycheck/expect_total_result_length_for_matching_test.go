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
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestResultTotalLengthForMatching(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				ListResources: map[string]testprovider.ListResource{
					"examplecloud_containerette": examplecloudListResource(),
					"examplecloud_bananette":     examplecloudListResourceBananette(),
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_containerette": examplecloudResource(),
					"examplecloud_bananette":     examplecloudResourceBananette(),
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

				list "examplecloud_bananette" "test" {
					provider = examplecloud
			
					config {
						resource_group_name = "bar"
 					}
				}

				list "examplecloud_containerette" "test2" {
					provider = examplecloud
			
					config {
						resource_group_name = "foo"
 					}
				}

				list "examplecloud_containerette" "test3" {
					provider = examplecloud
			
					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectTotalLengthForMatching(regexp.MustCompile("examplecloud_(.*)ette.test[1-9]"), 12),
				},
			},
		},
	})
}

func TestResultTotalLengthForMatching_WrongAmount(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				ListResources: map[string]testprovider.ListResource{
					"examplecloud_containerette": examplecloudListResource(),
					"examplecloud_bananette":     examplecloudListResourceBananette(),
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_containerette": examplecloudResource(),
					"examplecloud_bananette":     examplecloudResourceBananette(),
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

				list "examplecloud_bananette" "test" {
					provider = examplecloud
			
					config {
						resource_group_name = "bar"
 					}
				}

				list "examplecloud_containerette" "test2" {
					provider = examplecloud
			
					config {
						resource_group_name = "foo"
 					}
				}

				list "examplecloud_containerette" "test3" {
					provider = examplecloud
			
					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectTotalLengthForMatching(regexp.MustCompile("examplecloud_(.*)ette.test[1-9]"), 10),
				},
				ExpectError: regexp.MustCompile("expected total of found resources to be 10, got 12"),
			},
		},
	})
}

func TestResultTotalLengthForMatching_NoMatches(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				ListResources: map[string]testprovider.ListResource{
					"examplecloud_containerette": examplecloudListResource(),
					"examplecloud_bananette":     examplecloudListResourceBananette(),
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_containerette": examplecloudResource(),
					"examplecloud_bananette":     examplecloudResourceBananette(),
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

				list "examplecloud_bananette" "test" {
					provider = examplecloud
			
					config {
						resource_group_name = "bar"
 					}
				}

				list "examplecloud_containerette" "test2" {
					provider = examplecloud
			
					config {
						resource_group_name = "foo"
 					}
				}

				list "examplecloud_containerette" "test3" {
					provider = examplecloud
			
					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectTotalLengthForMatching(regexp.MustCompile("examplecloud_(.*)ette.test[4-9]"), 10),
				},
				ExpectError: regexp.MustCompile("no list resources matching the provided regex pattern .* were found in the query results"),
			},
		},
	})
}
