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
		if f.Address == address && f.Kind == "resource" {
			add(&f.Completed, operation)
			return
		}
	}
	// Import steps without a config still identify the resource by address.
	parts := strings.Split(address, ".")
	if len(parts) > 0 {
		f := c.feature("resource", parts[0], address)
		if operation == "import" {
			add(&f.Completed, operation)
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
		if f.Provider == "" {
			f.Provider = p.Provider
		}
		if !destroy {
			add(&f.Completed, "apply")
		}
		for _, operation := range p.Planned {
			if operation == "unknown" {
				continue
			}
			if destroy {
				if operation != "delete" {
					continue
				}
				add(&f.Completed, "destroy")
			}
			add(&f.Completed, operation)
		}
	}
	for _, p := range c.actions {
		f := c.feature(p.Kind, p.Name, p.Address)
		if f.Provider == "" {
			f.Provider = p.Provider
		}
		add(&f.Completed, "invoke")
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
		if f.Kind == "list_resource" {
			add(&f.Completed, "query")
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
				if f.Kind == "list_resource" && f.Name == m.ListResourceFound.ResourceType {
					add(&f.Completed, "result")
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
		if f.Kind == "resource" {
			add(&f.Completed, "refresh")
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
		if f.Kind == "state_store" {
			add(&f.Completed, "state_store")
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
		if f.Kind == "provider" {
			add(&f.Completed, "step")
		}
	}
}
