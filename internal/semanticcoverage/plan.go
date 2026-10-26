// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// Plan records the planned operation of each resource change and action
// invocation in plan. A nil plan means the saved plan could not be read.
func (c *Collector) Plan(plan *tfjson.Plan) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	if plan == nil {
		c.warnings.add(WarningPlanUnavailable)
		return
	}
	if plan.Config != nil {
		for _, provider := range plan.Config.ProviderConfigs {
			if provider != nil && provider.FullName != "" {
				c.feature(KindProvider, provider.FullName, "")
			}
		}
	}
	// Several plans run during one step. Only the pre-apply plan can justify
	// completed actions; follow-up no-op plans must not replace that evidence.
	for _, rc := range plan.ResourceChanges {
		if rc == nil || rc.Change == nil {
			continue
		}
		kind, ok := modeKinds[rc.Mode]
		if !ok {
			continue
		}
		f := c.feature(kind, rc.Type, rc.Address)
		c.setProvider(kind, rc.Type, rc.Address, rc.ProviderName)
		f.planned.add(action(rc.Change.Actions))
	}
	for _, invocation := range plan.ActionInvocations {
		if invocation == nil {
			continue
		}
		f := c.feature(KindAction, invocation.Type, invocation.Address)
		c.setProvider(KindAction, invocation.Type, invocation.Address, invocation.ProviderName)
		f.planned.add("invoke")
	}
	// PlannedValues is a plan, not the state produced by a completed apply.
	// Computed coverage is taken from post-apply or refresh state instead.
}

// PreApplyPlan selects plan as the one that the next apply executes.
// ApplyCompleted marks its changes completed. A nil plan means the saved plan
// could not be read.
func (c *Collector) PreApplyPlan(plan *tfjson.Plan) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	c.lastPlan = nil
	c.actions = nil
	if plan == nil {
		c.warnings.add(WarningPlanUnavailable)
		return
	}
	for _, rc := range plan.ResourceChanges {
		if rc != nil && rc.Change != nil && rc.Mode == tfjson.ManagedResourceMode {
			c.lastPlan = append(c.lastPlan, Feature{Kind: KindResource, Name: rc.Type, Provider: rc.ProviderName, Address: safeAddress(rc.Address), Planned: []string{action(rc.Change.Actions)}})
		}
	}
	for _, invocation := range plan.ActionInvocations {
		if invocation != nil {
			c.actions = append(c.actions, Feature{Kind: KindAction, Name: invocation.Type, Provider: invocation.ProviderName, Address: safeAddress(invocation.Address)})
		}
	}
}

func (c *Collector) ImportPlan(plan *tfjson.Plan) {
	if c == nil || plan == nil {
		return
	}
	if !c.active {
		return
	}
	for _, rc := range plan.ResourceChanges {
		if rc != nil && rc.Change != nil && rc.Change.Importing != nil {
			f := c.feature(KindResource, rc.Type, rc.Address)
			if f.provider == "" {
				f.provider = rc.ProviderName
			}
			f.planned.add("import")
		}
	}
}

// modeKinds maps the resource modes in plans and state to feature kinds.
var modeKinds = map[tfjson.ResourceMode]Kind{
	tfjson.ManagedResourceMode: KindResource,
	tfjson.DataResourceMode:    KindDataSource,
}

func action(actions tfjson.Actions) string {
	if actions.Create() {
		return "create"
	}
	if actions.Delete() {
		return "delete"
	}
	if actions.Update() {
		return "update"
	}
	if actions.Replace() {
		return "replace"
	}
	if actions.NoOp() {
		return "no_op"
	}
	if actions.Read() {
		return "read"
	}
	return "unknown"
}
