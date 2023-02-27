package resource

// TODO: is this a good name?
// TODO: document type and all aliases below
type DiffChangeType byte

const (
	DiffInvalid DiffChangeType = iota
	DiffNoop
	DiffCreate
	DiffRead
	DiffUpdate
	DiffDestroy
	DiffDestroyBeforeCreate
	DiffCreateBeforeDestroy
	// TODO: document, this is more of a helper to detect either of the two above
	DiffReplace
)
