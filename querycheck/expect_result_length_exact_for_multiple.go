// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"regexp"
)

var _ QueryResultCheck = expectLengthForMultiple{}

type expectLengthForMultiple struct {
	regex *regexp.Regexp
	check int
}

// CheckQuery implements the query check logic.
func (e expectLengthForMultiple) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil && req.QuerySummaries == nil {
		resp.Error = fmt.Errorf("no query summary information available")
		return
	}

	total := 0
	matchFound := false
	for _, summary := range req.QuerySummaries {
		matches, err := regexp.MatchString(e.regex.String(), summary.Address)
		if err != nil {
			resp.Error = fmt.Errorf("invalid regex pattern provided: %s, error: %s", e.regex.String(), err)
			return
		}

		if matches {
			total += summary.Total
			matchFound = true
		}
	}

	if !matchFound {
		resp.Error = fmt.Errorf("no list resources matching the provided regex pattern %s were found in the query results", e.regex.String())
		return
	}

	if total != e.check {
		resp.Error = fmt.Errorf("number of found resources %v - expected but got %v.", e.check, total)
	}
}

// ExpectLengthForMultiple returns a query check that asserts that the sum of query result lengths
// produced by multiple list blocks is exactly the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectLengthForMultiple(regex *regexp.Regexp, length int) QueryResultCheck {
	return expectLengthForMultiple{
		regex: regex,
		check: length,
	}
}
