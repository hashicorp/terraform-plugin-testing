// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package query

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

var _ querycheck.QueryCheck = &queryCheckSpy{}

type queryCheckSpy struct {
	err    error
	called bool
}

func (s *queryCheckSpy) CheckQuery(ctx context.Context, req querycheck.CheckQueryRequest, resp *querycheck.CheckQueryResponse) {
	s.called = true
	resp.Error = s.err
}
