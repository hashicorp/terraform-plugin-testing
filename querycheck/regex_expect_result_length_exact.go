// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"regexp"
)

var _ QueryResultCheck = regexExpectLength{}

type regexExpectLength struct {
	regex *regexp.Regexp
	check int
}

// CheckQuery implements the query check logic.
func (e regexExpectLength) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil && req.QuerySummaries == nil {
		resp.Error = fmt.Errorf("no query summary information available")
		return
	}

	total := 0
	for _, summary := range req.QuerySummaries {
		matches, err := regexp.MatchString(e.regex.String(), summary.Address)
		if err != nil {
			resp.Error = fmt.Errorf("invalid regex pattern provided: %s, error: %s", e.regex.String(), err)
			return
		}

		if matches {
			total += summary.Total
		}
	}

	if total == e.check {
		return
	}

	resp.Error = fmt.Errorf("number of found resources matching regex %s - expected %v but got %v.", e.regex.String(), e.check, total)
}

// ExpectLength returns a query check that asserts that the length of the query result is exactly the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func RegexExpectLength(regex *regexp.Regexp, length int) QueryResultCheck {
	return regexExpectLength{
		regex: regex,
		check: length,
	}
}
