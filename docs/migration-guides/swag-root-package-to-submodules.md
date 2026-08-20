# Migration guide: the `swag` root package → submodules

`github.com/go-openapi/swag` v0.29 split its root package into per-concern Go
submodules (published from the same repository, versioned in lockstep) and
**deprecated all 124 exported symbols at the root**.

There is **no announced removal date** — upstream froze the root package for
backwards compatibility and it still forwards to the submodules — so nothing is
broken today. What you get by not migrating is a permanent stream of
`staticcheck SA1019` warnings, and drift from go-swagger's generated code, which
already emits the new import paths (e.g. `typeutils.IsZero` in `internal/types/`).

go-starter itself has completed this migration; this guide is for the
application code in your project.

## 1. Find your call sites

```bash
grep -rn '"github.com/go-openapi/swag"' --include='*.go' . | grep -v '/internal/types/'
```

`internal/types/` is generated — `make swagger` regenerates it with the new imports,
so never hand-edit those files.

Or let the linter enumerate them. Note the default caps hide the true count:

```bash
golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 | grep SA1019
```

> In go-starter this was the difference between "26 findings" and the real ~85.

## 2. Watch out for the import collision

If you also import `github.com/go-openapi/strfmt/conv`, it collides with the new
`github.com/go-openapi/swag/conv`.

**Recommended: drop `strfmt/conv` entirely.** Its helpers are type-specific pointer
constructors (`conv.UUID4(v) *strfmt.UUID4`, `conv.DateTime(v) *strfmt.DateTime`, …),
and the generic `conv.Pointer[T any](v T) *T` from `swag/conv` subsumes every one
of them:

```go
conv.UUID4(strfmt.UUID4("..."))     // before (strfmt/conv)
conv.Pointer(strfmt.UUID4("..."))   // after  (swag/conv) — same type, same result
```

That leaves a single unaliased `conv` import. Aliasing one of the two also works
(`swagconv "github.com/go-openapi/swag/conv"`), but you then carry the alias at
every new call site.

## 3. Rewrite the calls

Most of the mapping is mechanical and `gofmt -r` can do it safely (it rewrites
expressions, not text, so it will not touch strings or comments):

```bash
# the common pointer/value helpers
for fn in String Int Int32 Int64 Bool Float32 Float64 Time Uint Uint16 Uint32 Uint64; do
  gofmt -r "swag.${fn}(a) -> conv.Pointer(a)" -w .
  gofmt -r "swag.${fn}Value(a) -> conv.Value(a)" -w .
done

# strfmt/conv, if you are dropping it (step 2)
for fn in UUID4 UUID3 UUID DateTime Date Base64 Email Hostname Duration; do
  gofmt -r "conv.${fn}(a) -> conv.Pointer(a)" -w .
done

gofmt -r 'swag.ReadJSON(a, b) -> jsonutils.ReadJSON(a, b)' -w .
gofmt -r 'swag.WriteJSON(a) -> jsonutils.WriteJSON(a)' -w .
gofmt -r 'swag.IsZero(a) -> typeutils.IsZero(a)' -w .
```

Two caveats on the snippet above:

- The patterns match on the **local** import name. If you alias your imports
  (`sconv "github.com/go-openapi/strfmt/conv"`), write the pattern with that alias
  (`sconv.UUID4(a) -> conv.Pointer(a)`).
- `gofmt -r` rewrites call sites only — it does **not** touch the import block, so
  the code will not compile until you do step 3b.

**3b.** Fix the imports: drop `github.com/go-openapi/swag`, add the submodules you
now use, then run `make build`. **The compiler finds every remaining site** — this
migration is type-checked end to end, so there is no silent-breakage risk.

### The one real gotcha: untyped numeric constants

The old helpers were type-specific; `conv.Pointer` is generic and infers from its
argument. An untyped constant therefore changes type:

```go
swag.Int64(3600)            // *int64
conv.Pointer(3600)          // *int   — WRONG, will not compile against an *int64 field
conv.Pointer(int64(3600))   // *int64 — correct
```

The same applies to `swag.Int32`, `swag.Float32`, `swag.Uint*`. Calls that already
pass a typed expression (`swag.Int64(u.Unix())`, `swag.Int64(int64(code))`) convert
cleanly. The compiler flags every case; there is nothing to hunt for by hand.

### Not drop-in replacements

Four symbols changed shape, not just location — review these by hand:

| Old | New | What changed |
| --- | --- | --- |
| `swag.AddInitialisms(words...)` | `mangling.WithAdditionalInitialisms(words...)` | Was a mutator on a package global; is now an **option** passed to a mangler constructor. |
| `swag.GoNamePrefixFunc` | `mangling.WithGoNamePrefixFunc` | Same — global variable becomes a constructor option. |
| `swag.LoadFromFileOrHTTPWithTimeout(p, t)` | `loading.LoadFileOrHTTP(p, loading.WithTimeout(t))` | Renamed; timeout became a functional option. |
| `swag.DynamicJSONToStruct(d, tgt)` | `jsonutils.FromDynamicJSON(d, tgt)` | Renamed only; signature is identical. |

