// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actioncheck

import (
	"context"
	"fmt"
	"strings"
)

var _ ActionCheck = expectProgressSequence{}

// expectProgressSequence is an ActionCheck that verifies progress messages appear in the expected sequence.
type expectProgressSequence struct {
	actionName       string
	expectedSequence []string
}

// CheckAction implements the ActionCheck interface.
func (e expectProgressSequence) CheckAction(ctx context.Context, req CheckActionRequest, resp *CheckActionResponse) {
	if req.ActionName != e.actionName {
		return
	}

	sequenceIndex := 0
	for _, message := range req.Messages {
		if sequenceIndex < len(e.expectedSequence) && strings.Contains(message.Message, e.expectedSequence[sequenceIndex]) {
			sequenceIndex++
		}
	}

	if sequenceIndex != len(e.expectedSequence) {
		resp.Error = fmt.Errorf("expected action %s to have progress messages in sequence %v, but only found %d of %d expected messages",
			e.actionName, e.expectedSequence, sequenceIndex, len(e.expectedSequence))
	}
}

// ExpectProgressSequence returns an ActionCheck that verifies progress messages appear in the expected sequence.
func ExpectProgressSequence(actionName string, expectedSequence []string) ActionCheck {
	return expectProgressSequence{
		actionName:       actionName,
		expectedSequence: expectedSequence,
	}
}
