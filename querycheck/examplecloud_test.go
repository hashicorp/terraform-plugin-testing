// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck_test

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
						"id":                  tftypes.String,
						"location":            tftypes.String,
						"name":                tftypes.String,
						"resource_group_name": tftypes.String,
						"instances":           tftypes.Number,
					},
				},
				map[string]tftypes.Value{
					"id":                  tftypes.NewValue(tftypes.String, "foo/banana"),
					"location":            tftypes.NewValue(tftypes.String, "westeurope"),
					"name":                tftypes.NewValue(tftypes.String, "banana"),
					"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
					"instances":           tftypes.NewValue(tftypes.Number, int64(5)),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"resource_group_name": tftypes.String,
						"name":                tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
					"name":                tftypes.NewValue(tftypes.String, "banana"),
				},
			)),
		},
		ReadResponse: &resource.ReadResponse{
			NewState: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":                  tftypes.String,
						"location":            tftypes.String,
						"name":                tftypes.String,
						"resource_group_name": tftypes.String,
						"instances":           tftypes.Number,
					},
				},
				map[string]tftypes.Value{
					"id":                  tftypes.NewValue(tftypes.String, "foo/banana"),
					"location":            tftypes.NewValue(tftypes.String, "westeurope"),
					"name":                tftypes.NewValue(tftypes.String, "banana"),
					"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
					"instances":           tftypes.NewValue(tftypes.Number, int64(5)),
				},
			),
			NewIdentity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"resource_group_name": tftypes.String,
						"name":                tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
					"name":                tftypes.NewValue(tftypes.String, "banana"),
				},
			)),
		},
		ImportStateResponse: &resource.ImportStateResponse{
			State: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":                  tftypes.String,
						"location":            tftypes.String,
						"name":                tftypes.String,
						"resource_group_name": tftypes.String,
						"instances":           tftypes.Number,
					},
				},
				map[string]tftypes.Value{
					"id":                  tftypes.NewValue(tftypes.String, "foo/banana"),
					"location":            tftypes.NewValue(tftypes.String, "westeurope"),
					"name":                tftypes.NewValue(tftypes.String, "banana"),
					"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
					"instances":           tftypes.NewValue(tftypes.Number, int64(5)),
				},
			),
			Identity: teststep.Pointer(tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"resource_group_name": tftypes.String,
						"name":                tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"resource_group_name": tftypes.NewValue(tftypes.String, "foo"),
					"name":                tftypes.NewValue(tftypes.String, "banana"),
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
						RequiredStringAttribute("resource_group_name"),
						OptionalNumberAttribute("instances"),
					},
				},
			},
		},
		IdentitySchemaResponse: &resource.IdentitySchemaResponse{
			Schema: &tfprotov6.ResourceIdentitySchema{
				Version: 1,
				IdentityAttributes: []*tfprotov6.ResourceIdentitySchemaAttribute{
					{
						Name:              "resource_group_name",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
					{
						Name:              "name",
						Type:              tftypes.String,
						RequiredForImport: true,
					},
				},
			},
		},
	}
}
