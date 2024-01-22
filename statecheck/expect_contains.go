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
var _ StateCheck = expectContains{}

type expectContains struct {
	resourceAddressOne string
	attributePathOne   tfjsonpath.Path
	resourceAddressTwo string
	attributePathTwo   tfjsonpath.Path
}

// CheckState implements the state check logic.
func (e expectContains) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var resourceOne, resourceTwo *tfjson.StateResource

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
		if e.resourceAddressOne == r.Address {
			resourceOne = r

			break
		}
	}

	if resourceOne == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddressOne)

		return
	}

	resultOne, err := tfjsonpath.Traverse(resourceOne.AttributeValues, e.attributePathOne)

	if err != nil {
		resp.Error = err

		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddressTwo == r.Address {
			resourceTwo = r

			break
		}
	}

	if resourceTwo == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddressTwo)

		return
	}

	resultTwo, err := tfjsonpath.Traverse(resourceTwo.AttributeValues, e.attributePathTwo)

	if err != nil {
		resp.Error = err

		return
	}

	resultOneCollection, ok := resultOne.([]any)

	if !ok {
		resp.Error = fmt.Errorf("%s.%s is not a collection", e.resourceAddressOne, e.attributePathOne.String())

		return
	}

	for _, v := range resultOneCollection {
		if reflect.DeepEqual(v, resultTwo) {
			return
		}
	}

	resp.Error = fmt.Errorf("value of %s.%s is not found in value of %s.%s", e.resourceAddressTwo, e.attributePathTwo.String(), e.resourceAddressOne, e.attributePathOne.String())
}

// ExpectContains returns a state check that asserts that the value of the second attribute is contained within the
// value of the first attribute, allowing checking of whether a set contains a value identified by another attribute.
func ExpectContains(resourceAddressOne string, attributePathOne tfjsonpath.Path, resourceAddressTwo string, attributePathTwo tfjsonpath.Path) StateCheck {
	return expectContains{
		resourceAddressOne: resourceAddressOne,
		attributePathOne:   attributePathOne,
		resourceAddressTwo: resourceAddressTwo,
		attributePathTwo:   attributePathTwo,
	}
}
