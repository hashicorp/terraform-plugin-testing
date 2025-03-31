// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package teststep_test

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
)

func TestVerifyImportPlan(t *testing.T) {
	t.Parallel()

	resource := &tfjson.StateResource{
		Address: "example_resource.instance-1",
		Type:    "example_resource",
		Name:    "instance-1",
		AttributeValues: map[string]interface{}{
			"attr1": "value1",
		},
	}

	state := new(tfjson.State)
	state.Values = &tfjson.StateValues{
		RootModule: &tfjson.StateModule{
			Resources: []*tfjson.StateResource{resource},
		},
	}

	plan := new(tfjson.Plan)
	plan.ResourceChanges = []*tfjson.ResourceChange{
		{
			Address: "example_resource.instance-1",
			Change: &tfjson.Change{
				Actions: []tfjson.Action{tfjson.ActionNoop},
				After: map[string]interface{}{
					"attr1": "value1",
				},
				Before: map[string]interface{}{
					"attr1": "value1",
				},
				Importing: &tfjson.Importing{
					ID: "instance-1",
				},
			},
		},
	}

	if err := teststep.VerifyImportPlan(plan, state); err != nil {
		t.Fatal(err)
	}

}

func TestVerifyImportPlan_AttributeValueMismatch(t *testing.T) {
	t.Parallel()

	resource := &tfjson.StateResource{
		Address: "example_resource.instance-1",
		Type:    "example_resource",
		Name:    "instance-1",
		AttributeValues: map[string]interface{}{
			"attr1": "value1",
		},
	}

	state := new(tfjson.State)
	state.Values = &tfjson.StateValues{
		RootModule: &tfjson.StateModule{
			Resources: []*tfjson.StateResource{resource},
		},
	}

	plan := new(tfjson.Plan)
	plan.ResourceChanges = []*tfjson.ResourceChange{
		{
			Address: "example_resource.instance-1",
			Change: &tfjson.Change{
				Actions: []tfjson.Action{tfjson.ActionNoop},
				After: map[string]interface{}{
					"attr1": "value5",
				},
				Before: map[string]interface{}{
					"attr1": "value5",
				},
				Importing: &tfjson.Importing{
					ID: "instance-1",
				},
			},
		},
	}

	if err := teststep.VerifyImportPlan(plan, state); err == nil {
		t.Fatal("expected error, got nil")
	}

}
