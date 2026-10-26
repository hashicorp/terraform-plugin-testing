# TODO: small simplifications in `internal/semanticcoverage/collector.go`

Source: review comments. Do these after the file split in `PLAN_FIX_COLLECTOR_SPLIT.md`, so that each change lands in its new file.

## 1. Make `kind` a typed enum

The reviewer decided this ("let's make kind an enum"). It replaces the "Kind constants" follow-up in `PLAN_FIX_COLLECTOR_SPLIT.md`.

- Add `type Kind string`, with constants for `provider`, `resource`, `data_source`, `ephemeral_resource`, `list_resource`, `action`, `state_store`, `function` and `identity`, in `finding.go` next to `Outcome`. Give it the same `Valid()` method.
- Use `Kind` for `Feature.Kind` and for the `kind` parameter of `feature` and `setProvider`. The `report` command uses it for `target.Kind` and its group table.
- `readFindings` should reject findings with an invalid kind, the same way it does for `Outcome`.
- **Derive the address prefix and the config keyword from `Kind` with one table.** This option needs the least code. The kind values stay as they are: they match the provider schema JSON keys, so the report joins on them unchanged. Kind and prefix can't be identical anyway, because `resource`, `provider`, `state_store` and `function` have no prefix.

  ```go
  var kinds = map[string]struct{ kind Kind; prefix string }{ // keyed by config block keyword
      "resource":  {KindResource, ""},
      "data":      {KindDataSource, "data."},
      "ephemeral": {KindEphemeralResource, "ephemeral."},
      "list":      {KindListResource, "list."},
      "action":    {KindAction, "action."},
      "provider":  {KindProvider, ""},
  }
  ```

  Use it in `blocks` (the nested `switch`), `jsonBlocks` (the `if kind == ...` chain) and in the `Plan` / `module` mode→kind mapping. One lookup replaces each chain. Check that the special cases stay correct: list addresses are `list.<type>.<label>`, and the provider address is empty.

## 2. Replace `add` with a real set

`add` does a linear search, then appends. Callers never rely on insertion order: `Finish` sorts every slice before writing (`collector.go`, `sort.Strings(*xs)`). So this is a **set**, not an ordered set. It is used about 30 times, so keep a helper and don't inline it.

Options, in order of preference:

1. **An internal set type.** `type set map[string]struct{}`, with `add(v)` (ignores `""`) and `sorted() []string`. Use it in a collector-internal `featureEvidence` struct. Convert to the wire `Feature` (sorted `[]string`) in `Finish` only. That also removes the sort loop from `Finish`, and `Feature` stays a plain DTO.
2. **The smallest change.** Keep `[]string`, but write the body as `if v != "" && !slices.Contains(*xs, v) { *xs = append(*xs, v) }`, and rename the function to `appendUnique` (see `TODO_COLLECTOR_NAMING.md`).

Go has no standard set type, and adding a dependency for one is not worth it. Prefer option 1.

## 3. Inline `jsonInputs`

It is a 3-line loop with 3 callers in `jsonBlocks`. Inline it as `for name := range obj { s.add(name) }`. That becomes trivial once item 2 lands.

Keep `asMap`. Inlining its 7 call sites, some of them nested, would add temporaries and make the code harder to read.

## 4. Remove the magic query file name

`QueryConfig` hardcodes `"terraform_plugin_test.tfquery.hcl"` twice. The same constant already exists twice elsewhere:

- `plugintest.QueryFileName` (`internal/plugintest/working_dir.go:25`)
- `teststep.rawQueryConfigFileName` (`internal/teststep/config.go:21`, unexported)

`semanticcoverage` must not import `plugintest`, because that would reverse the dependency that `TODO_WORKINGDIR_COVERAGE_COUPLING.md` removes. Change `QueryConfig(dir string)` to `QueryConfig(path string)` instead. The caller then passes `filepath.Join(wd.baseDir, plugintest.QueryFileName)`, and `Parse` gets `filepath.Base(path)`. The collector then holds no file names.

Also look at the `".tfquery.hcl"` suffix checks in `Config` and `blocks`. They are extensions, not file names, so they can stay as they are. If they spread further, add an unexported `queryFileExt` constant.
