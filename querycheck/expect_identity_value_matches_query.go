// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ QueryCheck = expectIdentityValueMatchesQuery{}

type expectIdentityValueMatchesQuery struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

// CheckQuery implements the query check logic.
func (e expectIdentityValueMatchesQuery) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
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

	if resource.IdentitySchemaVersion == nil || len(resource.IdentityValues) == 0 {
		resp.Error = fmt.Errorf("%s - Identity not found in query. Either the resource does not support identity or the Terraform version running the test does not support identity. (must be v1.12+)", e.resourceAddress)

		return
	}

	identityResult, err := tfjsonpath.Traverse(resource.IdentityValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	queryResult, err := tfjsonpath.Traverse(resource.AttributeValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if !reflect.DeepEqual(identityResult, queryResult) {
		resp.Error = fmt.Errorf("expected identity and query value at path to match, but they differ: %s.%s, identity value: %v, query value: %v", e.resourceAddress, e.attributePath.String(), identityResult, queryResult)

		return
	}
}

// ExpectIdentityValueMatchesQuery returns a query check that asserts that the specified identity attribute at the given resource
// matches the same attribute in query. This is useful when an identity attribute is in sync with a query attribute of the same path.
//
// This query check can only be used with managed resources that support resource identity. Resource identity is only supported in Terraform v1.12+
func ExpectIdentityValueMatchesQuery(resourceAddress string, attributePath tfjsonpath.Path) QueryCheck {
	return expectIdentityValueMatchesQuery{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
