package plancheck

import (
	"context"
	"fmt"
)

var _ PlanCheck = expectResourceAction{}

type expectResourceAction struct {
	resourceAddress string
	actionType      ResourceActionType
}

func (e expectResourceAction) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	foundResource := false

	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress == rc.Address {
			switch e.actionType {
			case ResourceActionNoop:
				if !rc.Change.Actions.NoOp() {
					resp.Error = fmt.Errorf("'%s' - expected NoOp, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionCreate:
				if !rc.Change.Actions.Create() {
					resp.Error = fmt.Errorf("'%s' - expected Create, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionRead:
				if !rc.Change.Actions.Read() {
					resp.Error = fmt.Errorf("'%s' - expected Read, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionUpdate:
				if !rc.Change.Actions.Update() {
					resp.Error = fmt.Errorf("'%s' - expected Update, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionDestroy:
				if !rc.Change.Actions.Delete() {
					resp.Error = fmt.Errorf("'%s' - expected Destroy, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionDestroyBeforeCreate:
				if !rc.Change.Actions.DestroyBeforeCreate() {
					resp.Error = fmt.Errorf("'%s' - expected DestroyBeforeCreate, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionCreateBeforeDestroy:
				if !rc.Change.Actions.CreateBeforeDestroy() {
					resp.Error = fmt.Errorf("'%s' - expected CreateBeforeDestroy, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			case ResourceActionReplace:
				if !rc.Change.Actions.Replace() {
					resp.Error = fmt.Errorf("%s - expected Replace, got action(s): %v", rc.Address, rc.Change.Actions)
					return
				}
			default:
				resp.Error = fmt.Errorf("%s - unexpected ResourceActionType byte: %d", rc.Address, e.actionType)
				return
			}

			foundResource = true
			break
		}
	}

	if !foundResource {
		resp.Error = fmt.Errorf("%s - Resource not found in plan ResourceChanges", e.resourceAddress)
		return
	}
}

// TODO: document
func ExpectResourceAction(resourceAddress string, actionType ResourceActionType) PlanCheck {
	return expectResourceAction{
		resourceAddress: resourceAddress,
		actionType:      actionType,
	}
}
