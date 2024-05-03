// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Resource interface {
	Create(context.Context, CreateRequest, *CreateResponse)
	Delete(context.Context, DeleteRequest, *DeleteResponse)
	ImportState(context.Context, ImportStateRequest, *ImportStateResponse)
	PlanChange(context.Context, PlanChangeRequest, *PlanChangeResponse)
	Read(context.Context, ReadRequest, *ReadResponse)
	Schema(context.Context, SchemaRequest, *SchemaResponse)
	Update(context.Context, UpdateRequest, *UpdateResponse)
	UpgradeState(context.Context, UpgradeStateRequest, *UpgradeStateResponse)
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}

type CreateRequest struct {
	Config tftypes.Value
}

type CreateResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	NewState    tftypes.Value
}

type DeleteRequest struct {
	PriorState tftypes.Value
}

type DeleteResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type ImportStateRequest struct {
	ID string
}

type ImportStateResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	State       tftypes.Value
}

type PlanChangeRequest struct {
	Config           tftypes.Value
	PriorState       tftypes.Value
	ProposedNewState tftypes.Value
}

type PlanChangeResponse struct {
	Deferred        *tfprotov6.Deferred
	Diagnostics     []*tfprotov6.Diagnostic
	PlannedState    tftypes.Value
	RequiresReplace []*tftypes.AttributePath
}

type ReadRequest struct {
	CurrentState tftypes.Value
}

type ReadResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	NewState    tftypes.Value
}

type SchemaRequest struct{}

type SchemaResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	Schema      *tfprotov6.Schema
}

type UpdateRequest struct {
	Config       tftypes.Value
	PlannedState tftypes.Value
	PriorState   tftypes.Value
}

type UpdateResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	NewState    tftypes.Value
}

type UpgradeStateRequest struct {
	RawState *tfprotov6.RawState
	Version  int64
}

type UpgradeStateResponse struct {
	Diagnostics   []*tfprotov6.Diagnostic
	UpgradedState tftypes.Value
}

type ValidateConfigRequest struct {
	Config tftypes.Value
}

type ValidateConfigResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}
