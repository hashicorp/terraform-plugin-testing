// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep

import (
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// VerifyImportPlan compares a Terraform plan against a known good state
func VerifyImportPlan(plan *tfjson.Plan, state *terraform.State) error {
	if state == nil {
		return fmt.Errorf("state is nil")
	}
	if plan == nil {
		return fmt.Errorf("plan is nil")
	}
	oldResources := make(map[string]*terraform.ResourceState)
	for logicalResourceName, resourceState := range state.RootModule().Resources {
		if !strings.HasPrefix(logicalResourceName, "data.") {
			oldResources[logicalResourceName] = resourceState
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

			oldResource := oldResources[rc.Address]
			if oldResource == nil {
				// does this matter?
				return fmt.Errorf("importing resource %s: expected resource %s to exist in known state", rc.Address, rc.Change.Importing.ID)
			}

			attr, ok := oldResource.Primary.Attributes[k]
			if !ok {
				return fmt.Errorf("importing resource %s: expected %s in known state to exist", rc.Address, k)
			}
			if attr != vs {
				return fmt.Errorf("importing resource %s: expected %s in known state to be %q, got %q", rc.Address, k, oldResource.Primary.Attributes[k], vs)
			}
		}
	}
	return nil
}
