// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestImportBlockForResourceWithADependency(t *testing.T) {
	t.Parallel()

	config := `
resource "examplecloud_zone" "zone" {
	name = "example.net"
}

resource "examplecloud_zone_record" "record" {
	zone_id = examplecloud_zone.zone.id
	name    = "www"
}
`
	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0), // ImportBlockWithID requires Terraform 1.5.0 or later
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"examplecloud": providerserver.NewProviderServer(testprovider.Provider{
				Resources: map[string]testprovider.Resource{
					"examplecloud_zone":        examplecloudZone(),
					"examplecloud_zone_record": examplecloudZoneRecord(),
				},
			}),
		},
		Steps: []r.TestStep{
			{
				Config: config,
			},
			{
				ImportState:     true,
				ImportStateKind: r.ImportBlockWithID,
				ResourceName:    "examplecloud_zone_record.record",
			},
		},
	})
}
