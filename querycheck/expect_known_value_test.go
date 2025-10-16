// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestExpectKnownValue(t *testing.T) {
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
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectKnownValue(
						"examplecloud_containerette.test",
						"banana",
						tfjsonpath.New("instances"),
						knownvalue.NumberExact(big.NewFloat(5)),
					),
				},
				ExpectError: regexp.MustCompile("examplecloud_containerette.test - the resource banana was not found"),
			},
		},
	})
}

// Let's add a test case that checks the failure scenario when the value is incorrect.
func TestExpectKnownValue_ValueIncorrect(t *testing.T) {
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
				`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectKnownValue(
						"examplecloud_containerette.test",
						"banana",
						tfjsonpath.New("instances"),
						knownvalue.NumberExact(big.NewFloat(4)),
					),
				},
				ExpectError: regexp.MustCompile("examplecloud_containerette.test - the resource banana was not found"),
			},
		},
	})
}
