// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// DiagnosticCheck defines an interface for implementing test logic that checks diagnostics and then returns an error
// if the diagnostics do not match what is expected.
type DiagnosticCheck interface {
	// CheckDiagnostics should perform the diagnostic check.
	CheckDiagnostics(context.Context, CheckDiagnosticsRequest, *CheckDiagnosticsResponse)
}

// CheckDiagnosticsRequest is a request for an invoke of the CheckDiagnostics function.
type CheckDiagnosticsRequest struct {
	// Diagnostics represents diagnostics from terraform validate, containing errors and warnings.
	Diagnostics []tfjson.Diagnostic
}

// CheckDiagnosticsResponse is a response to an invoke of the CheckDiagnostics function.
type CheckDiagnosticsResponse struct {
	// Error is used to report the failure of a diagnostic check assertion and is combined with other DiagnosticCheck errors
	// to be reported as a test failure.
	Error error
}
