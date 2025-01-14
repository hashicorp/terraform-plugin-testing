// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ StateCheck = extractStringValue{}

type extractStringValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	targetVar       *string
}

func (e extractStringValue) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var resource *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")

		return
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")

		return
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")

		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(resource.AttributeValues, e.attributePath)
	if err != nil {
		resp.Error = err

		return
	}

	if result == nil {
		resp.Error = fmt.Errorf("nil: result for attribute '%s' in '%s'", e.attributePath, e.resourceAddress)

		return
	}

	switch t := result.(type) {
	case string:
		*e.targetVar = t
		return
	default:
		resp.Error = fmt.Errorf("invalid type for attribute '%s' in '%s'. Expected: string, Got: %T", e.attributePath, e.resourceAddress, t)

		return
	}
}

func ExtractStringValue(resourceAddress string, attributePath tfjsonpath.Path, targetVar *string) StateCheck {
	return extractStringValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		targetVar:       targetVar,
	}
}
