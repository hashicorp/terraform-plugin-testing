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
//
// The following is an example of using ExpectMatchingValues.
//
//	package example_test
//
//	import (
//		"testing"
//
//		"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//		"github.com/hashicorp/terraform-plugin-testing/statecheck"
//		"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
//	)
//
//	func TestExpectMatchingValues_CheckState_AttributeValuesEqual_Bool(t *testing.T) {
//		t.Parallel()
//
//		resource.Test(t, resource.TestCase{
//			// Provider definition omitted.
//			Steps: []resource.TestStep{
//				{
//					Config: `resource "test_resource" "one" {
//		          bool_attribute = true
//		        }
//
//		        resource "test_resource" "two" {
//		          bool_attribute = true
//		        }`,
//					ConfigStateChecks: resource.ConfigStateChecks{
//						statecheck.ExpectMatchingValues(
//							"test_resource.one",
//							tfjsonpath.New("bool_attribute"),
//							"test_resource.two",
//							tfjsonpath.New("bool_attribute"),
//						),
//					},
//				},
//			},
//		})
//	}
func ExpectMatchingValues(resourceAddressOne string, attributePathOne tfjsonpath.Path, resourceAddressTwo string, attributePathTwo tfjsonpath.Path) StateCheck {
	return expectMatchingValues{
		resourceAddressOne: resourceAddressOne,
		attributePathOne:   attributePathOne,
		resourceAddressTwo: resourceAddressTwo,
		attributePathTwo:   attributePathTwo,
	}
}
