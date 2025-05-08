package resource

// Deprecated. This is an undocumented compatibility flag. When
// `RefreshAfterApply` is set to non-empty, a `Config`-mode test step will
// invoke a refresh before successful completion. This is intended as a
// compatibility measure for test cases that have different -- but
// semantically-equal -- state representations in their test steps. When
// comparing two states, the testing framework is not aware of semantic
// equality or set equality, as that would rely on provider logic.
var RefreshAfterApply string
