// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// Operations that don't come from a plan are recorded by the test step mode
// that performs them, right after Terraform reports success. A finding is
// written only if the step as a whole passed (see Step.Finish).

// ImportCompleted records that the import command imported the resource at
// address. state is the state after the import, which gives the resource's
// type.
func (s *Step) ImportCompleted(state *tfjson.State, address string) {
// <review>
// question: this only collects whether any kind of importing happened, could we start differentiating between import by id and by resource identity (if supported)?
// </review>
	if s == nil {
		return
	}
	stateResources(state, func(r *tfjson.StateResource) {
		if r.Address == address && r.Mode == tfjson.ManagedResourceMode {
			add(&s.feature(KindResource, r.Type, r.ProviderName).Completed, OperationImport)
		}
	})
}

// RefreshCompleted records that a refresh completed for every resource in the
// step.
func (s *Step) RefreshCompleted() {
	s.completedForKind(KindResource, OperationRefresh)
}

// QueryCompleted records that a query completed for every list resource in
// the step.
func (s *Step) QueryCompleted() {
	s.completedForKind(KindListResource, OperationQuery)
}

func (s *Step) completedForKind(kind Kind, operation Operation) {
	if s == nil {
		return
	}
	for _, f := range s.features {
		if f.Kind == kind {
			add(&f.Completed, operation)
		}
	}
}
