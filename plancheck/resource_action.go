package plancheck

// ResourceActionType is a string enum type that routes to a specific terraform-json.Actions function for asserting resource changes.
// https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions
type ResourceActionType string

const (
	// ResourceActionNoop is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.NoOp
	ResourceActionNoop ResourceActionType = "NoOp"

	// ResourceActionCreate is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Create
	ResourceActionCreate ResourceActionType = "Create"

	// ResourceActionRead is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Read
	ResourceActionRead ResourceActionType = "Read"

	// ResourceActionUpdate is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Update
	ResourceActionUpdate ResourceActionType = "Update"

	// ResourceActionDestroy is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Delete
	ResourceActionDestroy ResourceActionType = "Destroy"

	// ResourceActionDestroyBeforeCreate is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.DestroyBeforeCreate
	ResourceActionDestroyBeforeCreate ResourceActionType = "DestroyBeforeCreate"

	// ResourceActionCreateBeforeDestroy is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.CreateBeforeDestroy
	ResourceActionCreateBeforeDestroy ResourceActionType = "CreateBeforeDestroy"

	// ResourceActionReplace is a string enum that represents: https://pkg.go.dev/github.com/hashicorp/terraform-json#Actions.Replace
	ResourceActionReplace ResourceActionType = "Replace"
)
