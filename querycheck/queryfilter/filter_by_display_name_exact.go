package queryfilter

import (
	"context"
)

type filterByDisplayNameExact struct {
	displayName string
}

func (f filterByDisplayNameExact) Filter(ctx context.Context, req FilterQueryRequest, resp *FilterQueryResponse) {
	if req.QueryItem.DisplayName == f.displayName {
		resp.Include = true
	}
}

func ByDisplayNameExact(displayName string) QueryFilter {
	return filterByDisplayNameExact{
		displayName: displayName,
	}
}
