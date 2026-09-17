// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck

import (
	"context"
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

var _ DiagnosticCheck = expectNoDiagnostics{}

type expectNoDiagnostics struct {
	allowErrors   bool
	allowWarnings bool
}

// CheckDiagnostics implements the diagnostic check logic.
func (e expectNoDiagnostics) CheckDiagnostics(ctx context.Context, req CheckDiagnosticsRequest, resp *CheckDiagnosticsResponse) {
	var result []error

	for _, diag := range req.Diagnostics {
		if diag.Severity == tfjson.DiagnosticSeverityError && !e.allowErrors {
			result = append(result, fmt.Errorf("unexpected error diagnostic: %s - %s", diag.Summary, diag.Detail))
		}
		if diag.Severity == tfjson.DiagnosticSeverityWarning && !e.allowWarnings {
			result = append(result, fmt.Errorf("unexpected warning diagnostic: %s - %s", diag.Summary, diag.Detail))
		}
	}

	resp.Error = errors.Join(result...)
}

// ExpectNoDiagnostics returns a diagnostic check that asserts that there are no diagnostics (errors or warnings).
// All diagnostics found will be aggregated and returned in a diagnostic check error.
func ExpectNoDiagnostics() DiagnosticCheck {
	return expectNoDiagnostics{
		allowErrors:   false,
		allowWarnings: false,
	}
}

// ExpectNoWarnings returns a diagnostic check that asserts that there are no warning diagnostics.
// Errors are allowed. All warnings found will be aggregated and returned in a diagnostic check error.
func ExpectNoWarnings() DiagnosticCheck {
	return expectNoDiagnostics{
		allowErrors:   true,
		allowWarnings: false,
	}
}

// ExpectNoErrors returns a diagnostic check that asserts that there are no error diagnostics.
// Warnings are allowed. All errors found will be aggregated and returned in a diagnostic check error.
func ExpectNoErrors() DiagnosticCheck {
	return expectNoDiagnostics{
		allowErrors:   false,
		allowWarnings: true,
	}
}
