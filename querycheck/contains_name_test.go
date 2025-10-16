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

func TestContainsResourceWithName(t *testing.T) {
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
			// We'll skip the first test step where we simulate creating the resource that will be returned when we query for it for simplicity.
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
					querycheck.ContainsResourceWithName("examplecloud_containerette.test", "banane"),
					querycheck.ContainsResourceWithName("examplecloud_containerette.test", "ananas"),
					querycheck.ContainsResourceWithName("examplecloud_containerette.test", "kiwi"),
					querycheck.ContainsResourceWithName("examplecloud_containerette.test2", "papaya"),
					querycheck.ContainsResourceWithName("examplecloud_containerette.test2", "birne"),
					querycheck.ContainsResourceWithName("examplecloud_containerette.test2", "kirsche"),
				},
			},
		},
	})
}

// Let's add a test case that checks the failure scenario when a resource of a given name is not found.
func TestContainsResourceWithName_NotFound(t *testing.T) {
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
					querycheck.ContainsResourceWithName("examplecloud_containerette.test", "pflaume"),
				},
				ExpectError: regexp.MustCompile("expected to find resource with display name \"pflaume\" in results but resource was not found"),
			},
		},
	})
}
