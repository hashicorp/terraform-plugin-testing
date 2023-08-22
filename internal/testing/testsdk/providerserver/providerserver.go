package providerserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

var _ tfprotov6.ProviderServer = ProviderServer{}

// NewProviderServer returns a lightweight protocol version 6 provider server
// for consumption with ProtoV6ProviderFactories.
func NewProviderServer(p provider.Provider) func() (tfprotov6.ProviderServer, error) {
	return NewProviderServerWithError(p, nil)
}

// NewProviderServerWithError returns a lightweight protocol version 6 provider
// server and an associated error for consumption with ProtoV6ProviderFactories.
func NewProviderServerWithError(p provider.Provider, err error) func() (tfprotov6.ProviderServer, error) {
	providerServer := ProviderServer{
		Provider: p,
	}

	return func() (tfprotov6.ProviderServer, error) {
		return providerServer, err
	}
}

// ProviderServer is a lightweight protocol version 6 provider server which
// is assumed to be well-behaved, e.g. does not return gRPC errors.
//
// This implementation intends to reduce the heaviest parts of
// terraform-plugin-go based provider development:
//
//   - Converting *tfprotov6.DynamicValue to tftypes.Value using schema
//   - Splitting ApplyResourceChange into Create/Update/Delete calls
//   - Set PlanResourceChange null config values of Computed attributes to unknown
//   - Roundtrip UpgradeResourceState with equal schema version
//
// By default, the following data is copied automatically:
//
//   - ApplyResourceChange (create): req.Config -> resp.NewState
//   - ApplyResourceChange (delete): req.PlannedState -> resp.NewState
//   - ApplyResourceChange (update): req.PlannedState -> resp.NewState
//   - PlanResourceChange: req.ProposedNewState -> resp.PlannedState
//   - ReadDataSource: req.Config -> resp.State
//   - ReadResource: req.CurrentState -> resp.NewState
type ProviderServer struct {
	Provider provider.Provider
}

