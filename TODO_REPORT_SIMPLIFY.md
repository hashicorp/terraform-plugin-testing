# TODO: simplify `cmd/semantic-coverage-report`

Source: review comments in `cmd/semantic-coverage-report/main.go`.

## 1. Use `tfjson.ProviderSchemas` instead of the local `schemaDoc` types

The review asked whether a schema type already exists. It does: `github.com/hashicorp/terraform-json`, which this module already depends on (v0.28.0, `go.mod:15`) and which terraform-ls also uses. It has `tfjson.ProviderSchemas` / `tfjson.ProviderSchema`, and every group the report reads is covered:

| local field | tfjson field |
|---|---|
| `Provider` | `ConfigSchema *Schema` |
| `Resources` | `ResourceSchemas` |
| `DataSources` | `DataSourceSchemas` |
| `Ephemeral` | `EphemeralResourceSchemas` |
| `Lists` | `ListResourceSchemas` |
| `Actions` | `ActionSchemas map[string]*ActionSchema` (a different type: it wraps `Schema`) |
| `Stores` | `StateStoreSchemas` |
| `Identities` | `ResourceIdentitySchemas map[string]*IdentitySchema` (has `Attributes` and no `Block`) |
| `Functions` | `Functions map[string]*FunctionSignature` |

What changes:

- Delete `schemaDoc`, `providerSchema`, `schema`, `block` and `attribute`. Use the tfjson types instead. `flatten` and `addAttribute` then take `*tfjson.SchemaBlock` / `*tfjson.SchemaAttribute` (nested blocks: `NestedBlocks map[string]*SchemaBlockType`; nested attributes: `AttributeNestedType`). Check that `flatten` keeps its current paths.
- Actions and identities have their own schema types, so they don't fit the uniform group table. This happens to fit item 3, which duplicates the per-kind code anyway.
- Replace the manual `strings.HasPrefix(doc.FormatVersion, "1.")` check with `doc.Validate()`. It checks the major format version and returns a better error.
- The report's tests need fixtures that unmarshal into tfjson. They should, because the JSON shape is unchanged.

## 2. Inline `readFindings` into `main`

It has a single caller. Move the `filepath.WalkDir` closure into `main` directly. Keep the validation checks: test name and step, outcome, and kind once `TODO_COLLECTOR_SIMPLIFY.md` item 1 lands. Keep the finding version check as well (see `TODO_COMMENTS.md`). The new warning codes are validated too (`TODO_WARNING_CODES.md`).

## 3. Inline the `register` closure in `aggregate`, and duplicate it per kind

The reviewer decided to inline and duplicate it. Replace the closure plus the anonymous `[]struct{kind; members}` table with one plain loop per schema group (provider, resource, data source, ephemeral, list, action, state store, identity, function). Each loop builds its `target` directly.

The duplication is deliberate. The groups already differ in type (see item 1), and explicit loops are easier to read than a closure that is fed from a heterogeneous table. Factor out only the 3-line `target` construction, if anything.

## 4. Rename the flags to follow Terraform CLI conventions

In Terraform, `-json` is always a boolean ("machine-readable output"). Output paths use `-out` / `-state-out`, and input files use `-var-file`.

| current | new | default |
|---|---|---|
| `-findings` | `-findings-dir` | `findings` |
| `-schema` | `-schema-file` | (required) |
| `-json` | `-json-out` | `coverage.json` |
| `-text` | `-text-out` | `coverage.txt` |

- Both reports are still written on every run, with the defaults unchanged.
- No aliases for the old names. The command is experimental and internal to this repo.
- Update the `-schema is required` error message, the usage example in `SEMANTIC_COVERAGE.md` (around lines 35-37), and any scripts or tests that call the command.

## Order

Do 1 and 3 together, since both touch the same lines. Do 2 and 4 separately.
