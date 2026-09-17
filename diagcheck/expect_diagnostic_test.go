// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck_test

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/diagcheck"
)

func TestExpectDiagnostic_Found(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnostic("warning", "Test Summary", "Test Detail")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Test Summary",
				Detail:   "Test Detail",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error when diagnostic found, got: %v", resp.Error)
	}
}

func TestExpectDiagnostic_NotFound(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnostic("warning", "Expected Summary", "Expected Detail")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Different Summary",
				Detail:   "Different Detail",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error when diagnostic not found, got nil")
	}
}

func TestExpectDiagnostic_WrongSeverity(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnostic("error", "Test Summary", "Test Detail")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Test Summary",
				Detail:   "Test Detail",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error when severity doesn't match, got nil")
	}
}

func TestExpectDiagnostic_FoundAmongMultiple(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnostic("warning", "Second Warning", "Second Detail")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "First Warning",
				Detail:   "First Detail",
			},
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Second Warning",
				Detail:   "Second Detail",
			},
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "An Error",
				Detail:   "Error Detail",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error when diagnostic found among multiple, got: %v", resp.Error)
	}
}

func TestExpectDiagnostic_EmptyList(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnostic("warning", "Test Summary", "Test Detail")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error when diagnostic list is empty, got nil")
	}
}

func TestExpectWarning_Found(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectWarning("Deprecation Warning", "Use new_field instead")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Deprecation Warning",
				Detail:   "Use new_field instead",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error when warning found, got: %v", resp.Error)
	}
}

func TestExpectWarning_ErrorInsteadOfWarning(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectWarning("Test Summary", "Test Detail")
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Test Summary",
				Detail:   "Test Detail",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error when severity is error instead of warning, got nil")
	}
}
