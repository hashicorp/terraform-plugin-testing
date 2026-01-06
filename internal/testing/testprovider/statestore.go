// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

var _ statestore.StateStore = &StateStore{}

type StateStore struct {
	configuredChunkSize int64

	SchemaResponse          *statestore.SchemaResponse
	ConfigureResponse       *statestore.ConfigureResponse
	ValidateConfigResponse  *statestore.ValidateConfigResponse
	GetStatesResponse       *statestore.GetStatesResponse
	DeleteStateResponse     *statestore.DeleteStateResponse
	LockStateResponse       *statestore.LockStateResponse
	UnlockStateResponse     *statestore.UnlockStateResponse
	ReadStateBytesResponse  *statestore.ReadStateBytesResponse
	WriteStateBytesResponse *statestore.WriteStateBytesResponse
}

func (s *StateStore) Schema(ctx context.Context, req statestore.SchemaRequest, resp *statestore.SchemaResponse) {
	if s.SchemaResponse != nil {
		resp.Diagnostics = s.SchemaResponse.Diagnostics
		resp.Schema = s.SchemaResponse.Schema
	}
}

func (s *StateStore) Configure(ctx context.Context, req statestore.ConfigureRequest, resp *statestore.ConfigureResponse) {
	if s.ConfigureResponse != nil {
		resp.Diagnostics = s.ConfigureResponse.Diagnostics
		resp.ServerCapabilities = s.ConfigureResponse.ServerCapabilities
	}

	// Store configured chunk size
	s.configuredChunkSize = resp.ServerCapabilities.ChunkSize
}

func (s *StateStore) ConfiguredChunkSize() int64 {
	return s.configuredChunkSize
}

func (s *StateStore) ValidateConfig(ctx context.Context, req statestore.ValidateConfigRequest, resp *statestore.ValidateConfigResponse) {
	if s.ValidateConfigResponse != nil {
		resp.Diagnostics = s.ValidateConfigResponse.Diagnostics
	}
}

// TODO:PSS: Probably need to adjust some of these to be callback functions, rather then field responses.
func (s *StateStore) GetStates(ctx context.Context, req statestore.GetStatesRequest, resp *statestore.GetStatesResponse) {
	if s.GetStatesResponse != nil {
		resp.Diagnostics = s.GetStatesResponse.Diagnostics
		resp.StateIDs = s.GetStatesResponse.StateIDs
	}
}

func (s *StateStore) DeleteState(ctx context.Context, req statestore.DeleteStateRequest, resp *statestore.DeleteStateResponse) {
	if s.DeleteStateResponse != nil {
		resp.Diagnostics = s.DeleteStateResponse.Diagnostics
	}
}

func (s *StateStore) LockState(ctx context.Context, req statestore.LockStateRequest, resp *statestore.LockStateResponse) {
	if s.LockStateResponse != nil {
		resp.LockID = s.LockStateResponse.LockID
		resp.Diagnostics = s.LockStateResponse.Diagnostics
	}
}

func (s *StateStore) UnlockState(ctx context.Context, req statestore.UnlockStateRequest, resp *statestore.UnlockStateResponse) {
	if s.UnlockStateResponse != nil {
		resp.Diagnostics = s.UnlockStateResponse.Diagnostics
	}
}

func (s *StateStore) ReadStateBytes(ctx context.Context, req statestore.ReadStateBytesRequest, resp *statestore.ReadStateBytesResponse) {
	if s.ReadStateBytesResponse != nil {
		resp.Diagnostics = s.ReadStateBytesResponse.Diagnostics
		resp.StateBytes = s.ReadStateBytesResponse.StateBytes
	}
}

func (s *StateStore) WriteStateBytes(ctx context.Context, req statestore.WriteStateBytesRequest, resp *statestore.WriteStateBytesResponse) {
	if s.WriteStateBytesResponse != nil {
		resp.Diagnostics = s.WriteStateBytesResponse.Diagnostics
	}
}
