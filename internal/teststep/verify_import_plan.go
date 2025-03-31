// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

// VerifyImportPlan compares a Terraform plan against a known good state
func VerifyImportPlan(plan *tfjson.Plan, state *tfjson.State) error {
	if state == nil {
		return fmt.Errorf("state is nil")
	}
	if plan == nil {
		return fmt.Errorf("plan is nil")
	}

	// TODO: if state.Values is nil ...

	resourcesInState := make(map[string]*tfjson.StateResource)
	for _, resource := range state.Values.RootModule.Resources {
		if !strings.HasPrefix(resource.Address, "data.") {
			resourcesInState[resource.Address] = resource
		}
	}
	for _, rc := range plan.ResourceChanges {
		if rc.Change == nil || rc.Change.Actions == nil {
			// does this matter?
			continue
		}

		if !rc.Change.Actions.NoOp() {
			return fmt.Errorf("importing resource %s: expected a no-op resource action, got %q action", rc.Address, rc.Change.Actions)
		}

		if rc.Change.Importing == nil {
			return fmt.Errorf("importing resource %s: expected importing to be true", rc.Address)
		}
	}

	for _, rc := range plan.ResourceChanges {
		after, ok := rc.Change.After.(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("unexpected type %T", rc.Change.After))
		}

		for k, v := range after {
			vs, ok := v.(string)
			if !ok {
				panic(fmt.Sprintf("unexpected type %T", v))
			}

			resourceInState := resourcesInState[rc.Address]
			if resourceInState == nil {
				// does this matter?
				return fmt.Errorf("importing resource %s: expected resource %s to exist in known state", rc.Address, rc.Change.Importing.ID)
			}

			attr, ok := resourceInState.AttributeValues[k]
			if !ok {
				return fmt.Errorf("importing resource %s: expected %s in known state to exist", rc.Address, k)
			}
			if attr != vs {
				return fmt.Errorf("importing resource %s: expected %s in known state to be %q, got %q", rc.Address, k, attr, vs)
			}
		}
	}
	return nil
}
