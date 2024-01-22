// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = expectMatchingValues{}

type expectMatchingValues struct {
	resourceAddressOne string
	attributePathOne   tfjsonpath.Path
	resourceAddressTwo string
	attributePathTwo   tfjsonpath.Path
}

// CheckState implements the state check logic.
func (e expectMatchingValues) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var rcOne, rcTwo *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")
	}

	for _, resourceChange := range req.State.Values.RootModule.Resources {
		if e.resourceAddressOne == resourceChange.Address {
			rcOne = resourceChange

			break
		}
	}

	if rcOne == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddressOne)

		return
	}

	resultOne, err := tfjsonpath.Traverse(rcOne.AttributeValues, e.attributePathOne)

	if err != nil {
		resp.Error = err

		return
	}

	for _, resourceChange := range req.State.Values.RootModule.Resources {
		if e.resourceAddressTwo == resourceChange.Address {
			rcTwo = resourceChange

			break
		}
	}

	if rcTwo == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddressTwo)

		return
	}

	resultTwo, err := tfjsonpath.Traverse(rcTwo.AttributeValues, e.attributePathTwo)

	if err != nil {
		resp.Error = err

		return
	}

	if !reflect.DeepEqual(resultOne, resultTwo) {
		resp.Error = fmt.Errorf("values are not equal: %v != %v", resultOne, resultTwo)

		return
	}
}

// ExpectMatchingValues returns a state check that asserts that the specified attributes at the given resources
// have a matching value.
func ExpectMatchingValues(resourceAddressOne string, attributePathOne tfjsonpath.Path, resourceAddressTwo string, attributePathTwo tfjsonpath.Path) StateCheck {
	return expectMatchingValues{
		resourceAddressOne: resourceAddressOne,
		attributePathOne:   attributePathOne,
		resourceAddressTwo: resourceAddressTwo,
		attributePathTwo:   attributePathTwo,
	}
}
