// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

var _ resource.Resource = Resource{}

type Resource struct {
	CreateResponse      *resource.CreateResponse
	DeleteResponse      *resource.DeleteResponse
	ImportStateResponse *resource.ImportStateResponse

	// Planning happens multiple ways during a single TestStep, so statically
	// defining only the response is very problematic.
	PlanChangeFunc func(context.Context, resource.PlanChangeRequest, *resource.PlanChangeResponse)

	ReadResponse           *resource.ReadResponse
	SchemaResponse         *resource.SchemaResponse
	UpdateResponse         *resource.UpdateResponse
	UpgradeStateResponse   *resource.UpgradeStateResponse
	ValidateConfigResponse *resource.ValidateConfigResponse
}

func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.CreateResponse != nil {
		resp.Diagnostics = r.CreateResponse.Diagnostics
		resp.NewState = r.CreateResponse.NewState
	}
}

func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.DeleteResponse != nil {
		resp.Diagnostics = r.DeleteResponse.Diagnostics
	}
}

func (r Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.ImportStateResponse != nil {
		resp.Diagnostics = r.ImportStateResponse.Diagnostics
		resp.State = r.ImportStateResponse.State
	}
}

func (r Resource) PlanChange(ctx context.Context, req resource.PlanChangeRequest, resp *resource.PlanChangeResponse) {
	if r.PlanChangeFunc != nil {
		r.PlanChangeFunc(ctx, req, resp)
	}
}

func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.ReadResponse != nil {
		resp.Diagnostics = r.ReadResponse.Diagnostics
		resp.NewState = r.ReadResponse.NewState
	}
}

func (r Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	if r.SchemaResponse != nil {
		resp.Diagnostics = r.SchemaResponse.Diagnostics
		resp.Schema = r.SchemaResponse.Schema
	}
}

func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.UpdateResponse != nil {
		resp.Diagnostics = r.UpdateResponse.Diagnostics
		resp.NewState = r.UpdateResponse.NewState
	}
}

func (r Resource) UpgradeState(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	if r.UpgradeStateResponse != nil {
		resp.Diagnostics = r.UpgradeStateResponse.Diagnostics
		resp.UpgradedState = r.UpgradeStateResponse.UpgradedState
	}
}

func (r Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r.ValidateConfigResponse != nil {
		resp.Diagnostics = r.ValidateConfigResponse.Diagnostics
	}
}
