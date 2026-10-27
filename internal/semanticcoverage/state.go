// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// ObserveState records only schema attribute names and null/known status from
// state explicitly obtained after a step operation, not setup or cleanup reads.
func (c *Collector) ObserveState(state *tfjson.State) {
	if c == nil || !c.active || state == nil || state.Values == nil || state.Values.RootModule == nil {
		return
	}
	c.observeStateModule(state.Values.RootModule)
}

func (c *Collector) observeStateModule(m *tfjson.StateModule) {
	for _, r := range m.Resources {
		if r == nil {
			continue
		}
		kind, ok := modeKinds[r.Mode]
		if !ok {
			continue
		}
		f := c.feature(kind, r.Type, r.Address)
		c.setProvider(kind, r.Type, r.Address, r.ProviderName)
		recordAttributeStatus(f, r.AttributeValues)
		if r.IdentitySchemaVersion != nil && len(r.IdentityValues) > 0 {
			identity := c.feature(KindIdentity, r.Type, r.Address)
			identity.provider = r.ProviderName
			recordAttributeStatus(identity, r.IdentityValues)
		}
	}
	for _, child := range m.ChildModules {
		c.observeStateModule(child)
	}
}

func recordAttributeStatus(f *featureEvidence, values map[string]any) {
	// Only top-level keys are unambiguously schema names. Nested map keys
	// could be sensitive user-provided values, so do not persist them.
	for name, value := range values {
		if value == nil {
			f.null.add(name)
		} else {
			f.known.add(name)
		}
	}
}
