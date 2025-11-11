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

// ByDisplayNameExact returns a query filter that only includes query items that match
// the specified display name.
func ByDisplayNameExact(displayName string) QueryFilter {
	return filterByDisplayNameExact{
		displayName: displayName,
	}
}
