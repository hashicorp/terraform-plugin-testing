// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package importstate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_ImportBlock_InConfigDirectory(t *testing.T) {
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
				ConfigDirectory: func(config.TestStepConfigRequest) string {
					return `testdata/1`
				},
			},
			{
				ResourceName:      "examplecloud_container.test",
				ImportState:       true,
				ImportStateKind:   r.ImportBlockWithID,
				ImportStateVerify: true,
				ConfigDirectory: func(config.TestStepConfigRequest) string {
					return `testdata/2`
				},
			},
		},
	})
}
