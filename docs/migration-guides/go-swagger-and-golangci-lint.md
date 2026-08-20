# Migration guide: go-swagger 0.36 and golangci-lint 2.13

Two upgrades that need attention in downstream projects: the forked swagger
templates were rebased, and the linter jump surfaces findings in code that
previously passed.

## go-swagger v0.29 → v0.36.4

### The forked server templates were rebased

`api/templates/server/builder.gotmpl` and `parameter.gotmpl` are forks of upstream
go-swagger templates. They had not been refreshed since 2020 and no longer work
with the current generator:

- the `comment` and `blockcomment` template functions were removed in favour of
  `lineComment` / `blockComment`
- `blockComment` now emits the `/* */` wrapper itself, so the old templates
  produced **nested block comments**, i.e. invalid Go
- parameter documentation moved from block comments to line comments, and the
  templates now iterate `.ServerParams` rather than `.Params`

Both files were therefore re-forked from upstream v0.36.4 with only go-starter's
customizations re-applied, each fenced with `BEGIN/END CUSTOM BY AAA`:

1. the `query` / `param` / `form` struct tags — needed because we bind through
   `echo.DefaultBinder` instead of go-swagger's generated server
2. the standalone `Validate(strfmt.Registry) error` method — this is what makes the
   generated `Params` types satisfy `go-openapi/runtime.Validatable`, the interface
   every `internal/util.BindAndValidate*` helper takes
3. the supporting `customsliceparamvalidator` / nullable / recursive-array
   validation plumbing

**If you modified these two templates in your project, your changes are not
preserved by the merge.** Diff your pre-upgrade versions against the new ones,
re-apply your changes on top of the new upstream base, then:

```bash
make swagger
git diff internal/types/    # review carefully — this is your request validation
```

If your customizations use template functions, check them against the current
function map: the built-ins come from `sprig` plus go-swagger's own set (see
`generator/internal/funcmaps/golang/funcmap.go` upstream). `comment` and
`blockcomment` are the two that most commonly break.

### The go-openapi stack moved with it

`runtime` 0.28 → 0.33, `strfmt` 0.23 → 0.27, `validate` 0.24 → 0.26,
`errors` 0.22.2 → 0.22.8. Generated code now imports the `swag` submodules
directly (`typeutils`, `jsonutils`) — see
[the swag migration guide](swag-root-package-to-submodules.md) for your
hand-written code.

## golangci-lint v2.4 → v2.13

Run `make lint` after upgrading. Expect findings in code that previously passed —
these are new or stricter rules, not regressions you introduced.

### Config change: `gomodguard` was renamed

Deprecated since golangci-lint v2.12.0. In `.golangci.yml`:

```yaml
linters:
  enable:
    - gomodguard_v2   # was: gomodguard
```

### The default caps hide the real count

`golangci-lint` defaults to `max-issues-per-linter: 50` and `max-same-issues: 3`,
so the summary undercounts badly. Get the true picture before you start:

```bash
golangci-lint run --max-issues-per-linter 0 --max-same-issues 0
```

> In go-starter this was the difference between a reported 26 `SA1019` findings
> and the actual ~85, and between 3 reported `gosec` G101 findings and 17.

### What go-starter hit, and how it was resolved

| Finding | Resolution |
| --- | --- |
| `staticcheck SA1019` — deprecated `swag.*` | Migrated, see [the swag guide](swag-root-package-to-submodules.md) |
| `govet inline` — `reflect.Ptr` should be inlined | Replaced with `reflect.Pointer` (a stdlib alias since Go 1.18) |
| `noctx` — `httptest.NewRequest` | Replaced with `httptest.NewRequestWithContext`; use `t.Context()` where a `*testing.T` is available |
| `staticcheck QF1012` | `fmt.Fprintf(&builder, ...)` instead of `builder.WriteString(fmt.Sprintf(...))` |
| `prealloc` | Preallocate with a known capacity, or seed the slice from the first call instead of appending into an empty one |
| `goconst` (48 of 50 in `_test.go`) | Config: `goconst.ignore-tests: true` — repeated literals are table-driven test case names, and extracting them hurts readability |
| `gosec G101` in fixtures/tests | Config: an exclusion rule scoped to `(_test\.go\|^internal/test/fixtures/)` — our static UUID test tokens are not secrets |

The two config decisions are in `.golangci.yml` with comments explaining why. If
your project disagrees, they are localized and easy to drop.

## Verify

```bash
make build   # includes swagger regeneration and lint
make test
```

Note that `make test` was also fixed in this release to actually fail on failing
tests. If your test suite has been quietly green, re-check it now.
