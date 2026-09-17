// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck_test

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/diagcheck"
)

func TestExpectNoDiagnostics_NoDiagnostics(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoDiagnostics()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error, got: %v", resp.Error)
	}
}

func TestExpectNoDiagnostics_WithWarning_Error(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoDiagnostics()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Test Warning",
				Detail:   "This is a test warning",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for warning diagnostic, got nil")
	}
}

func TestExpectNoDiagnostics_WithError_Error(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoDiagnostics()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Test Error",
				Detail:   "This is a test error",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for error diagnostic, got nil")
	}
}

func TestExpectNoDiagnostics_WithMultiple_Error(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoDiagnostics()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 1",
				Detail:   "First warning",
			},
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Error 1",
				Detail:   "First error",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for multiple diagnostics, got nil")
	}
}

func TestExpectNoWarnings_NoWarnings(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoWarnings()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Test Error",
				Detail:   "This is a test error",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error when only errors present, got: %v", resp.Error)
	}
}

func TestExpectNoWarnings_WithWarning_Error(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoWarnings()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Test Warning",
				Detail:   "This is a test warning",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for warning diagnostic, got nil")
	}
}

func TestExpectNoErrors_NoErrors(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoErrors()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Test Warning",
				Detail:   "This is a test warning",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error when only warnings present, got: %v", resp.Error)
	}
}

func TestExpectNoErrors_WithError_Error(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectNoErrors()
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Test Error",
				Detail:   "This is a test error",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for error diagnostic, got nil")
	}
}
