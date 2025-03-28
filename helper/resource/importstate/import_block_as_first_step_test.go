// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_ImportBlock_AsFirstStep(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_container": examplecloudResource(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				ResourceName:    "examplecloud_container.test",
				ImportStateId:   "examplecloud_container.test",
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithID,
				// ImportStateVerify: true,
				Config: `resource "examplecloud_container" "test" {
					name = "somevalue"
					location = "westeurope"
				}`,
				ImportStatePersist: true,
				ImportPlanChecks: r.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("examplecloud_container.test", plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue("examplecloud_container.test", tfjsonpath.New("id"), knownvalue.StringExact("westeurope/somevalue")),
						plancheck.ExpectKnownValue("examplecloud_container.test", tfjsonpath.New("name"), knownvalue.StringExact("somevalue")),
						plancheck.ExpectKnownValue("examplecloud_container.test", tfjsonpath.New("location"), knownvalue.StringExact("westeurope")),
					},
				},
			},
		},
	})
}
