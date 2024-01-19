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

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for attribute at path: %s.%s, err: %s", e.resourceAddress, e.attributePath.String(), err)
	}
}

// ExpectKnownValue returns a state check that asserts that the specified attribute at the given resource
// has a known type and value.
//
// The following is an example of using ExpectKnownValue.
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
//	func TestExpectKnownValue_CheckState_Bool(t *testing.T) {
//		t.Parallel()
//
//		resource.Test(t, resource.TestCase{
//			// Provider definition omitted.
//			Steps: []resource.TestStep{
//				{
//					Config: `resource "test_resource" "one" {
//		          bool_attribute = true
//		        }
//		        `,
//					ConfigStateChecks: resource.ConfigStateChecks{
//						statecheck.ExpectKnownValue(
//							"test_resource.one",
//							tfjsonpath.New("bool_attribute"),
//							knownvalue.BoolExact(true),
//						),
//					},
//				},
//			},
//		})
//	}
//
// The following is an example of using ExpectKnownValue in combination
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
//	func TestExpectKnownValue_CheckState_String_Custom(t *testing.T) {
//		t.Parallel()
//
//		resource.Test(t, resource.TestCase{
//			// Provider definition omitted.
//			Steps: []resource.TestStep{
//				{
//					Config: `resource "test_resource" "one" {
//		           string_attribute = "string"
//		         }
//		         `,
//					ConfigStateChecks: resource.ConfigStateChecks{
//						statecheck.ExpectKnownValue(
//							"test_resource.one",
//							tfjsonpath.New("string_attribute"),
//							StringContains("tri")),
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
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) StateCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
