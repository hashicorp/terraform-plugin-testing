// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package queryfilter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

type filterByResourceIdentityFunc struct {
	identityFunc func() map[string]knownvalue.Check
}

func (f filterByResourceIdentityFunc) Filter(ctx context.Context, req FilterQueryRequest, resp *FilterQueryResponse) {
	inner := ByResourceIdentity(f.identityFunc())
	inner.Filter(ctx, req, resp)
}

// ByResourceIdentityFunc returns a query filter that only includes query items that match
// the given resource identity.
//
// Errors thrown by the given known value checks are only used to filter out non-matching query
// items and are otherwise ignored.
func ByResourceIdentityFunc(identityFunc func() map[string]knownvalue.Check) QueryFilter {
	return filterByResourceIdentityFunc{
		identityFunc: identityFunc,
	}
}
