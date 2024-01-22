// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = expectValueExists{}

type expectValueExists struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckState implements the state check logic.
func (e expectValueExists) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var rc *tfjson.StateResource

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

	for _, resourceChange := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == resourceChange.Address {
			rc = resourceChange

			break
		}
	}

	if rc == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	_, err := tfjsonpath.Traverse(rc.AttributeValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}
}

// ExpectValueExists returns a state check that asserts that the specified
// attribute at the given resource exists.
func ExpectValueExists(resourceAddress string, attributePath tfjsonpath.Path) StateCheck {
	return expectValueExists{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
