// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ QueryResultCheck = expectIdentity{}

type expectIdentity struct {
	resourceAddress string
	check           map[string]knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectIdentity) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	for _, res := range *req.Query {
		var errCollection []error

		for attribute := range e.check {
			var val any
			var unmarshalledVal any

			val, ok := res.Identity[attribute]
			if !ok {
				resp.Error = fmt.Errorf("%s - expected attribute %q not in actual identity object", e.resourceAddress, attribute)
				return
			}

			rawMessage, ok := val.(json.RawMessage)
			if !ok {
				resp.Error = fmt.Errorf("%s - expected json.RawMessage but got %T", e.resourceAddress, val)
				return
			}
			err := json.Unmarshal(rawMessage, &unmarshalledVal)

			if err != nil {
				resp.Error = fmt.Errorf("%s - Error decoding message type: %s", e.resourceAddress, err)
				return
			}

			if err = e.check[attribute].CheckValue(unmarshalledVal); err != nil {
				errCollection = append(errCollection, fmt.Errorf("%s - %q identity attribute: %s\n", e.resourceAddress, e.check, err))
			}
		}
		if errCollection == nil {
			return
		}
	}

	var errCollection []error

	errCollection = append(errCollection, fmt.Errorf("An identity with the following attributes was not found:"))

	// wrap errors for each check
	for attr, check := range e.check {
		errCollection = append(errCollection, fmt.Errorf("Attribute %s: %s", attr, check))
	}
	errCollection = append(errCollection, fmt.Errorf("Address: %s\n", e.resourceAddress))

	resp.Error = errors.Join(errCollection...)

	return
}

// ExpectIdentity returns a query check that asserts that the identity at the given resource matches a known object, where each
// map key represents an identity attribute name. The identity in query must exactly match the given object.
//
// This query check can only be used with managed resources that support resource identity and query. Query is only supported in Terraform v1.14+
func ExpectIdentity(resourceAddress string, identity map[string]knownvalue.Check) QueryResultCheck {
	return expectIdentity{
		resourceAddress: resourceAddress,
		check:           identity,
	}
}
