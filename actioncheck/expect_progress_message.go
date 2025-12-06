// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actioncheck

import (
	"context"
	"fmt"
	"strings"
)

var _ ActionCheck = expectProgressMessageContains{}

// expectProgressMessageContains is an ActionCheck that verifies that at least one progress message contains the expected content.
type expectProgressMessageContains struct {
	actionName      string
	expectedContent string
}

// CheckAction implements the ActionCheck interface.
func (e expectProgressMessageContains) CheckAction(ctx context.Context, req CheckActionRequest, resp *CheckActionResponse) {
	if req.ActionName != e.actionName {
		return
	}

	for _, message := range req.Messages {
		if strings.Contains(message.Message, e.expectedContent) {
			return // Found the expected content
		}
	}

	resp.Error = fmt.Errorf("expected action %s to have progress message containing %q, but no matching message found", e.actionName, e.expectedContent)
}

// ExpectProgressMessageContains returns an ActionCheck that verifies that at least one progress message contains the expected content.
func ExpectProgressMessageContains(actionName, expectedContent string) ActionCheck {
	return expectProgressMessageContains{
		actionName:      actionName,
		expectedContent: expectedContent,
	}
}
