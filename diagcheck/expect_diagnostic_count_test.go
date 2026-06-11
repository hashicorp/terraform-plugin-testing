// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package diagcheck_test

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/diagcheck"
)

func TestExpectDiagnosticCount_ZeroCount(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnosticCount("warning", 0)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error for zero count, got: %v", resp.Error)
	}
}

func TestExpectDiagnosticCount_CorrectCount(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnosticCount("warning", 2)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 1",
				Detail:   "Detail 1",
			},
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 2",
				Detail:   "Detail 2",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error for correct count, got: %v", resp.Error)
	}
}

func TestExpectDiagnosticCount_IncorrectCount_TooFew(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnosticCount("warning", 3)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 1",
				Detail:   "Detail 1",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for incorrect count (too few), got nil")
	}
}

func TestExpectDiagnosticCount_IncorrectCount_TooMany(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnosticCount("error", 1)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Error 1",
				Detail:   "Detail 1",
			},
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Error 2",
				Detail:   "Detail 2",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for incorrect count (too many), got nil")
	}
}

func TestExpectDiagnosticCount_OnlyCountsMatchingSeverity(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectDiagnosticCount("warning", 2)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 1",
				Detail:   "Detail 1",
			},
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Error 1",
				Detail:   "Detail 1",
			},
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 2",
				Detail:   "Detail 2",
			},
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Error 2",
				Detail:   "Detail 2",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error when counting only matching severity, got: %v", resp.Error)
	}
}

func TestExpectWarningCount_Zero(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectWarningCount(0)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityError,
				Summary:  "Error 1",
				Detail:   "Detail 1",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error for zero warning count with errors present, got: %v", resp.Error)
	}
}

func TestExpectWarningCount_Correct(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectWarningCount(3)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 1",
				Detail:   "Detail 1",
			},
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 2",
				Detail:   "Detail 2",
			},
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 3",
				Detail:   "Detail 3",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error != nil {
		t.Errorf("expected no error for correct warning count, got: %v", resp.Error)
	}
}

func TestExpectWarningCount_Incorrect(t *testing.T) {
	t.Parallel()

	check := diagcheck.ExpectWarningCount(2)
	req := diagcheck.CheckDiagnosticsRequest{
		Diagnostics: []tfjson.Diagnostic{
			{
				Severity: tfjson.DiagnosticSeverityWarning,
				Summary:  "Warning 1",
				Detail:   "Detail 1",
			},
		},
	}
	resp := &diagcheck.CheckDiagnosticsResponse{}

	check.CheckDiagnostics(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("expected error for incorrect warning count, got nil")
	}
}
