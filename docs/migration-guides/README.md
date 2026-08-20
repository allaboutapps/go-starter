# Migration guides

Step-by-step instructions for go-starter changes that **require action in projects
based on this template**, beyond merging the diff.

These live in the repository (rather than only in the
[wiki](https://github.com/allaboutapps/go-starter/wiki)) so that a project which
merges go-starter also receives the guide it needs, at the revision it needs it.

A guide belongs here when the change cannot be absorbed by merging alone — because
it breaks compilation of application code, needs a data migration, or silently
changes runtime behaviour.

| Guide | Applies when |
| --- | --- |
| [Postgres 17 → 18](postgres-17-to-18.md) | Always. Requires a database major upgrade, and a stale volume mount silently discards writes. |
| [`swag` root package → submodules](swag-root-package-to-submodules.md) | Your code calls `swag.String`, `swag.Int64`, `swag.IsZero`, … (124 deprecated symbols). |
| [go-swagger 0.36 and golangci-lint 2.13](go-swagger-and-golangci-lint.md) | You modified `api/templates/server/*.gotmpl`, or your lint pipeline turns red after the upgrade. |
