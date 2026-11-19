// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// Operations that don't come from a plan are recorded by the test step mode
// that performs them, right after Terraform reports success. A finding is
// written only if the step as a whole passed (see Step.End).

// ImportCompleted records that the import command imported the resource at
// address. state is the state after the import, which gives the resource's
// type. The import command imports by ID; import blocks are recorded from
// the plan (see ObservePlan).
func (s *Step) ImportCompleted(state *tfjson.State, address string) {
	if s == nil {
		return
	}
	stateResources(state, func(r *tfjson.StateResource) {
		if r.Address == address && r.Mode == tfjson.ManagedResourceMode {
			appendNonEmpty(&s.feature(KindResource, r.Type, r.ProviderName).Completed, OperationImportID)
		}
	})
}

// RefreshCompleted records a read of each resource and data source in state,
// the state that a refresh produced. Pass the state after a refresh, or the
// prior state of a plan that refreshed. A nil state records nothing.
func (s *Step) RefreshCompleted(state *tfjson.State) {
	if s == nil {
		return
	}
	stateResources(state, func(r *tfjson.StateResource) {
		if kind, ok := modeKinds[r.Mode]; ok {
			appendNonEmpty(&s.feature(kind, r.Type, r.ProviderName).Completed, OperationRead)
		}
	})
}

// DestroyCompleted records a delete of each resource in state, the state
// before a destroy that succeeded. A nil state records nothing.
func (s *Step) DestroyCompleted(state *tfjson.State) {
	if s == nil {
		return
	}
	stateResources(state, func(r *tfjson.StateResource) {
		if r.Mode == tfjson.ManagedResourceMode {
			appendNonEmpty(&s.feature(KindResource, r.Type, r.ProviderName).Completed, OperationDelete)
		}
	})
}

// QueryCompleted records that a query completed for every list resource in
// the step.
func (s *Step) QueryCompleted() {
	s.completedForKind(KindListResource, OperationQuery)
}

// StateStoreCompleted records the operations of a StateStore step that
// passed for every state store in the step. The step initializes the store,
// lists its states, and writes and reads state. With workspaces, it also
// creates and deletes workspaces, which deletes their states.
func (s *Step) StateStoreCompleted(workspaces bool) {
	s.completedForKind(KindStateStore, OperationConfigure, OperationListStates, OperationWrite, OperationRead)
	if workspaces {
		s.completedForKind(KindStateStore, OperationDeleteState)
	}
}

// StateStoreLockVerified records that a client locked and unlocked the state
// of every state store in the step, while another client could not lock it.
func (s *Step) StateStoreLockVerified() {
	s.completedForKind(KindStateStore, OperationLock, OperationUnlock)
}

func (s *Step) completedForKind(kind Kind, operations ...Operation) {
	if s == nil {
		return
	}
	for _, f := range s.features {
		if f.Kind == kind {
			for _, operation := range operations {
				appendNonEmpty(&f.Completed, operation)
			}
		}
	}
}
