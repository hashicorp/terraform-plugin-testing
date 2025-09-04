// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
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
					"examplecloud_containerette": {
						SchemaResponse: &list.SchemaResponse{
							Schema: &tfprotov6.Schema{
								Block: &tfprotov6.SchemaBlock{
									Attributes: []*tfprotov6.SchemaAttribute{
										ComputedStringAttribute("id"),
									},
								},
							},
						},
						ListResultsStream: &list.ListResultsStream{
							Results: func(push func(list.ListResult) bool) {
							},
						},
					},
				},
				Resources: map[string]testprovider.Resource{
					"examplecloud_containerette": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{ // config mode step 1, creates something we can list later, need tf file with terraform providers block
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
				// run query in terraform itself maybe by moving this provider to the corner provider to check if it works

				Query: true,
				Config: `
				provider "examplecloud" {} 
				list "examplecloud_containerette" "test" {
					provider = examplecloud

					config {
						id = "bat"
					}
				}`,
			},
		},
	})
}
