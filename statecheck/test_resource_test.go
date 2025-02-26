// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

type aTestResource struct{}

func (t aTestResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) PlanChange(ctx context.Context, request resource.PlanChangeRequest, response *resource.PlanChangeResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	schema := tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:     "bool_attribute",
					Type:     tftypes.Bool,
					Computed: true,
				},
			},
			Description:     "",
			DescriptionKind: 0,
		},
	}
	response.Schema = &schema
}

func (t aTestResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) UpgradeState(ctx context.Context, request resource.UpgradeStateRequest, response *resource.UpgradeStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t aTestResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
}
