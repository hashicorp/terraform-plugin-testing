// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck_test

import (
	"context"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

type testResource struct {}

func (t testResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) PlanChange(ctx context.Context, request resource.PlanChangeRequest, response *resource.PlanChangeResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) UpgradeState(ctx context.Context, request resource.UpgradeStateRequest, response *resource.UpgradeStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t testResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	//TODO implement me
	panic("implement me")
}


