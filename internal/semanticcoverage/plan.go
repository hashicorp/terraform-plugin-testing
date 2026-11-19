// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// ObservePlan records the providers in plan and the operations it plans. A nil
// plan means the plan could not be read.
//
// A plan that could be read succeeded, so ObservePlan also records what
// planning did:
//   - It configured each provider in the plan's configuration.
//   - It imported each resource that an import block imports.
//   - It planned each action invocation in the plan.
//   - It opened and closed each ephemeral resource, and called each
//     provider-defined function, in the configuration that the step
//     observed. The plan doesn't list these, so they are inferred: call
//     ObservePlan only with plans of that configuration.
//
// Plan values are not recorded: populated attributes must come from state
// that an operation produced (see ObserveState).
func (s *Step) ObservePlan(plan *tfjson.Plan) {
	if s == nil {
		return
	}
	if plan == nil {
		s.warn(WarningPlanUnavailable)
		return
	}
	if plan.Config != nil {
		for _, provider := range plan.Config.ProviderConfigs {
			if provider != nil && provider.FullName != "" {
				appendNonEmpty(&s.feature(KindProvider, provider.FullName, "").Completed, OperationConfigure)
			}
		}
	}
	s.plannedOperations(plan, func(f *Feature, operation Operation) {
		appendNonEmpty(&f.Planned, operation)
		switch operation {
		case OperationImportID, OperationImportIdentity:
			appendNonEmpty(&f.Completed, operation)
		case OperationInvoke:
			appendNonEmpty(&f.Completed, OperationPlan)
		}
	})
	for _, f := range s.features {
		switch f.Kind {
		case KindEphemeralResource:
			appendNonEmpty(&f.Completed, OperationOpen)
			appendNonEmpty(&f.Completed, OperationClose)
		case KindFunction:
			appendNonEmpty(&f.Completed, OperationCall)
		}
	}
}

// ApplyCompleted records the operations of plan as completed, after it was
// applied successfully. A nil plan records nothing; ObservePlan warns about
// it.
func (s *Step) ApplyCompleted(plan *tfjson.Plan) {
	if s == nil || plan == nil {
		return
	}
	s.plannedOperations(plan, func(f *Feature, operation Operation) {
		appendNonEmpty(&f.Completed, operation)
	})
}

// plannedOperations calls fn for each resource change and action invocation
// in plan. A change that imports a resource has two operations: the import,
// and the change to the imported object. A change without an operation, such
// as a no-op, still declares its feature, with an empty operation.
func (s *Step) plannedOperations(plan *tfjson.Plan, fn func(f *Feature, operation Operation)) {
	for _, rc := range plan.ResourceChanges {
		if rc == nil || rc.Change == nil {
			continue
		}
		kind, ok := modeKinds[rc.Mode]
		if !ok {
			continue
		}
		f := s.feature(kind, rc.Type, rc.ProviderName)
		fn(f, operationForActions(rc.Change.Actions))
		if rc.Change.Importing != nil {
			fn(f, importOperation(rc.Change.Importing))
		}
	}
	for _, invocation := range plan.ActionInvocations {
		if invocation != nil {
			fn(s.feature(KindAction, invocation.Type, invocation.ProviderName), OperationInvoke)
		}
	}
}

// importOperation returns the operation of an import block: an import by
// identity if the block gives one, and by ID otherwise.
func importOperation(importing *tfjson.Importing) Operation {
	if importing.Identity != nil {
		return OperationImportIdentity
	}
	return OperationImportID
}

// modeKinds maps the resource modes in plans and state to feature kinds.
var modeKinds = map[tfjson.ResourceMode]Kind{
	tfjson.ManagedResourceMode: KindResource,
	tfjson.DataResourceMode:    KindDataSource,
}

// operationForActions returns the operation of a resource change, or "" if
// it changes nothing.
func operationForActions(actions tfjson.Actions) Operation {
	switch {
	case actions.Create():
		return OperationCreate
	case actions.Delete():
		return OperationDelete
	case actions.Update():
		return OperationUpdate
	case actions.Replace():
		return OperationReplace
	case actions.Read():
		return OperationRead
	}
	return ""
}
