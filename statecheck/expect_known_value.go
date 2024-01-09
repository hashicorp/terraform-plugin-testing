// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = expectKnownValue{}

type expectKnownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	knownValue      knownvalue.Check
}

// CheckState implements the state check logic.
func (e expectKnownValue) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var rc *tfjson.StateResource

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
		if e.resourceAddress == resourceChange.Address {
			rc = resourceChange

			break
		}
	}

	if rc == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(rc.AttributeValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if result == nil {
		resp.Error = fmt.Errorf("value is null")

		return
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Bool,
		reflect.Map,
		reflect.Slice,
		reflect.String:
		if err := e.knownValue.CheckValue(result); err != nil {
			resp.Error = err

			return
		}
	default:
		errorStr := fmt.Sprintf("unrecognised attribute type: %T, known value type is %T", result, e.knownValue)
		errorStr += "\n\nThis is an error in statecheck.ExpectKnownValue.\nPlease report this to the maintainers."

		resp.Error = fmt.Errorf(errorStr)

		return
	}
}

// ExpectKnownValue returns a state check that asserts that the specified attribute at the given resource
// has a known type and value.
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) StateCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
