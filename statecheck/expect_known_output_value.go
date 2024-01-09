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
var _ StateCheck = expectKnownOutputValue{}

type expectKnownOutputValue struct {
	outputAddress string
	knownValue    knownvalue.Check
}

// CheckState implements the state check logic.
func (e expectKnownOutputValue) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var output *tfjson.StateOutput

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")
	}

	for address, oc := range req.State.Values.Outputs {
		if e.outputAddress == address {
			output = oc

			break
		}
	}

	if output == nil {
		resp.Error = fmt.Errorf("%s - Output not found in state", e.outputAddress)

		return
	}

	result, err := tfjsonpath.Traverse(output.Value, tfjsonpath.Path{})

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
		errorStr := fmt.Sprintf("unrecognised output type: %T, known value type is %T", result, e.knownValue)
		errorStr += "\n\nThis is an error in statecheck.ExpectKnownOutputValue.\nPlease report this to the maintainers."

		resp.Error = fmt.Errorf(errorStr)

		return
	}
}

// ExpectKnownOutputValue returns a state check that asserts that the specified value
// has a known type, and value.
func ExpectKnownOutputValue(outputAddress string, knownValue knownvalue.Check) StateCheck {
	return expectKnownOutputValue{
		outputAddress: outputAddress,
		knownValue:    knownValue,
	}
}
