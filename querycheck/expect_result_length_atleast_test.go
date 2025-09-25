package querycheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestResultLengthAtLeast(t *testing.T) {
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

// Let's add a test case that checks the failure scenario when there are too few results.
func TestResultLengthAtLeast_TooFewResults(t *testing.T) {
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
					querycheck.ExpectLengthAtLeast("examplecloud_containerette.test", 1),
					querycheck.ExpectLengthAtLeast("examplecloud_containerette.test2", 0),
				},
			},
		},
	})
}
