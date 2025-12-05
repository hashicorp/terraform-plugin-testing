// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

var _ statecheck.StateCheck = &stateCheckSpy{}

type stateCheckSpy struct {
	err    error
	called bool
}

func (s *stateCheckSpy) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	s.called = true
	resp.Error = s.err
}
