package teststep_test

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func Test_VerifyImportPlan(t *testing.T) {
	t.Parallel()

	state := &terraform.State{
		Version: 6,
		Modules: []*terraform.ModuleState{
			{
				Path: []string{"root"},
				Resources: map[string]*terraform.ResourceState{
					"example_resource.instance-1": {
						Primary: &terraform.InstanceState{
							Attributes: map[string]string{
								"attr1": "value1",
							},
						},
					},
				},
			},
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
