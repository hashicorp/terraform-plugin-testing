// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func (c *Collector) OperationForAddress(address, operation string) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	address = safeAddress(address)
	for _, f := range c.features {
		if f.address == address && f.kind == KindResource {
			f.completed.add(operation)
			return
		}
	}
	// Import steps without a config still identify the resource by address.
	parts := strings.Split(address, ".")
	if len(parts) > 0 {
		f := c.feature(KindResource, parts[0], address)
		if operation == "import" {
			f.completed.add(operation)
		}
	}
}

func (c *Collector) ApplyCompleted(destroy bool) {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, p := range c.lastPlan {
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
	for _, p := range c.actions {
		f := c.feature(p.Kind, p.Name, p.Address)
		if f.provider == "" {
			f.provider = p.Provider
		}
		f.completed.add("invoke")
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

func (c *Collector) ProviderStepCompleted() {
	if c == nil {
		return
	}
	if !c.active {
		return
	}
	for _, f := range c.features {
		if f.kind == KindProvider {
			f.completed.add("step")
		}
	}
}
