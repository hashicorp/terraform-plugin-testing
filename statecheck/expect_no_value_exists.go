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
var _ StateCheck = expectNoValueExists{}

type expectNoValueExists struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckState implements the state check logic.
func (e expectNoValueExists) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
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

	// Resource doesn't exist
	if rc == nil {
		return
	}

	_, err := tfjsonpath.Traverse(rc.AttributeValues, e.attributePath)

	if err == nil {
		resp.Error = fmt.Errorf("attribute found at path: %s.%s", e.resourceAddress, e.attributePath.String())

		return
	}
}

// ExpectNoValueExists returns a state check that asserts that the specified attribute at the given resource
// does not exist.
func ExpectNoValueExists(resourceAddress string, attributePath tfjsonpath.Path) StateCheck {
	return expectNoValueExists{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
