// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = expectKnownOutputValueAtPath{}

type expectKnownOutputValueAtPath struct {
	outputAddress string
	outputPath    tfjsonpath.Path
	knownValue    knownvalue.Check
}

// CheckState implements the state check logic.
func (e expectKnownOutputValueAtPath) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var output *tfjson.StateOutput

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")

		return
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")

		return
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

	result, err := tfjsonpath.Traverse(output.Value, e.outputPath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for output at path: %s.%s, err: %s", e.outputAddress, e.outputPath.String(), err)

		return
	}
}

// ExpectKnownOutputValueAtPath returns a state check that asserts that the specified output at the given path
// has a known type and value.
//
// The following is an example of using ExpectKnownOutputValueAtPath.
//
//	package example_test
//
//	import (
//		"testing"
//
//		"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//		"github.com/hashicorp/terraform-plugin-testing/knownvalue"
//		"github.com/hashicorp/terraform-plugin-testing/statecheck"
//		"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
//	)
//
//	func TestExpectKnownOutputValueAtPath_CheckState_Bool(t *testing.T) {
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
//		        output test_resource_one_output {
//		          value = test_resource.one
//		        }
//		        `,
//					ConfigStateChecks: resource.ConfigStateChecks{
//						statecheck.ExpectKnownOutputValueAtPath(
//							"test_resource_one_output",
//							tfjsonpath.New("bool_attribute"),
//							knownvalue.BoolExact(true),
//						),
//					},
//				},
//			},
//		})
//	}
//
// The following is an example of using ExpectKnownOutputValueAtPath in combination
// with a custom knownvalue.Check.
//
//	package example_test
//
//	import (
//		"fmt"
//		"strings"
//		"testing"
//
//		"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//		"github.com/hashicorp/terraform-plugin-testing/knownvalue"
//		"github.com/hashicorp/terraform-plugin-testing/statecheck"
//		"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
//	)
//
//	func TestExpectKnownOutputValueAtPath_CheckState_String_Custom(t *testing.T) {
//		t.Parallel()
//
//		resource.Test(t, resource.TestCase{
//			// Provider definition omitted.
//			Steps: []resource.TestStep{
//				{
//					Config: `resource "test_resource" "one" {
//							string_attribute = "string"
//						}
//
//						output test_resource_one_output {
//							value = test_resource.one
//						}
//						`,
//					ConfigStateChecks: resource.ConfigStateChecks{
//						statecheck.ExpectKnownOutputValueAtPath(
//							"test_resource_one_output",
//							tfjsonpath.New("string_attribute"),
//							StringContains("str")),
//					},
//				},
//			},
//		})
//	}
//
//	var _ knownvalue.Check = stringContains{}
//
//	type stringContains struct {
//		value string
//	}
//
//	func (v stringContains) CheckValue(other any) error {
//		otherVal, ok := other.(string)
//
//		if !ok {
//			return fmt.Errorf("expected string value for StringContains check, got: %T", other)
//		}
//
//		if !strings.Contains(otherVal, v.value) {
//			return fmt.Errorf("expected string %q to contain %q for StringContains check", otherVal, v.value)
//		}
//
//		return nil
//	}
//
//	func (v stringContains) String() string {
//		return v.value
//	}
//
//	func StringContains(value string) stringContains {
//		return stringContains{
//			value: value,
//		}
//	}
func ExpectKnownOutputValueAtPath(outputAddress string, outputPath tfjsonpath.Path, knownValue knownvalue.Check) StateCheck {
	return expectKnownOutputValueAtPath{
		outputAddress: outputAddress,
		outputPath:    outputPath,
		knownValue:    knownValue,
	}
}
