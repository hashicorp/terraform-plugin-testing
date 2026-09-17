// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck

import (
	"context"
	"fmt"
)

var _ DiagnosticCheck = expectDiagnosticCount{}

type expectDiagnosticCount struct {
	severity string
	count    int
}

// CheckDiagnostics implements the diagnostic check logic.
func (e expectDiagnosticCount) CheckDiagnostics(ctx context.Context, req CheckDiagnosticsRequest, resp *CheckDiagnosticsResponse) {
	actualCount := 0

	for _, diag := range req.Diagnostics {
		if string(diag.Severity) == e.severity {
			actualCount++
		}
	}

	if actualCount != e.count {
		resp.Error = fmt.Errorf("expected %d %s diagnostic(s), got %d", e.count, e.severity, actualCount)
	}
}

// ExpectDiagnosticCount returns a diagnostic check that asserts that there are exactly the specified number of
// diagnostics with the given severity. The severity should be "error" or "warning".
func ExpectDiagnosticCount(severity string, count int) DiagnosticCheck {
	return expectDiagnosticCount{
		severity: severity,
		count:    count,
	}
}

// ExpectWarningCount returns a diagnostic check that asserts that there are exactly the specified number of
// warning diagnostics. This is a convenience wrapper for ExpectDiagnosticCount with severity="warning".
func ExpectWarningCount(count int) DiagnosticCheck {
	return expectDiagnosticCount{
		severity: "warning",
		count:    count,
	}
}
