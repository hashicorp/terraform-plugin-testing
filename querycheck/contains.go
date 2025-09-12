package querycheck

import (
	"context"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

var _ QueryResultCheck = contains{}

type contains struct {
	resourceAddress string
	check           string
}

func (c contains) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	// TODO refactor below
	foundResources := make([]string, 0)

	if req.Query == nil {
		resp.Error = fmt.Errorf("Query is nil")
		return
	}

	for _, v := range *req.Query {
		switch i := v.(type) {
		case tfjson.ListResourceFoundMessage:
			prefix := "list."
			if strings.TrimPrefix(i.ListResourceFound.Address, prefix) == c.resourceAddress {
				foundResources = append(foundResources, i.ListResourceFound.DisplayName)
			}
		default:
			continue
		}
	}

	if len(foundResources) == 0 {
		resp.Error = fmt.Errorf("%s - no resources found by query.", c.resourceAddress)

		return
	}
	// TODO refactor above

	for _, res := range foundResources {
		if c.check == res {
			return
		}
	}

	resp.Error = fmt.Errorf("expected to find resource with display name %q in results but resource was not found", c.check)

	return
}

// Contains returns a query check that asserts that a resource with a given display name exists within the returned results of the query.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func Contains(resourceAddress string, displayName string) QueryResultCheck {
	return contains{
		resourceAddress: resourceAddress,
		check:           displayName,
	}
}