func (s ProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	resp := &tfprotov6.ApplyResourceChangeResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	plannedState, diag := DynamicValueToValue(schemaResp.Schema, req.PlannedState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	priorState, diag := DynamicValueToValue(schemaResp.Schema, req.PriorState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	if priorState.IsNull() {
		createReq := resource.CreateRequest{
			Config: config,
		}
		createResp := &resource.CreateResponse{
			NewState: config.Copy(),
		}

		r.Create(ctx, createReq, createResp)

		resp.Diagnostics = createResp.Diagnostics

		if len(resp.Diagnostics) > 0 {
			return resp, nil
		}

		newState, diag := ValuetoDynamicValue(schemaResp.Schema, createResp.NewState)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.NewState = newState
	} else if plannedState.IsNull() {
		deleteReq := resource.DeleteRequest{
			PriorState: priorState,
		}
		deleteResp := &resource.DeleteResponse{}

		r.Delete(ctx, deleteReq, deleteResp)

		resp.Diagnostics = deleteResp.Diagnostics

		if len(resp.Diagnostics) > 0 {
			return resp, nil
		}

		resp.NewState = req.PlannedState
	} else {
		updateReq := resource.UpdateRequest{
			Config:       config,
			PlannedState: plannedState,
			PriorState:   priorState,
		}
		updateResp := &resource.UpdateResponse{
			NewState: plannedState.Copy(),
		}

		r.Update(ctx, updateReq, updateResp)

		resp.Diagnostics = updateResp.Diagnostics

		if len(resp.Diagnostics) > 0 {
			return resp, nil
		}

		newState, diag := ValuetoDynamicValue(schemaResp.Schema, updateResp.NewState)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.NewState = newState
	}

	return resp, nil
}

func (s ProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	resp := &tfprotov6.ConfigureProviderResponse{}

	schemaReq := provider.SchemaRequest{}
	schemaResp := &provider.SchemaResponse{}

	s.Provider.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	configureReq := provider.ConfigureRequest{
		Config: config,
	}
	configureResp := &provider.ConfigureResponse{}

	s.Provider.Configure(ctx, configureReq, configureResp)

	resp.Diagnostics = configureResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	providerReq := provider.SchemaRequest{}
	providerResp := &provider.SchemaResponse{}

	s.Provider.Schema(ctx, providerReq, providerResp)

	resp := &tfprotov6.GetProviderSchemaResponse{
		DataSourceSchemas: map[string]*tfprotov6.Schema{},
		Diagnostics:       providerResp.Diagnostics,
		Provider:          providerResp.Schema,
		ResourceSchemas:   map[string]*tfprotov6.Schema{},
		ServerCapabilities: &tfprotov6.ServerCapabilities{
			PlanDestroy: true,
		},
	}

	for typeName, d := range s.Provider.DataSourcesMap() {
		schemaReq := datasource.SchemaRequest{}
		schemaResp := &datasource.SchemaResponse{}

		d.Schema(ctx, schemaReq, schemaResp)

		resp.Diagnostics = append(resp.Diagnostics, schemaResp.Diagnostics...)

		resp.DataSourceSchemas[typeName] = schemaResp.Schema
	}

	for typeName, r := range s.Provider.ResourcesMap() {
		schemaReq := resource.SchemaRequest{}
		schemaResp := &resource.SchemaResponse{}

		r.Schema(ctx, schemaReq, schemaResp)

		resp.Diagnostics = append(resp.Diagnostics, schemaResp.Diagnostics...)

		resp.ResourceSchemas[typeName] = schemaResp.Schema
	}

	return resp, nil
}

func (s ProviderServer) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	resp := &tfprotov6.ImportResourceStateResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	importReq := resource.ImportStateRequest{
		ID: req.ID,
	}
	importResp := &resource.ImportStateResponse{}

	r.ImportState(ctx, importReq, importResp)

	resp.Diagnostics = importResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	if importResp.State.IsNull() {
		return resp, nil
	}

	state, diag := ValuetoDynamicValue(schemaResp.Schema, importResp.State)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.ImportedResources = []*tfprotov6.ImportedResource{
		{
			State:    state,
			TypeName: req.TypeName,
		},
	}

	return resp, nil
}

func (s ProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	resp := &tfprotov6.PlanResourceChangeResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	priorState, diag := DynamicValueToValue(schemaResp.Schema, req.PriorState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	proposedNewState, diag := DynamicValueToValue(schemaResp.Schema, req.ProposedNewState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	if !proposedNewState.IsNull() && !proposedNewState.Equal(priorState) {
		modifiedProposedNewState, err := tftypes.Transform(proposedNewState, func(path *tftypes.AttributePath, val tftypes.Value) (tftypes.Value, error) {
			// we are only modifying attributes, not the entire resource
			if len(path.Steps()) < 1 {
				return val, nil
			}

			configValIface, _, err := tftypes.WalkAttributePath(config, path)

			if err != nil && err != tftypes.ErrInvalidStep {
				return val, fmt.Errorf("error walking attribute/block path during unknown marking: %w", err)
			}

			configVal, ok := configValIface.(tftypes.Value)

			if !ok {
				return val, fmt.Errorf("unexpected type during unknown marking: %T", configValIface)
			}

			if !configVal.IsNull() {
				return val, nil
			}

			attribute := SchemaAttributeAtPath(schemaResp.Schema, path)

			if attribute == nil {
				return val, nil
			}

			if !attribute.Computed {
				return val, nil
			}

			return tftypes.NewValue(val.Type(), tftypes.UnknownValue), nil
		})

		if err != nil {
			diag := &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error Modifying ProposedNewState",
				Detail:   err.Error(),
			}

			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
		}

		proposedNewState = modifiedProposedNewState
	}

	planReq := resource.PlanChangeRequest{
		Config:           config,
		PriorState:       priorState,
		ProposedNewState: proposedNewState,
	}
	planResp := &resource.PlanChangeResponse{
		PlannedState: proposedNewState.Copy(),
	}

	r.PlanChange(ctx, planReq, planResp)

	resp.Diagnostics = planResp.Diagnostics
	resp.RequiresReplace = planResp.RequiresReplace

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	plannedState, diag := ValuetoDynamicValue(schemaResp.Schema, planResp.PlannedState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.PlannedState = plannedState

	return resp, nil
}

func (s ProviderServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	resp := &tfprotov6.ReadDataSourceResponse{}

	d, diag := ProviderDataSource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := &datasource.SchemaResponse{}

	d.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	readReq := datasource.ReadRequest{
		Config: config,
	}
	readResp := &datasource.ReadResponse{
		State: config.Copy(),
	}

	d.Read(ctx, readReq, readResp)

	resp.Diagnostics = readResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	state, diag := ValuetoDynamicValue(schemaResp.Schema, readResp.State)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.State = state

	return resp, nil
}

func (s ProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	resp := &tfprotov6.ReadResourceResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	currentState, diag := DynamicValueToValue(schemaResp.Schema, req.CurrentState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	readReq := resource.ReadRequest{
		CurrentState: currentState,
	}
	readResp := &resource.ReadResponse{
		NewState: currentState.Copy(),
	}

	r.Read(ctx, readReq, readResp)

	resp.Diagnostics = readResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	newState, diag := ValuetoDynamicValue(schemaResp.Schema, readResp.NewState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.NewState = newState

	return resp, nil
}

func (s ProviderServer) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	providerReq := provider.StopRequest{}
	providerResp := &provider.StopResponse{}

	s.Provider.Stop(ctx, providerReq, providerResp)

	resp := &tfprotov6.StopProviderResponse{}

	if providerResp.Error != nil {
		resp.Error = providerResp.Error.Error()
	}

	return resp, nil
}

func (s ProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	// Define options to be used when unmarshalling raw state.
	// IgnoreUndefinedAttributes will silently skip over fields in the JSON
	// that do not have a matching entry in the schema.
	unmarshalOpts := tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{
			IgnoreUndefinedAttributes: true,
		},
	}

	// Terraform CLI can call UpgradeResourceState even if the stored state
	// version matches the current schema. Presumably this is to account for
	// the previous terraform-plugin-sdk implementation, which handled some
	// state fixups on behalf of Terraform CLI. This will attempt to roundtrip
	// the prior RawState to a state matching the current schema.
	if req.Version == schemaResp.Schema.Version {
		rawStateValue, err := req.RawState.UnmarshalWithOpts(schemaResp.Schema.ValueType(), unmarshalOpts)

		if err != nil {
			diag := &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unable to Read Previously Saved State for UpgradeResourceState",
				Detail:   "There was an error reading the saved resource state using the current resource schema: " + err.Error(),
			}

			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
		}

		upgradedState, diag := ValuetoDynamicValue(schemaResp.Schema, rawStateValue)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.UpgradedState = upgradedState

		return resp, nil
	}

	upgradeReq := resource.UpgradeStateRequest{}
	upgradeResp := &resource.UpgradeStateResponse{}

	r.UpgradeState(ctx, upgradeReq, upgradeResp)

	resp.Diagnostics = upgradeResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	upgradedState, diag := ValuetoDynamicValue(schemaResp.Schema, upgradeResp.UpgradedState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.UpgradedState = upgradedState

	return resp, nil
}

func (s ProviderServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateDataResourceConfigResponse{}

	d, diag := ProviderDataSource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := &datasource.SchemaResponse{}

	d.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	validateReq := datasource.ValidateConfigRequest{
		Config: config,
	}
	validateResp := &datasource.ValidateConfigResponse{}

	d.ValidateConfig(ctx, validateReq, validateResp)

	resp.Diagnostics = validateResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	providerReq := provider.ValidateConfigRequest{}
	providerResp := &provider.ValidateConfigResponse{}

	s.Provider.ValidateConfig(ctx, providerReq, providerResp)

	resp := &tfprotov6.ValidateProviderConfigResponse{
		Diagnostics:    providerResp.Diagnostics,
		PreparedConfig: req.Config,
	}

	return resp, nil
}

func (s ProviderServer) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateResourceConfigResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	validateReq := resource.ValidateConfigRequest{
		Config: config,
	}
	validateResp := &resource.ValidateConfigResponse{}

	r.ValidateConfig(ctx, validateReq, validateResp)

	resp.Diagnostics = validateResp.Diagnostics

	return resp, nil
}
