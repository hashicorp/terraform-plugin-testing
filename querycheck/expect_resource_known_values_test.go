// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

import (
	"math/big"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectResourceKnownValues(t *testing.T) {
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
					provider         = examplecloud
					include_resource = true

					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceKnownValues(
						"examplecloud_containerette.test", queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"name":                knownvalue.StringExact("banane"),
							"resource_group_name": knownvalue.StringExact("foo"),
						}), []querycheck.KnownValueCheck{
							{
								tfjsonpath.New("instances"),
								knownvalue.NumberExact(big.NewFloat(5)),
							},
							{
								tfjsonpath.New("location"),
								knownvalue.StringExact("westeurope"),
							},
						},
					),
				},
			},
		},
	})
}

func TestExpectResourceKnownValues_ValueIncorrect(t *testing.T) {
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
					provider         = examplecloud
					include_resource = true

					config {
						resource_group_name = "foo"
 					}
				}
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectResourceKnownValues(
						"examplecloud_containerette.test", queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
							"name":                knownvalue.StringExact("banane"),
							"resource_group_name": knownvalue.StringExact("foo"),
						}), []querycheck.KnownValueCheck{
							{
								tfjsonpath.New("location"),
								knownvalue.StringExact("westeurope"),
							},
							{
								tfjsonpath.New("instances"),
								knownvalue.NumberExact(big.NewFloat(4)),
							},
						},
					),
				},
				ExpectError: regexp.MustCompile("the following errors were found while checking values: error checking value for attribute at path: instances for resource with identity .*, err: expected value 4 for NumberExact check, got: 5;"),
			},
		},
	})
}
