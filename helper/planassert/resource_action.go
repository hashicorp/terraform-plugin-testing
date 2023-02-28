package planassert

// TODO: is this a good name?
// TODO: document type and all aliases below
type ResourceActionType byte

const (
	ResourceActionInvalid ResourceActionType = iota
	ResourceActionNoop
	ResourceActionCreate
	ResourceActionRead
	ResourceActionUpdate
	ResourceActionDestroy
	ResourceActionDestroyBeforeCreate
	ResourceActionCreateBeforeDestroy
	// TODO: document, this is more of a helper to detect either of the two above
	ResourceActionReplace
)
