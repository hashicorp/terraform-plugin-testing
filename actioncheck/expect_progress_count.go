// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actioncheck

import (
	"context"
	"fmt"
)

var _ ActionCheck = expectProgressCount{}

// expectProgressCount is an ActionCheck that verifies the expected number of progress messages.
type expectProgressCount struct {
	actionName    string
	expectedCount int
}

// CheckAction implements the ActionCheck interface.
func (e expectProgressCount) CheckAction(ctx context.Context, req CheckActionRequest, resp *CheckActionResponse) {
	if req.ActionName != e.actionName {
		return
	}

	actualCount := len(req.Messages)
	if actualCount != e.expectedCount {
		resp.Error = fmt.Errorf("expected action %s to have %d progress messages, but got %d", e.actionName, e.expectedCount, actualCount)
	}
}

// ExpectProgressCount returns an ActionCheck that verifies the expected number of progress messages.
func ExpectProgressCount(actionName string, expectedCount int) ActionCheck {
	return expectProgressCount{
		actionName:    actionName,
		expectedCount: expectedCount,
	}
}
