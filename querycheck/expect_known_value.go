package querycheck

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"strings"
)

var _ QueryResultCheck = expectKnownValue{}

type expectKnownValue struct {
	listResourceAddress string
	resourceName        string
	attributePath       tfjsonpath.Path
	knownValue          knownvalue.Check
}

func (e expectKnownValue) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	for _, res := range *req.Query {
		diags := make([]error, 0)

		if e.listResourceAddress == strings.TrimPrefix(res.Address, "list.") && e.resourceName == res.DisplayName {
			if res.ResourceObject == nil {
				resp.Error = fmt.Errorf("%s - no resource object was returned, ensure `include_resource` has been set to `true` in the list resource config`", e.listResourceAddress)
				return
			}

			// Ideally we can do the check like below which is identical to the expect known value state check but... terraform-json hasn't
			// defined the resource object as a map[string]interface, we so we need to iterate over it manually
			//resource, err := tfjsonpath.Traverse(res.ResourceObject, e.attributePath)
			//if err != nil {
			//	resp.Error = err
			//	return
			//}
			//
			//if err := e.knownValue.CheckValue(resource); err != nil {
			//	diags = append(diags, fmt.Errorf("error checking value for attribute at path: %s for resource %s, err: %s", e.attributePath.String(), e.resourceName, err))
			//}
			//
			//if diags == nil {
			//	return
			//}

			for k, v := range res.ResourceObject {
				if k == e.attributePath.String() {
					var val any

					err := json.Unmarshal(v, &val)
					if err != nil {
						resp.Error = fmt.Errorf("%s - Error decoding message type: %s", e.listResourceAddress, err)
						return
					}

					if err := e.knownValue.CheckValue(val); err != nil {
						diags = append(diags, fmt.Errorf("error checking value for attribute at path: %s for resource %s, err: %s", e.attributePath.String(), e.resourceName, err))
					}
				}
			}
		}
	}

	resp.Error = fmt.Errorf("%s - the resource %s was not found", e.listResourceAddress, e.resourceName)

	return
}

func ExpectKnownValue(listResourceAddress string, resourceName string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) QueryResultCheck {
	return expectKnownValue{
		listResourceAddress: listResourceAddress,
		resourceName:        resourceName,
		attributePath:       attributePath,
		knownValue:          knownValue,
	}
}
