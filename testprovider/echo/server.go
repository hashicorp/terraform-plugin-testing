// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package echo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func NewProviderServer() func() (tfprotov6.ProviderServer, error) {
	return func() (tfprotov6.ProviderServer, error) {
		return &echoProviderServer{}, nil
	}
}

type echoProviderServer struct {
	// The value of the "data" attribute during provider configuration. Will be directly echoed to the echo_test.data attribute.
	providerConfigData tftypes.Value
}

func (e *echoProviderServer) providerSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Description: "This provider is used to output the data attribute provided to the provider configuration into all resources instances of echo_test. " +
				"This is only useful for testing ephemeral resources where the data isn't stored to state.",
			DescriptionKind: tfprotov6.StringKindPlain,
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:            "data",
					Type:            tftypes.DynamicPseudoType,
					Description:     "Dynamic data to provide to the echo_test resource.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Required:        true,
				},
			},
		},
	}
}

func (e *echoProviderServer) testResourceSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:            "data",
					Type:            tftypes.DynamicPseudoType,
					Description:     "Dynamic data that was provided to the provider configuration.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Computed:        true,
				},
			},
		},
	}
}

func (e *echoProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	resp := &tfprotov6.ApplyResourceChangeResponse{}

	if req.TypeName != "echo_test" {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   fmt.Sprintf("ApplyResourceChange was called for a resource type that is not supported by this provider: %q", req.TypeName),
			},
		}

		return resp, nil
	}

	echoTestSchema := e.testResourceSchema()

	plannedState, diag := dynamicValueToValue(echoTestSchema, req.PlannedState)
	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	// Destroy Op, just return planned state, which is null
	if plannedState.IsNull() {
		resp.NewState = req.PlannedState
		return resp, nil
	}

	// Take the provider config "data" attribute verbatim and put back into state. It shares the same type (DynamicPseudoType)
	// as the echo_test "data" attribute.
	newVal := tftypes.NewValue(echoTestSchema.ValueType(), map[string]tftypes.Value{
		"data": e.providerConfigData,
	})

	newState, diag := valuetoDynamicValue(echoTestSchema, newVal)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.NewState = newState

	return resp, nil
}

func (e *echoProviderServer) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	return &tfprotov6.CallFunctionResponse{}, nil
}

func (e *echoProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	resp := &tfprotov6.ConfigureProviderResponse{}

	configVal, diags := dynamicValueToValue(e.providerSchema(), req.Config)
	if diags != nil {
		resp.Diagnostics = append(resp.Diagnostics, diags)
		return resp, nil
	}

	objVal := map[string]tftypes.Value{}
	err := configVal.As(&objVal)
	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading Config",
			Detail:   err.Error(),
		}
		resp.Diagnostics = append(resp.Diagnostics, diag)
		return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
	}

	dynamicDataVal, ok := objVal["data"]
	if !ok {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Config data not found",
		}
		resp.Diagnostics = append(resp.Diagnostics, diag)
		return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
	}

	e.providerConfigData = dynamicDataVal.Copy()

	return resp, nil
}

func (e *echoProviderServer) GetFunctions(ctx context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	return &tfprotov6.GetFunctionsResponse{}, nil
}

func (e *echoProviderServer) GetMetadata(ctx context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	return &tfprotov6.GetMetadataResponse{
		Resources: []tfprotov6.ResourceMetadata{
			{
				TypeName: "echo_test",
			},
		},
	}, nil
}

func (e *echoProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	return &tfprotov6.GetProviderSchemaResponse{
		Provider: e.providerSchema(),
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"echo_test": e.testResourceSchema(),
		},
	}, nil
}

func (e *echoProviderServer) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	return &tfprotov6.ImportResourceStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource Operation",
				Detail:   "ImportResourceState is not supported by this provider.",
			},
		},
	}, nil
}

func (e *echoProviderServer) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	return &tfprotov6.MoveResourceStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource Operation",
				Detail:   "MoveResourceState is not supported by this provider.",
			},
		},
	}, nil
}

func (e *echoProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	// Just return proposed new state, since echo_test doesn't need to modify the plan.
	return &tfprotov6.PlanResourceChangeResponse{
		PlannedState: req.ProposedNewState,
	}, nil
}

func (e *echoProviderServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	return &tfprotov6.ReadDataSourceResponse{}, nil
}

func (e *echoProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	// Just return current state, since the data doesn't need to be refreshed.
	return &tfprotov6.ReadResourceResponse{
		NewState: req.CurrentState,
	}, nil
}

func (e *echoProviderServer) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	return &tfprotov6.StopProviderResponse{}, nil
}

func (e *echoProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}

	if req.TypeName != "echo_test" {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   fmt.Sprintf("UpgradeResourceState was called for a resource type that is not supported by this provider: %q", req.TypeName),
			},
		}

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

	providerSchema := e.providerSchema()

	if req.Version != providerSchema.Version {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   "UpgradeResourceState was called for echo_test, which does not support multiple schema versions",
			},
		}

		return resp, nil
	}

	// Terraform CLI can call UpgradeResourceState even if the stored state
	// version matches the current schema. Presumably this is to account for
	// the previous terraform-plugin-sdk implementation, which handled some
	// state fixups on behalf of Terraform CLI. This will attempt to roundtrip
	// the prior RawState to a state matching the current schema.
	rawStateValue, err := req.RawState.UnmarshalWithOpts(providerSchema.ValueType(), unmarshalOpts)

	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Read Previously Saved State for UpgradeResourceState",
			Detail:   "There was an error reading the saved resource state using the current resource schema: " + err.Error(),
		}

		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
	}

	upgradedState, diag := valuetoDynamicValue(providerSchema, rawStateValue)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.UpgradedState = upgradedState

	return resp, nil
}

func (e *echoProviderServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	return &tfprotov6.ValidateDataResourceConfigResponse{}, nil
}

func (e *echoProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	return &tfprotov6.ValidateProviderConfigResponse{}, nil
}

func (e *echoProviderServer) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	return &tfprotov6.ValidateResourceConfigResponse{}, nil
}

func (e *echoProviderServer) OpenEphemeralResource(ctx context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	return &tfprotov6.OpenEphemeralResourceResponse{}, nil
}

func (e *echoProviderServer) RenewEphemeralResource(ctx context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	return &tfprotov6.RenewEphemeralResourceResponse{}, nil
}

func (e *echoProviderServer) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	return &tfprotov6.CloseEphemeralResourceResponse{}, nil
}

func (e *echoProviderServer) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	return &tfprotov6.ValidateEphemeralResourceConfigResponse{}, nil
}
