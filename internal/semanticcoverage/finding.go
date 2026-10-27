// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Package semanticcoverage records structural evidence from acceptance-test steps.
// It never serializes configuration or Terraform values; see redact.go.
package semanticcoverage

const EnvDirectory = "TF_ACC_SEMANTIC_COVERAGE_DIR"

// FindingVersion is the schema version of the JSON findings written per step.
//
// Findings live in a user-owned directory (EnvDirectory) that outlives the
// test run. They are written by the test binary and read by
// cmd/semantic-coverage-report, and the two can be built from different
// versions of this module. Bump FindingVersion whenever the meaning or shape
// of Finding changes incompatibly, so that the report rejects stale findings
// instead of misreading them.
const FindingVersion = 1

type Feature struct {
	Kind      Kind     `json:"kind"`
	Name      string   `json:"name"`
	Provider  string   `json:"provider,omitempty"`
	Address   string   `json:"address,omitempty"`
	Inputs    []string `json:"inputs,omitempty"`
	Known     []string `json:"known,omitempty"`
	Null      []string `json:"null,omitempty"`
	Planned   []string `json:"planned,omitempty"`
	Completed []string `json:"completed,omitempty"`
}

// Kind is the type of a provider feature, such as a resource or a data
// source.
type Kind string

const (
	KindProvider          Kind = "provider"
	KindResource          Kind = "resource"
	KindDataSource        Kind = "data_source"
	KindEphemeralResource Kind = "ephemeral_resource"
	KindListResource      Kind = "list_resource"
	KindAction            Kind = "action"
	KindStateStore        Kind = "state_store"
	KindFunction          Kind = "function"
	KindIdentity          Kind = "identity"
)

// Valid reports whether k is one of the known kinds.
func (k Kind) Valid() bool {
	switch k {
	case KindProvider, KindResource, KindDataSource, KindEphemeralResource, KindListResource, KindAction, KindStateStore, KindFunction, KindIdentity:
		return true
	}
	return false
}

// Outcome is the final status of a test step, as written to a Finding.
type Outcome string

const (
	// OutcomeSuccess means the step completed and passed its checks.
	OutcomeSuccess Outcome = "success"
	// OutcomeExpectedError means the step failed with an error matching
	// TestStep.ExpectError.
	OutcomeExpectedError Outcome = "expected_error"
	// OutcomeFailed means the test failed during the step.
	OutcomeFailed Outcome = "failed"
)

// Valid reports whether o is one of the known outcomes.
func (o Outcome) Valid() bool {
	switch o {
	case OutcomeSuccess, OutcomeExpectedError, OutcomeFailed:
		return true
	}
	return false
}

// Completed reports whether the step ran to completion, either passing or
// failing with the expected error.
func (o Outcome) Completed() bool {
	return o == OutcomeSuccess || o == OutcomeExpectedError
}

// Warning means that one type of evidence for a step may be incomplete. The
// codes don't say why; use the logs or the collector tests for that.
type Warning string

const (
	// WarningConfigUnavailable means some configuration could not be read or
	// parsed, so input and feature evidence may be missing.
	WarningConfigUnavailable Warning = "config_unavailable"
	// WarningPlanUnavailable means a plan could not be read, so planned and
	// completed operation evidence may be missing.
	WarningPlanUnavailable Warning = "plan_unavailable"
	// WarningStateUnavailable means the post-operation state could not be
	// read, so known and null attribute evidence may be missing.
	WarningStateUnavailable Warning = "state_unavailable"
)

// Valid reports whether w is one of the known warnings.
func (w Warning) Valid() bool {
	switch w {
	case WarningConfigUnavailable, WarningPlanUnavailable, WarningStateUnavailable:
		return true
	}
	return false
}

type Finding struct {
	Version int `json:"version"`
	// Package is the import path of the package under test.
	Package  string    `json:"package,omitempty"`
	Test     string    `json:"test"`
	Step     int       `json:"step"`
	Outcome  Outcome   `json:"outcome"`
	Features []Feature `json:"features,omitempty"`
	Warnings []Warning `json:"warnings,omitempty"`
}
