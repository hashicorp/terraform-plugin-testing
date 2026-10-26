// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

// ImportCompleted marks the resource at address as imported by a successful
// import command.
func (c *Collector) ImportCompleted(address string) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	address = safeAddress(address)
	for _, f := range c.features {
		if f.address == address && f.kind == KindResource {
			f.completed.add("import")
			return
		}
	}
	// Import steps without a config still identify the resource by address.
	resourceType, _, _ := strings.Cut(address, ".")
	c.feature(KindResource, resourceType, address).completed.add("import")
}

// ApplyCompleted records a successful apply of the plan selected by
// PreApplyPlan: each planned change is marked "apply" plus its operation, and
// each action invocation "invoke". After a destroy apply, only deletions are
// marked, as "destroy" and "delete". Every provider in the step is marked
// "step".
func (c *Collector) ApplyCompleted(destroy bool) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, p := range c.pendingApplyChanges {
		f := c.feature(p.Kind, p.Name, p.Address)
		if f.provider == "" {
			f.provider = p.Provider
		}
		if !destroy {
			f.completed.add("apply")
		}
		for _, operation := range p.Planned {
			if operation == "unknown" {
				continue
			}
			if destroy {
				if operation != "delete" {
					continue
				}
				f.completed.add("destroy")
			}
			f.completed.add(operation)
		}
	}
	for _, p := range c.pendingApplyActions {
		f := c.feature(p.Kind, p.Name, p.Address)
		if f.provider == "" {
			f.provider = p.Provider
		}
		f.completed.add("invoke")
	}
	for _, f := range c.features {
		if f.kind == KindProvider {
			f.completed.add("step")
		}
	}
}

// QueryCompleted marks list resources in a query step that passed its checks.
func (c *Collector) QueryCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.kind == KindListResource {
			f.completed.add("query")
		}
	}
}

// QueryResults reads query message types without persisting result values,
// display names, instance keys, or provider-generated configuration.
func (c *Collector) QueryResults(messages []tfjson.LogMsg) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, message := range messages {
		switch m := message.(type) {
		case tfjson.ListResourceFoundMessage:
			for _, f := range c.features {
				if f.kind == KindListResource && f.name == m.ListResourceFound.ResourceType {
					f.completed.add("result")
				}
			}
		}
	}
}

// RefreshCompleted marks every resource in a refresh step that passed its
// checks as refreshed.
func (c *Collector) RefreshCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.kind == KindResource {
			f.completed.add("refresh")
		}
	}
}

// StateStoreCompleted marks every state store in a state store step that
// passed its checks as used.
func (c *Collector) StateStoreCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.kind == KindStateStore {
			f.completed.add("state_store")
		}
	}
}
