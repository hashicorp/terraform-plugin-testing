// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"
	"sort"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = &compareValueContains{}

type compareValueContains struct {
	resourceAddressOne string
	attributePathOne   tfjsonpath.Path
	resourceAddressTwo string
	attributePathTwo   tfjsonpath.Path
	comparer           compare.ValueComparer
}

// CheckState implements the state check logic.
func (e *compareValueContains) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var resourceOne *tfjson.StateResource
	var resourceTwo *tfjson.StateResource

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

	if listOrSet, ok := resultOne.([]any); ok {
		var errs []error

		for _, v := range listOrSet {
			errs = append(errs, e.comparer.CompareValues(v, resultTwo))
		}

		for _, err = range errs {
			if err == nil {
				return
			}
		}

		resp.Error = err

		return
	}

	if mapOrObject, ok := resultOne.(map[string]any); ok {
		keys := make([]string, 0, len(mapOrObject))

		for k := range mapOrObject {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		var errs []error

		for _, key := range keys {
			errs = append(errs, e.comparer.CompareValues(mapOrObject[key], resultTwo))
		}

		for _, err = range errs {
			if err == nil {
				return
			}
		}

		resp.Error = err

		return
	}

	resp.Error = fmt.Errorf("expected []any or map[string]any value for CompareValueContains check, got: %T", resultOne)
}

func CompareValueContains(resourceAddressOne string, attributePathOne tfjsonpath.Path, resourceAddressTwo string, attributePathTwo tfjsonpath.Path, comparer compare.ValueComparer) StateCheck {
	return &compareValueContains{
		resourceAddressOne: resourceAddressOne,
		attributePathOne:   attributePathOne,
		resourceAddressTwo: resourceAddressTwo,
		attributePathTwo:   attributePathTwo,
		comparer:           comparer,
	}
}
