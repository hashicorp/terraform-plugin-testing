// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Package semanticcoverage records which provider features each
// acceptance-test step uses, and writes one Finding per step.
//
// # Privacy
//
// Findings are written to a user-owned directory that outlives the test run
// and may be shared, so the collector records structure and never values:
//
//   - It does not persist configuration text, attribute values, plan or state
//     JSON, query results, resource names, addresses, instance keys, or
//     assertion matchers.
//   - From state, it reads only top-level attribute names and whether each
//     value is non-null, as nested map keys can be user data.
//
// It does persist provider source addresses and feature type names, schema
// attribute names, operations, test names, and test package import paths.
package semanticcoverage

import (
	"errors"
	"fmt"
	"slices"
)

// FindingVersion is the schema version of the JSON findings written per step.
//
// Findings live in a user-owned directory that outlives the
// test run. They are written by the test binary and read by
// cmd/semantic-coverage-report, and the two can be built from different
// versions of this module. Bump FindingVersion whenever the meaning or shape
// of Finding changes incompatibly, so that the report rejects stale findings
// instead of misreading them.
const FindingVersion = 1

// Finding is the evidence written for one test step. Findings are written
// only for steps that passed without TestStep.ExpectError: a step that
// expects an error exercises the provider's error handling rather than the
// feature, and a failed step can stop at any point.
type Finding struct {
	Version int `json:"version"`
	// Package is the import path of the package under test.
	Package string `json:"package,omitempty"`
	Test    string `json:"test"`
	// Step is the 1-based number of the test step. The destroy that runs
	// after the last step has the number after it.
	Step int `json:"step"`
	// PostTestDestroy reports whether the finding is of the destroy that
	// runs after the last step, rather than of a test step.
	PostTestDestroy bool `json:"post_test_destroy,omitempty"`
	// ServedProviders lists the source addresses under which the test
	// harness served provider factories from the test binary during the step.
	// The harness chooses these addresses (registry.terraform.io/hashicorp/
	// followed by the factory name, unless TF_ACC_PROVIDER_HOST or
	// TF_ACC_PROVIDER_NAMESPACE are set), so they can differ from the address
	// a provider is published under. Providers that Terraform installs, such
	// as ExternalProviders, are not listed.
	ServedProviders []string  `json:"served_providers,omitempty"`
	Features        []Feature `json:"features,omitempty"`
	Warnings        []Warning `json:"warnings,omitempty"`
}

// Validate reports whether f was written by this version of the collector and
// uses only the defined kinds, operations, and warnings.
func (f Finding) Validate() error {
	switch {
	case f.Version != FindingVersion:
		return fmt.Errorf("finding version %d, this report reads version %d; re-run the tests with the same terraform-plugin-testing version", f.Version, FindingVersion)
	case f.Test == "" || f.Step < 1:
		return errors.New("missing test name or step")
	}
	for _, w := range f.Warnings {
		if !w.valid() {
			return fmt.Errorf("invalid warning %q", w)
		}
	}
	for _, feature := range f.Features {
		if !feature.Kind.valid() {
			return fmt.Errorf("invalid feature kind %q", feature.Kind)
		}
		for _, operation := range slices.Concat(feature.Planned, feature.Completed) {
			if !operation.valid() {
				return fmt.Errorf("invalid operation %q", operation)
			}
		}
	}
	return nil
}

// Feature is the evidence for one provider schema entry in a step. All uses of
// the same entry in the step, such as several resources of one type, are
// merged into one Feature.
type Feature struct {
	Kind Kind `json:"kind"`
	// Name is the type name of the feature. For a provider, it is the local
	// name from a provider block or the source address from a plan.
	Name string `json:"name"`
	// Provider is the provider's source address (from plans and state), its
	// local name (for functions), or empty (from configuration).
	Provider string `json:"provider,omitempty"`
	// Configured lists the top-level attributes and blocks that the step's
	// configuration sets.
	Configured []string `json:"configured,omitempty"`
	// Populated lists the top-level attributes that are non-null in the state
	// that the step's operations produced.
	Populated []string `json:"populated,omitempty"`
	// Planned lists the operations in the step's plans.
	Planned []Operation `json:"planned,omitempty"`
	// Completed lists the operations that the step performed successfully.
	Completed []Operation `json:"completed,omitempty"`
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
)

func (k Kind) valid() bool {
	switch k {
	case KindProvider, KindResource, KindDataSource, KindEphemeralResource, KindListResource, KindAction, KindStateStore, KindFunction:
		return true
	}
	return false
}

// Operation is a lifecycle event of a feature, named after the provider
// protocol RPC that implements it. Some operations are inferred rather than
// observed; see the Step methods that record them.
type Operation string

const (
	// Resource operations. A replace runs both Create and Delete.
	OperationCreate  Operation = "create"
	OperationUpdate  Operation = "update"
	OperationReplace Operation = "replace"
	OperationDelete  Operation = "delete"
	// OperationRead is a read of a resource, a data source, or the state
	// of a state store.
	OperationRead Operation = "read"
	// OperationImportID is an import by ID, with the import command or an
	// import block with an id.
	OperationImportID Operation = "import_id"
	// OperationImportIdentity is an import with an import block that gives
	// the resource identity.
	OperationImportIdentity Operation = "import_identity"

	// OperationOpen and OperationClose are the lifecycle of an ephemeral
	// resource.
	OperationOpen  Operation = "open"
	OperationClose Operation = "close"

	// OperationQuery lists the resources of a list resource.
	OperationQuery Operation = "query"

	// OperationPlan and OperationInvoke are the lifecycle of an action.
	OperationPlan   Operation = "plan"
	OperationInvoke Operation = "invoke"

	// OperationCall is a call of a provider-defined function.
	OperationCall Operation = "call"

	// OperationConfigure configures a provider or a state store.
	OperationConfigure Operation = "configure"

	// State store operations, besides configure and read.
	OperationWrite       Operation = "write"
	OperationListStates  Operation = "list_states"
	OperationDeleteState Operation = "delete_state"
	OperationLock        Operation = "lock"
	OperationUnlock      Operation = "unlock"
)

func (o Operation) valid() bool {
	switch o {
	case OperationCreate, OperationUpdate, OperationReplace, OperationDelete, OperationRead,
		OperationImportID, OperationImportIdentity,
		OperationOpen, OperationClose,
		OperationQuery,
		OperationPlan, OperationInvoke,
		OperationCall,
		OperationConfigure,
		OperationWrite, OperationListStates, OperationDeleteState, OperationLock, OperationUnlock:
		return true
	}
	return false
}

// Warning means that one source of evidence for a step was unavailable. The
// codes don't record the cause, and the collector doesn't log it. Reproduce
// the step's evidence in the collector tests to find out why.
type Warning string

const (
	// WarningConfigUnavailable means some configuration could not be read or
	// parsed, so features and configured attributes may be missing.
	WarningConfigUnavailable Warning = "config_unavailable"
	// WarningPlanUnavailable means a plan could not be read, so planned and
	// completed operations may be missing.
	WarningPlanUnavailable Warning = "plan_unavailable"
	// WarningStateUnavailable means the state after an operation could not be
	// read, so populated attributes may be missing.
	WarningStateUnavailable Warning = "state_unavailable"
)

func (w Warning) valid() bool {
	switch w {
	case WarningConfigUnavailable, WarningPlanUnavailable, WarningStateUnavailable:
		return true
	}
	return false
}
