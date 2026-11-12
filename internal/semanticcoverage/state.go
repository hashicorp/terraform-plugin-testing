// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

// ObserveState records the features in state and the populated attributes of
// root module resources.
// Only pass state that one of the step's operations produced. A nil state
// means the state could not be read.
func (s *Step) ObserveState(state *tfjson.State) {
	if s == nil {
		return
	}
	if state == nil {
		s.warn(WarningStateUnavailable)
		return
	}
	stateResources(state, func(r *tfjson.StateResource) {
		kind, ok := modeKinds[r.Mode]
		if !ok {
			return
		}
		f := s.feature(kind, r.Type, r.ProviderName)
		// Module configuration is not read, so values that a module sets
		// would count as computed by the provider.
		if strings.HasPrefix(r.Address, "module.") {
			return
		}
		// Only top-level keys are schema names. Nested map keys can be
		// user data, so they are not read.
		for name, value := range r.AttributeValues {
			if value != nil {
				appendNonEmpty(&f.Populated, name)
			}
		}
	})
}

// stateResources calls fn for each resource in state, in all modules.
func stateResources(state *tfjson.State, fn func(r *tfjson.StateResource)) {
	if state == nil || state.Values == nil {
		return
	}
	var walk func(m *tfjson.StateModule)
	walk = func(m *tfjson.StateModule) {
		if m == nil {
			return
		}
		for _, r := range m.Resources {
			if r != nil {
				fn(r)
			}
		}
		for _, child := range m.ChildModules {
			walk(child)
		}
	}
	walk(state.Values.RootModule)
}
