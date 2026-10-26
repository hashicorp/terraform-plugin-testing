// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestStateRequiresExplicitPostOperationObservation(t *testing.T) {
	t.Parallel()

	c := &Collector{}
	c.BeginStep(1)
	c.Plan(&tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}, After: map[string]any{"id": "secret"}}}}})
	f := c.feature(KindResource, "example_item", "example_item.test")
	if len(f.known) != 0 {
		t.Fatalf("plan claimed post-operation state: %v", f.known)
	}
	c.ObserveState(&tfjson.State{Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{Resources: []*tfjson.StateResource{{Address: "example_item.test", Type: "example_item", Mode: tfjson.ManagedResourceMode, AttributeValues: map[string]any{"id": "secret"}}}}}})
	if !contains(f.known, "id") {
		t.Fatalf("missing post-operation state: %v", f.known)
	}
}
