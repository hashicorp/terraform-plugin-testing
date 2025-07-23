// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query_test

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

func examplecloudResource() testprovider.Resource {
	return testprovider.Resource{
		CreateResponse: &resource.CreateResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":       tftypes.String,
						"location": tftypes.String,
						"name":     tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
					"location": tftypes.NewValue(tftypes.String, "westeurope"),
					"name":     tftypes.NewValue(tftypes.String, "somevalue"),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
				},
			)),
		},
		ReadResponse: &resource.ReadResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":       tftypes.String,
						"location": tftypes.String,
						"name":     tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
					"location": tftypes.NewValue(tftypes.String, "westeurope"),
					"name":     tftypes.NewValue(tftypes.String, "somevalue"),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
				},
			)),
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":       tftypes.String,
						"location": tftypes.String,
						"name":     tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id":       tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
					"location": tftypes.NewValue(tftypes.String, "westeurope"),
					"name":     tftypes.NewValue(tftypes.String, "somevalue"),
				},
			),
			Identity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, "westeurope/somevalue"),
				},
			)),
		},
		SchemaResponse: &resource.SchemaResponse{
			Schema: &tfprotov6.Schema{
				Block: &tfprotov6.SchemaBlock{
					Attributes: []*tfprotov6.SchemaAttribute{
						ComputedStringAttribute("id"),
						RequiredStringAttribute("location"),
						RequiredStringAttribute("name"),
					},
				},
			},
		},
		IdentitySchemaResponse: &resource.IdentitySchemaResponse{
			Schema: &tfprotov6.ResourceIdentitySchema{
				Version: 1,
				IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
					{
						Name:              "id",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
				},
			},
		},
	}
}
