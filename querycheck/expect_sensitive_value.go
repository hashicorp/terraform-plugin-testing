// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ QueryCheck = expectSensitiveValue{}

type expectSensitiveValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckQuery implements the query check logic.
func (e expectSensitiveValue) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	var resource *tfjson.QueryResource

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")

		return
	}

	if req.Query.Values == nil {
		resp.Error = fmt.Errorf("query does not contain any query values")

		return
	}

	if req.Query.Values.RootModule == nil {
		resp.Error = fmt.Errorf("query does not contain a root module")

		return
	}

	for _, r := range req.Query.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddress)

		return
	}

	var data map[string]any

	err := json.Unmarshal(resource.SensitiveValues, &data)

	if err != nil {
		resp.Error = fmt.Errorf("could not unmarshal SensitiveValues: %s", err)

		return
	}

	result, err := tfjsonpath.Traverse(data, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	isSensitive, ok := result.(bool)
	if !ok {
		resp.Error = fmt.Errorf("invalid path: the path value cannot be asserted as bool")

		return
	}

	if !isSensitive {
		resp.Error = fmt.Errorf("attribute at path is not sensitive")

		return
	}
}

// ExpectSensitiveValue returns a query check that asserts that the specified attribute at the given resource has a sensitive value.
//
// Due to implementation differences between the terraform-plugin-sdk and the terraform-plugin-framework, representation of sensitive
// values may differ. For example, terraform-plugin-sdk based providers may have less precise representations of sensitive values, such
// as marking whole maps as sensitive rather than individual element values.
func ExpectSensitiveValue(resourceAddress string, attributePath tfjsonpath.Path) QueryCheck {
	return expectSensitiveValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
