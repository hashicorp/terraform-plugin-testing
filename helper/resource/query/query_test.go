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

	// 		--- FAIL: TestQuery (0.16s)
	//     query_test.go:20: Step 1/1 error running query: error running terraform query command: exit status 1

	//         Error: Inconsistent dependency lock file

	//         The following dependency selections recorded in the lock file are
	//         inconsistent with the current configuration:
	//           - provider registry.terraform.io/hashicorp/examplecloud: required by this configuration but no version is selected

	//         To make the initial dependency selections that will initialize the dependency
	//         lock file, run:
	//           terraform init
	t.Skip()

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
			{
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
