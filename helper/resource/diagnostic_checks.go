// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/diagcheck"
	"github.com/mitchellh/go-testing-interface"
)

func runDiagnosticChecks(ctx context.Context, t testing.T, diagnostics []tfjson.Diagnostic, diagChecks []diagcheck.DiagnosticCheck) error {
	t.Helper()

	var result []error

	for _, diagCheck := range diagChecks {
		resp := diagcheck.CheckDiagnosticsResponse{}
		diagCheck.CheckDiagnostics(ctx, diagcheck.CheckDiagnosticsRequest{Diagnostics: diagnostics}, &resp)

		result = append(result, resp.Error)
	}

	return errors.Join(result...)
}
