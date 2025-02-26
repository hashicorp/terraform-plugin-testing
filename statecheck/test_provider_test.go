// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

var anTestProvider = testprovider.Provider{
	Resources: map[string]testprovider.Resource{
		"test_resource": {
			SchemaResponse: &resource.SchemaResponse{
				Schema: &tfprotov6.Schema{
					Block: &tfprotov6.SchemaBlock{
						Attributes: []*tfprotov6.SchemaAttribute{
							{
								Name:     "bool_attribute",
								Type:     tftypes.Bool,
								Optional: true,
							},
							{
								Name:     "float_attribute",
								Type:     tftypes.Number,
								Optional: true,
							},
							{
								Name:     "int_attribute",
								Type:     tftypes.Number,
								Optional: true,
							},
							{
								Name:     "string_attribute",
								Type:     tftypes.String,
								Optional: true,
							},
							{
								Name:     "string_computed_attribute",
								Type:     tftypes.String,
								Computed: true,
							},
						},
					},
				},
			},
		},
	},
}