Also note `swag.ContainsStrings(coll, item)` — upstream points you at the standard
library's `slices.Contains(coll, item)` (identical argument order), though
`stringutils.ContainsStrings` also exists.

## 4. Verify

```bash
make build   # compiles, formats, and lints
make test
```

`make build` runs `golangci-lint`, so a clean run means no `SA1019` deprecation
warnings are left.

---

## Appendix: complete symbol mapping

All 124 deprecated root symbols, grouped by their new home. Generated from the
compatibility shims in `github.com/go-openapi/swag@v0.29.0` (`*_iface.go`).

### `conv` — `github.com/go-openapi/swag/conv`

95 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.Bool` | `conv.Pointer` |
| `swag.BoolMap` | `conv.PointerMap` |
| `swag.BoolSlice` | `conv.PointerSlice` |
| `swag.BoolValue` | `conv.Value` |
| `swag.BoolValueMap` | `conv.ValueMap` |
| `swag.BoolValueSlice` | `conv.ValueSlice` |
| `swag.ConvertBool` | `conv.ConvertBool` |
| `swag.ConvertFloat32` | `conv.ConvertFloat32` |
| `swag.ConvertFloat64` | `conv.ConvertFloat64` |
| `swag.ConvertInt16` | `conv.ConvertInt16` |
| `swag.ConvertInt32` | `conv.ConvertInt32` |
| `swag.ConvertInt64` | `conv.ConvertInt64` |
| `swag.ConvertInt8` | `conv.ConvertInt8` |
| `swag.ConvertUint16` | `conv.ConvertUint16` |
| `swag.ConvertUint32` | `conv.ConvertUint32` |
| `swag.ConvertUint64` | `conv.ConvertUint64` |
| `swag.ConvertUint8` | `conv.ConvertUint8` |
| `swag.Float32` | `conv.Pointer` |
| `swag.Float32Map` | `conv.PointerMap` |
| `swag.Float32Slice` | `conv.PointerSlice` |
| `swag.Float32Value` | `conv.Value` |
| `swag.Float32ValueMap` | `conv.ValueMap` |
| `swag.Float32ValueSlice` | `conv.ValueSlice` |
| `swag.Float64` | `conv.Pointer` |
| `swag.Float64Map` | `conv.PointerMap` |
| `swag.Float64Slice` | `conv.PointerSlice` |
| `swag.Float64Value` | `conv.Value` |
| `swag.Float64ValueMap` | `conv.ValueMap` |
| `swag.Float64ValueSlice` | `conv.ValueSlice` |
| `swag.FormatBool` | `conv.FormatBool` |
| `swag.FormatFloat32` | `conv.FormatFloat` |
| `swag.FormatFloat64` | `conv.FormatFloat` |
| `swag.FormatInt16` | `conv.FormatInteger` |
| `swag.FormatInt32` | `conv.FormatInteger` |
| `swag.FormatInt64` | `conv.FormatInteger` |
| `swag.FormatInt8` | `conv.FormatInteger` |
| `swag.FormatUint16` | `conv.FormatUinteger` |
| `swag.FormatUint32` | `conv.FormatUinteger` |
| `swag.FormatUint64` | `conv.FormatUinteger` |
| `swag.FormatUint8` | `conv.FormatUinteger` |
| `swag.Int` | `conv.Pointer` |
| `swag.Int32` | `conv.Pointer` |
| `swag.Int32Map` | `conv.PointerMap` |
| `swag.Int32Slice` | `conv.PointerSlice` |
| `swag.Int32Value` | `conv.Value` |
| `swag.Int32ValueMap` | `conv.ValueMap` |
| `swag.Int32ValueSlice` | `conv.ValueSlice` |
| `swag.Int64` | `conv.Pointer` |
| `swag.Int64Map` | `conv.PointerMap` |
| `swag.Int64Slice` | `conv.PointerSlice` |
| `swag.Int64Value` | `conv.Value` |
| `swag.Int64ValueMap` | `conv.ValueMap` |
| `swag.Int64ValueSlice` | `conv.ValueSlice` |
| `swag.IntMap` | `conv.PointerMap` |
| `swag.IntSlice` | `conv.PointerSlice` |
| `swag.IntValue` | `conv.Value` |
| `swag.IntValueMap` | `conv.ValueMap` |
| `swag.IntValueSlice` | `conv.ValueSlice` |
| `swag.IsFloat64AJSONInteger` | `conv.IsFloat64AJSONInteger` |
| `swag.String` | `conv.Pointer` |
| `swag.StringMap` | `conv.PointerMap` |
| `swag.StringSlice` | `conv.PointerSlice` |
| `swag.StringValue` | `conv.Value` |
| `swag.StringValueMap` | `conv.ValueMap` |
| `swag.StringValueSlice` | `conv.ValueSlice` |
| `swag.Time` | `conv.Pointer` |
| `swag.TimeMap` | `conv.PointerMap` |
| `swag.TimeSlice` | `conv.PointerSlice` |
| `swag.TimeValue` | `conv.Value` |
| `swag.TimeValueMap` | `conv.ValueMap` |
| `swag.TimeValueSlice` | `conv.ValueSlice` |
| `swag.Uint` | `conv.Pointer` |
| `swag.Uint16` | `conv.Pointer` |
| `swag.Uint16Map` | `conv.PointerMap` |
| `swag.Uint16Slice` | `conv.PointerSlice` |
| `swag.Uint16Value` | `conv.Value` |
| `swag.Uint16ValueMap` | `conv.ValueMap` |
| `swag.Uint16ValueSlice` | `conv.ValueSlice` |
| `swag.Uint32` | `conv.Pointer` |
| `swag.Uint32Map` | `conv.PointerMap` |
| `swag.Uint32Slice` | `conv.PointerSlice` |
| `swag.Uint32Value` | `conv.Value` |
| `swag.Uint32ValueMap` | `conv.ValueMap` |
| `swag.Uint32ValueSlice` | `conv.ValueSlice` |
| `swag.Uint64` | `conv.Pointer` |
| `swag.Uint64Map` | `conv.PointerMap` |
| `swag.Uint64Slice` | `conv.PointerSlice` |
| `swag.Uint64Value` | `conv.Value` |
| `swag.Uint64ValueMap` | `conv.ValueMap` |
| `swag.Uint64ValueSlice` | `conv.ValueSlice` |
| `swag.UintMap` | `conv.PointerMap` |
| `swag.UintSlice` | `conv.PointerSlice` |
| `swag.UintValue` | `conv.Value` |
| `swag.UintValueMap` | `conv.ValueMap` |
| `swag.UintValueSlice` | `conv.ValueSlice` |

### `jsonutils` — `github.com/go-openapi/swag/jsonutils`

7 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.ConcatJSON` | `jsonutils.ConcatJSON` |
| `swag.DynamicJSONToStruct` | `jsonutils.FromDynamicJSON` |
| `swag.FromDynamicJSON` | `jsonutils.FromDynamicJSON` |
| `swag.JSONMapItem` | `jsonutils.JSONMapItem` |
| `swag.JSONMapSlice` | `jsonutils.JSONMapSlice` |
| `swag.ReadJSON` | `jsonutils.ReadJSON` |
| `swag.WriteJSON` | `jsonutils.WriteJSON` |

