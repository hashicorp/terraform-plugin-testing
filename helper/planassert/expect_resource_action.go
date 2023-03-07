package planassert

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ resource.PlanAssert = expectResourceAction{}

type expectResourceAction struct {
	resourceAddress string
	actionType      ResourceActionType
}

func (e expectResourceAction) RunAssert(plan *tfjson.Plan) error {
	foundResource := false

	for _, rc := range plan.ResourceChanges {
		if e.resourceAddress == rc.Address {
			switch e.actionType {
			case ResourceActionNoop:
				if !rc.Change.Actions.NoOp() {
					return fmt.Errorf("'%s' - expected NoOp, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionCreate:
				if !rc.Change.Actions.Create() {
					return fmt.Errorf("'%s' - expected Create, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionRead:
				if !rc.Change.Actions.Read() {
					return fmt.Errorf("'%s' - expected Read, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionUpdate:
				if !rc.Change.Actions.Update() {
					return fmt.Errorf("'%s' - expected Update, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionDestroy:
				if !rc.Change.Actions.Delete() {
					return fmt.Errorf("'%s' - expected Destroy, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionDestroyBeforeCreate:
				if !rc.Change.Actions.DestroyBeforeCreate() {
					return fmt.Errorf("'%s' - expected DestroyBeforeCreate, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionCreateBeforeDestroy:
				if !rc.Change.Actions.CreateBeforeDestroy() {
					return fmt.Errorf("'%s' - expected CreateBeforeDestroy, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			case ResourceActionReplace:
				if !rc.Change.Actions.Replace() {
					return fmt.Errorf("%s - expected Replace, got action(s): %v", rc.Address, rc.Change.Actions)
				}
			default:
				return fmt.Errorf("%s - unexpected ResourceActionType byte: %d", rc.Address, e.actionType)
			}

			foundResource = true
			break
		}
	}

	if !foundResource {
		return fmt.Errorf("%s - Resource not found in plan ResourceChanges", e.resourceAddress)
	}

	return nil
}

// TODO: document
func ExpectResourceAction(resourceAddress string, actionType ResourceActionType) resource.PlanAssert {
	return expectResourceAction{
		resourceAddress: resourceAddress,
		actionType:      actionType,
	}
}
