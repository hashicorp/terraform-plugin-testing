// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_TestStep_ImportCommand_AsFirstStep(t *testing.T) {
	t.Parallel()

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0), // ProtoV6ProviderFactories needs Terraform 1.0.0 or later
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
				ImportStateKind: r.ImportCommandWithId,
				// ImportStateVerify: true,
				Config: `resource "examplecloud_container" "test" {
					name = "somevalue"
					location = "westeurope"
				}`,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state; got %d", len(states))
					}
					if states[0].ID != "westeurope/somevalue" {
						return fmt.Errorf("unexpected ID: %s", states[0].ID)
					}
					if states[0].Attributes["name"] != "somevalue" {
						return fmt.Errorf("unexpected name: %s", states[0].Attributes["name"])
					}
					if states[0].Attributes["location"] != "westeurope" {
						return fmt.Errorf("unexpected location: %s", states[0].Attributes["location"])
					}
					return nil
				},
			},
			{
				RefreshState: true,
				Check:        r.TestCheckResourceAttr("examplecloud_container.test", "name", "somevalue"),
			},
		},
	})
}
