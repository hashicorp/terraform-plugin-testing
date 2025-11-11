package queryfilter

import (
	"context"
	"regexp"
)

type filterByDisplayNameRegexp struct {
	regexp *regexp.Regexp
}

func (f filterByDisplayNameRegexp) Filter(ctx context.Context, req FilterQueryRequest, resp *FilterQueryResponse) {
	if f.regexp.MatchString(req.QueryItem.DisplayName) {
		resp.Include = true
	}
}

// ByDisplayNameRegexp returns a query filter that only includes query items that match
// the specified regular expression.
func ByDisplayNameRegexp(regexp *regexp.Regexp) QueryFilter {
	return filterByDisplayNameRegexp{
		regexp: regexp,
	}
}
