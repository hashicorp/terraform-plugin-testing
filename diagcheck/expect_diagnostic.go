// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck

import (
	"context"
	"fmt"
)

var _ DiagnosticCheck = expectDiagnostic{}

type expectDiagnostic struct {
	severity string
	summary  string
	detail   string
}

// CheckDiagnostics implements the diagnostic check logic.
func (e expectDiagnostic) CheckDiagnostics(ctx context.Context, req CheckDiagnosticsRequest, resp *CheckDiagnosticsResponse) {
	for _, diag := range req.Diagnostics {
		if string(diag.Severity) == e.severity &&
			diag.Summary == e.summary &&
			diag.Detail == e.detail {
			// Found matching diagnostic
			return
		}
	}

	// Diagnostic not found
	resp.Error = fmt.Errorf("expected %s diagnostic not found: summary=%q detail=%q", e.severity, e.summary, e.detail)
}

// ExpectDiagnostic returns a diagnostic check that asserts that a specific diagnostic exists with the given
// severity, summary, and detail. The severity should be "error" or "warning".
func ExpectDiagnostic(severity, summary, detail string) DiagnosticCheck {
	return expectDiagnostic{
		severity: severity,
		summary:  summary,
		detail:   detail,
	}
}

// ExpectWarning returns a diagnostic check that asserts that a specific warning diagnostic exists with the given
// summary and detail. This is a convenience wrapper for ExpectDiagnostic with severity="warning".
func ExpectWarning(summary, detail string) DiagnosticCheck {
	return expectDiagnostic{
		severity: "warning",
		summary:  summary,
		detail:   detail,
	}
}
