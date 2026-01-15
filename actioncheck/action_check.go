// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actioncheck

import (
	"context"
	"time"
)

// ActionCheck defines an interface for implementing test logic that checks action progress messages.
type ActionCheck interface {
	// CheckAction should perform the action check.
	CheckAction(context.Context, CheckActionRequest, *CheckActionResponse)
}

// CheckActionRequest is a request for an invoke of the CheckAction function.
type CheckActionRequest struct {
	// ActionName is the name of the action being checked (e.g., "aws_lambda_invoke.test").
	ActionName string

	// Messages contains all progress messages captured for this action.
	Messages []ProgressMessage
}

// CheckActionResponse is a response to an invoke of the CheckAction function.
type CheckActionResponse struct {
	// Error is used to report the failure of an action check assertion.
	Error error
}

// ProgressMessage represents a single progress message from an action.
type ProgressMessage struct {
	// Message is the progress message content.
	Message string

	// Timestamp is when the message was captured.
	Timestamp time.Time
}