### `typeutils` — `github.com/go-openapi/swag/typeutils`

1 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.IsZero` | `typeutils.IsZero` |

### `stringutils` — `github.com/go-openapi/swag/stringutils`

3 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.ContainsStringsCI` | `stringutils.ContainsStringsCI` |
| `swag.JoinByFormat` | `stringutils.JoinByFormat` |
| `swag.SplitByFormat` | `stringutils.SplitByFormat` |

### `slices` — `slices (Go standard library)`

1 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.ContainsStrings` | `slices.Contains` |

### `loading` — `github.com/go-openapi/swag/loading`

6 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.LoadFromFileOrHTTP` | `loading.LoadFromFileOrHTTP` |
| `swag.LoadFromFileOrHTTPWithTimeout` | `loading.LoadFileOrHTTP` |
| `swag.LoadStrategy` | `loading.LoadStrategy` |
| `swag.YAMLData` | `loading.YAMLData` |
| `swag.YAMLDoc` | `loading.YAMLDoc` |
| `swag.YAMLMatcher` | `loading.YAMLMatcher` |

### `yamlutils` — `github.com/go-openapi/swag/yamlutils`

2 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.BytesToYAMLDoc` | `yamlutils.BytesToYAMLDoc` |
| `swag.YAMLToJSON` | `yamlutils.YAMLToJSON` |

### `fileutils` — `github.com/go-openapi/swag/fileutils`

5 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.File` | `fileutils.File` |
| `swag.FindInGoSearchPath` | `fileutils.FindInGoSearchPath` |
| `swag.FindInSearchPath` | `fileutils.FindInSearchPath` |
| `swag.FullGoSearchPath` | `fileutils.FullGoSearchPath` |
| `swag.GOPATHKey` | `fileutils.GOPATHKey` |

### `netutils` — `github.com/go-openapi/swag/netutils`

1 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.SplitHostPort` | `netutils.SplitHostPort` |

### `mangling` — `github.com/go-openapi/swag/mangling`

2 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.AddInitialisms` | `mangling.WithAdditionalInitialisms` |
| `swag.GoNamePrefixFunc` | `mangling.WithGoNamePrefixFunc` |

### `cmdutils` — `github.com/go-openapi/swag/cmdutils`

1 symbol(s).

| Old (root `swag`) | New |
| --- | --- |
| `swag.CommandLineOptionsGroup` | `cmdutils.CommandLineOptionsGroup` |
