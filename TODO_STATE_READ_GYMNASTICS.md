# TODO: simplify the state reads in the import and refresh steps

Source: review comments ("why these gymnastics?") in `testing_new_import_state.go` and `testing_new_refresh_state.go`.

## What

Both steps contain this code:

```go
var importState *terraform.State
err = runProviderCommand(ctx, t, workingDir, providers, func() error {
	jsonState, state, stateErr := getState(ctx, t, workingDir)
	importState, err = state, stateErr
	if stateErr == nil {
		workingDir.Coverage.ObserveState(jsonState)
	}
	if err != nil {
		return err
	}
	return nil
})
```

## Why it is a problem

- The closure assigns the **outer** `err`, and that same `err` is then overwritten by `runProviderCommand`'s return value. It works, but it's confusing, and it's easy to break.
- `stateErr` and `err` hold the same value, yet the code checks them separately.
- `if err != nil { return err }; return nil` is just `return err`.
- Before coverage was added, the code was `importState, err = getState(...)`. The new code bends the closure so that it can also keep `jsonState`.

## Fix

```go
var jsonState *tfjson.State
var importState *terraform.State
err = runProviderCommand(ctx, t, workingDir, providers, func() error {
	var err error
	jsonState, importState, err = getState(ctx, t, workingDir)
	return err
})
if err != nil {
	t.Fatalf("Error getting state: %s", err)
}
coverage.ObserveState(jsonState)
```

Observe the state after the closure. Only the provider command needs to run inside it. `t.Fatalf` already handles the error path, so `ObserveState` never sees a failed read. Do the same in `testing_new_refresh_state.go`.

## Notes

- If `PLAN_FIX_SHOTGUN_SURGERY.md` lands first, `coverage` is the `stepCoverage` parameter, not `workingDir.Coverage`.
- `getState` calls `t.Fatal` on a shim error. That behaviour is unchanged.
