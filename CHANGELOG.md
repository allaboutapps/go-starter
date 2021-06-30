# Changelog

- All notable changes to this project will be documented in this file.
- The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
- We **do not follow [semantic versioning](https://semver.org/)**.
- There are **no git tags**. 
- All changes are solely **tracked by date**. 
- The latest `master` is considered **stable** and should be periodically merged into our customer projects.

## Unreleased
### Changed

## 2021-06-30
### Changed
- **BREAKING** Switched from [`golint`](https://github.com/golang/lint) to [`revive`](https://github.com/mgechev/revive)
  - [`golint` is deprecated](https://github.com/golang/go/issues/38968).
  - [`revive`](https://github.com/mgechev/revive) is considered to be a drop-in replacement for `golint`, however this change still might lead to breaking changes in your codebase.
- **BREAKING** `make lint` no longer uses `--fast` when calling `golangci-lint`
  - Up until now, `make lint` also ran `golangci-lint` using the `--fast` flag to remain consistent with the linting performed by VSCode automatically.
  - As running only fast linters in both steps meant skipping quite a few validations (only 4/13 enabled linters are actually active), a decision has been made to break consistency between the two lint-steps and perform "full" linting during the build pipeline.
  - This change could potentially bring up additional warnings and thus fail your build until fixed.
- **BREAKING** `gosec` is now also applied to test packages
  - All linters are now applied to every source code file in this project, removing the previous exclusion of `gosec` from test files/packages
  - As `gosec` might (incorrectly) detect some hardcoded credentials in your tests (variable names such as `passwordResetLink` get flagged), this change might require some fixes after merging.
- Extended auth middleware to allow for multiple auth token sources
  - Default token validator uses access token table, maintaining previous behavior without any changes required.
  - Token validator can be changed to e.g. use separate API keys for specific endpoints, allowing for more flexibility if so desired.
- Changed `util.LogFromContext` to always return a valid logger
  - Helper no longer returns a disabled logger if context provided did not have an associated logger set (e.g. by middleware). If you still need to disable the logger for a certain context/function, use `util.DisableLogger(ctx, true)` to force-disable it.
  - Added request ID to context in logger middleware.
- Extended DB query helpers
  - Fixed TSQuery escaping, should now properly handle all type of user input.
  - Implemented helper for JSONB queries (see `ExampleWhereJSON` for implementation details).
  - Added `LeftOuterJoin` helper, similar to already existing `LeftJoin` variants.
  - Managed transactions (via `WithTransaction`) can now have their options configured via `WithConfiguredTransaction`.
  - Added util to combine query mods with `OR` expression.
- Implemented middleware for parsing `Cache-Control` header
  - Allows for cache handling in relevant services, parsed directive is stored in request context.
  - New middleware is enabled by default, can be disabled via env var (`SERVER_ECHO_ENABLE_CACHE_CONTROL_MIDDLEWARE`).
- Added extra misc. helpers
  - Extra helpers for slice handling and generating random strings from a given character set have been included (`util.ContainsAllString`, `util.UniqueString`, `util.GenerateRandomString`).
  - Added util to check whether current execution runs inside a test environment (`util.RunningInTest`).
- Test and snapshot util improvements
  - Added `snapshoter.SaveU` as a shorthand for updating a single test
  - Implemented `GenericArrayPayload` with respective request helpers for array payloads in tests
  - Added VScode launch task for updating all snapshots in a single test file

## 2021-06-29
### Changed
- We now directly bake the `gsdev` cli "bridge" (it actually just runs `go run -tags scripts /app/scripts/main.go "$@"`) into the `development` stage of our `Dockerfile` and create it at `/usr/bin/gsdev` (requires `./docker-helper.sh --rebuild`).
  - `gsdev` was previously symlinked to `/app/bin` from `/app/scripts/gsdev` (within the projects' workspace) and `chmod +x` via the `Makefile` during `init`.
  - However this lead to problems with WSL2 VSCode related development setups (always dirty git workspaces as WSL2 tries to prevent `+x` flags). 
  - **BREAKING** encountered at **2021-06-30**: Upgrading your project via `make git-merge-go-starter` if you already have installed our previous `gsdev` approach from **2021-06-22** may require additional steps:
    - It might be necessary to unlink the current `gsdev` symlink residing at `/app/bin/gsdev` before merging up (as this symlinked file will no longer exist)!
    - Do this by issuing `rm -f /app/bin/gsdev` which will remove the symlink which pointed to the previous (now gone bash script) at `/app/scripts/gsdev`.
    - It might also be handy to install the newer variant directly into your container (without requiring a image rebuild). Do this by:
      - `sudo su` to become root in the container,
      - issuing the following command: `printf '#!/bin/bash\nset -Eeo pipefail\ncd /app && go run -tags scripts ./scripts/main.go "$@"' > /usr/bin/gsdev && chmod 755 /usr/bin/gsdev` (in sync with what we do in our `Dockerfile`) and
      - `[CTRL + c]` to return to being the `development` user within your container.

## 2021-06-24
### Changed
- Introduces GitHub Actions docker layer caching via docker buildx. For details see `.github/workflows/build-test.yml`.
- Upgrades:
  - Bump golang from 1.16.4 to [1.16.5](https://groups.google.com/g/golang-announce/c/RgCMkAEQjSI/m/r_EP-NlKBgAJ)
  - golangci-lint@[v1.41.1](https://github.com/golangci/golangci-lint/releases/tag/v1.41.1)
  - Bump github.com/rs/zerolog from 1.22.0 to [1.23.0](https://github.com/allaboutapps/go-starter/pull/92)
  - Bump github.com/go-openapi/runtime from 0.19.28 to 0.19.29
  - Bump github.com/volatiletech/sqlboiler/v4 from 4.5.0 to [4.6.0](https://github.com/volatiletech/sqlboiler/blob/HEAD/CHANGELOG.md#v460---2021-06-06)
  - Bump github.com/rubenv/sql-migrate v0.0.0-20210408115534-a32ed26c37ea to v0.0.0-20210614095031-55d5740dbbcc
  - Bump github.com/spf13/viper v1.7.1 to v1.8.0
  - Bump golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a to v0.0.0-20210616213533-5ff15b29337e
  - Bump golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea to v0.0.0-20210616094352-59db8d763f22
  - Bump google.golang.org/api v0.47.0 to v0.49.0
- Fixes linting within `/scripts/**/*.go`, now activated by default.

## 2021-06-22
### Changed
- Development scripts are no longer called via `go run [script]` but via `gsdev`:
  - The `gsdev` cli is our new entrypoint for development workflow specific scripts, these scripts are not available in the final `app` binary.
  - All previous `go run` scripts have been moved to their respective `/scripts/cmd` cli entrypoint + internal implementation within `/scripts/internal/**`.
  - Please use `gsdev --help` to get an overview of available development specific commands.
  - `gsdev` relys on a tiny helper bash script `scripts/gsdev` which gets symlinked to `/app/bin` on `make init`.
  - Use `make test-scripts` to run tests regarding these internal scripts within `/scripts/**/*_test.go`.
  - We now enforce that all `/scripts/**/*.go` files set the `// +build scripts` build tag. We do this to ensure these files are not directly depended upon from the actual `app` source-code within `/internal`.
- VSCode's `.devcontainer/devcontainer.json` now defines that the go tooling must use the `scripts` build tag for its IntelliSense. This is neccessary to still get proper code-completion when modifying resources at `/scripts/**/*.go`. You may need to reattach VSCode and/or run `./docker-helper.sh --rebuild`.

### Added
- Scaffolding tool to quickly generate generic CRUD endpoint stubs. Usage: `gsdev scaffold [resource name] [flags]`, also see `gsdev scaffold --help`.

## 2021-05-26
### Changed
- Scans for [CVE-2020-26160](https://nvd.nist.gov/vuln/detail/CVE-2020-26160) also match for our final `app` binary, however, we do not use `github.com/dgrijalva/jwt-go` as part of our auth logic. This dependency is mostly here because of child dependencies, that yet need to upgrade to `>=v4.0.0`. Therefore, we currently disable this CVE for scans in this project (via `.trivyignore`).
- Upgrades `Dockerfile`: [`watchexec@v1.16.1`](https://github.com/watchexec/watchexec/releases/tag/cli-v1.16.1), [`lichen@v0.1.4`](https://github.com/uw-labs/lichen/releases/tag/v0.1.4) (requires `./docker-helper.sh --rebuild`).

## 2021-05-18
### Changed
- Upgraded `Dockerfile` to `golang:1.16.4`, `gotestsum@v1.6.4`, `golangci-lint@v1.40.1`, `watchexec@v1.16.0` (requires `./docker-helper.sh --rebuild`).
- Upgraded `go.mod`:
  - [github.com/labstack/echo/v4@v4.3.0](https://github.com/labstack/echo/releases/tag/v4.3.0)
  - [github.com/lib/pq@v1.10.2](https://github.com/lib/pq/releases/tag/v1.10.2)
  - [github.com/gabriel-vasile/mimetype@v1.3.0](https://github.com/gabriel-vasile/mimetype/releases/tag/v1.3.0)
  - `github.com/go-openapi/runtime@v0.19.28`
  - [github.com/rs/zerolog@v1.22.0](https://github.com/rs/zerolog/releases/tag/v1.22.0)
  - `github.com/rubenv/sql-migrate@v0.0.0-20210408115534-a32ed26c37ea`
  - `golang.org/x/crypto@v0.0.0-20210513164829-c07d793c2f9a`
  - `golang.org/x/sys@v0.0.0-20210514084401-e8d321eab015`
  - [google.golang.org/api@v0.46.0](https://github.com/googleapis/google-api-go-client/releases/tag/v0.46.0)
- GitHub Actions:
  -  Pin to `actions/checkout@v2.3.4`.
  -  Remove unnecessary `git checkout HEAD^2` in CodeQL step (Code Scanning recommends analyzing the merge commit for best results).
  -  Limit trivy and codeQL actions to `push` against `master` and `pull_request` against `master` to overcome read-only access workflow errors.

## 2021-04-27
### Added
- Adds `test.WithTestDatabaseFromDump*`, `test.WithTestServerFromDump` methods for writing tests based on a database dump file that needs to be imported first:
  - We dynamically setup IntegreSQL pools for all combinations passed through a `test.DatabaseDumpConfig{}` object:
    - `DumpFile string` is required, absolute path to dump file
    - `ApplyMigrations bool` optional, default `false`, automigrate after installing the dump
    - `ApplyTestFixtures bool` optional, default `false`, import fixtures after (migrating) installing the dump
  - `test.ApplyDump(ctx context.Context, t *testing.T, db *sql.DB, dumpFile string) error` may be used to apply a dump to an existing database connection.
  - As we have dedicated IntegreSQL pools for each combination, testing performance should be on par with the default IntegreSQL database pool. 
- Adds `test.WithTestDatabaseEmpty*` methods for writing tests based on an empty database (also a dedicated IntegreSQL pool). 
- Adds context aware `test.WithTest*Context` methods reusing the provided `context.Context` (first arg).
- Adds `make sql-dump` command to easily create a dump of the local `development` database to `/app/dumps/development_YYYY-MM-DD-hh-mm-ss.sql` (.gitignored).

### Changed
- `test.ApplyMigrations(t *testing.T, db *sql.DB) (countMigrations int, err error)` is now public (e.g. for usage with `test.WithTestDatabaseEmpty*` or `test.WithTestDatabaseFromDump*`)
- `test.ApplyTestFixtures(ctx context.Context, t *testing.T, db *sql.DB) (countFixtures int, err error)` is now public (e.g. for usage with `test.WithTestDatabaseEmpty*` or `test.WithTestDatabaseFromDump*`)
- `internal/test/test_database_test.go` and `/app/internal/test/test_server_test.go` were massively refactored to allow for better extensibility later on (non breaking, all method signatures are backward-compatible).  

## 2021-04-12
### Added
- Adds echo `NoCache` middleware: Use `middleware.NoCache()` and `middleware.NoCacheWithConfig(Skipper)` to explicitly force browsers to never cache calls to these handlers/groups.

### Changed
- `/swagger.yml` and `/-/*` now explicity set no-cache headers by default, forcing browsers to re-execute calls each and every time.
- Upgrade [watchexec@v1.15.0](https://github.com/watchexec/watchexec/releases/tag/1.15.0) (requires `./docker-helper.sh --rebuild`).

## 2021-04-08
### Added
- Live-Reload for our swagger-ui is now available out of the box: 
  - [allaboutapps/browser-sync](https://hub.docker.com/r/allaboutapps/browser-sync) acts as proxy at [localhost:8081](http://localhost:8081/).
  - Requires `./docker-helper.sh --up`.
  - Best used in combination with `make watch-swagger` (still refreshes `make all` or `make swagger` of course).

### Changed
- Upgrades to [swaggerapi/swagger-ui:v3.46.0](https://github.com/swagger-api/swagger-ui/tree/v3.46.0) from [swaggerapi/swagger-ui:v3.28.0](https://github.com/swagger-api/swagger-ui/compare/v3.28.0...v3.46.0)
- Upgrades to [github.com/labstack/echo@v4.2.2](https://github.com/labstack/echo/releases/tag/v4.2.2)
- `golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2`
- Upgrades to [google.golang.org/api@v0.44.0](https://github.com/googleapis/google-api-go-client/releases/tag/v0.44.0)

## 2021-04-07
### Changed
-  Moved `/api/main.yml` to `/api/config/main.yml` to overcome path resolve issues (`../definitions`) with the VSCode [42crunch.vscode-openapi](https://github.com/42Crunch/vscode-openapi) extension (auto-included in our devContainer) and our go-swagger concat behaviour. 
- Updated [api/README.md](https://github.com/allaboutapps/go-starter/blob/master/api/README.md) information about `/api/swagger.yml` generation logic and changed `make swagger-concat` accordingly

## 2021-04-02
### Changed
- Bump [golang from v1.16.2 to v1.16.3](https://github.com/golang/go/issues?q=milestone%3AGo1.16.3+label%3ACherryPickApproved) (requires `./docker-helper.sh --rebuild`).

## 2021-04-01
### Changed
- Bump golang.org/x/crypto@v0.0.0-20210322153248-0c34fe9e7dc2
- Bump golang.org/x/sys@v0.0.0-20210331175145-43e1dd70ce54
- Bump [github.com/go-openapi/swag@v0.19.15](https://github.com/allaboutapps/go-starter/pull/71)
- Bump [github.com/go-openapi/strfmt@v0.20.1](https://github.com/allaboutapps/go-starter/pull/70)

## 2021-03-30
### Changed
- Bump [github.com/gotestyourself/gotestsum@v1.6.3](https://github.com/gotestyourself/gotestsum/releases/tag/v1.6.3) (requires `./docker-helper.sh --rebuild`).

## 2021-03-26
### Changed
- Bump [golangci-lint@v1.39.0](https://github.com/golangci/golangci-lint/releases/tag/v1.39.0) (requires `./docker-helper.sh --rebuild`).

## 2021-03-25
### Changed
- Bump github.com/rs/zerolog from [1.20.0 to 1.21.0](https://github.com/allaboutapps/go-starter/pull/69)
- Bump google.golang.org/api from [0.42.0 to 0.43.0](https://github.com/allaboutapps/go-starter/pull/68)

## 2021-03-24

### Changed
- We no longer do explicit calls to `t.Parallel()` in our go-starter tests (except autogenerated code). For the reasons why see [FAQ: Should I use `t.Parallel()` in my tests?](https://github.com/allaboutapps/go-starter/wiki/FAQ#should-i-use-tparallel-in-my-tests).
- Switched to [github.com/uw-labs/lichen](https://github.com/uw-labs/lichen) for getting license information of embedded dependencies in our final `./bin/app` binary.
- The following make targets are no longer flagged as `(opt)` and thus move into the main `make help` target (use `make help-all` to see all targets): 
  - `make lint`: Runs golangci-lint and make check-*.
  - `make go-test-print-slowest`: Print slowest running tests (must be done after running tests).
  - `make get-licenses`: Prints licenses of embedded modules in the compiled bin/app.
  - `make get-embedded-modules`: Prints embedded modules in the compiled bin/app.
  - `make clean`: Cleans ./tmp and ./api/tmp folder.
  - `make get-module-name`: Prints current go module-name (pipeable).
- `make check-gen-dirs` now ignores `.DS_Store` within `/internal/models/**/*` and `/internal/types/**/*` and echo an errors detailing what happened.
- Upgrade to [`github.com/go-openapi/runtime@v0.19.27`](https://github.com/go-openapi/runtime/compare/v0.19.26...v0.19.27)
## 2021-03-16
### Changed
- `make all` no longer executes `make info` as part of its targets chain.
  - It's very common to use `make all` multiple times per day during development and thats fine! However, the output of `make info` is typically ignored by our engineers (if they explicitly want this information, they use `make info`). So `make all` was just too spammy in it's previous form.
  - `make info` does network calls and typically takes around 5sec to execute. This slowdown is not acceptable when running `make all`, especially if the information it provides isn't used anyways.
  - Thus: Just trigger `make info` manually if you need the information of the `[spec DB]` structure, current `[handlers]` and `[go.mod]` information. Furthermore you may also visit `tmp/.info-db`, `tmp/.info-handlers` and `tmp/.info-go` after triggering `make info` as we store this information there after a run. 

## 2021-03-15
### Changed
- Upgrades `go.mod`:
  - [`github.com/volatiletech/sqlboiler/v4@v4.5.0`](https://github.com/volatiletech/sqlboiler/blob/master/CHANGELOG.md#v450---2021-03-14)
  - [`github.com/rogpeppe/go-internal@v1.8.0`](https://github.com/rogpeppe/go-internal/releases/tag/v1.8.0)
  - `golang.org/x/crypto@v0.0.0-20210314154223-e6e6c4f2bb5b`
  - ~`golang.org/x/sys@v0.0.0-20210314195730-07df6a141424`~
  - `golang.org/x/sys@v0.0.0-20210315160823-c6e025ad8005`
  - [`google.golang.org/api@v0.42.0`](https://github.com/googleapis/google-api-go-client/releases/tag/v0.42.0)
- `make help` no longer reports `(opt)` flagged targets, use `make help-all` instead.
- `make tools` now executes `go install {}` in parallel
- `make info` now fetches information in parallel
- Seeding: Switch to `db|dbUtil.WithTransaction` instead of manually managing the db transaction. *Note*: We will enforce using `WithTransaction` instead of manually managing the life-cycle of db transactions through a custom linter in an upcoming change. It's way safer and manually managing db transactions only makes sense in very very special cases (where you will be able to opt-out via linter excludes). Also see [What's `WithTransaction`, shouldn't I use `db.BeginTx`, `db.Commit`, and `db.Rollback`?](https://github.com/allaboutapps/go-starter/wiki/FAQ#whats-withtransaction-shouldnt-i-use-dbbegintx-dbcommit-and-dbrollback).

### Fixed
- The correct implementation of `(util|scripts).GetProjectRootDir() string` now gets automatically selected based on the `scripts` build tag.
  - We currently have 2 different `GetProjectRootDir()` implementations and each one is useful on its own:
    - `util.GetProjectRootDir()` gets used while `app` or `go test` runs and resolves in the following way: use `PROJECT_ROOT_DIR` (if set), else default to the resolved path to the executable unless we can't resolve that, then **panic**!
    - `scripts.GetProjectRootDir()` gets used while **generation time** (`make go-generate`) and resolves in the following way: use `PROJECT_ROOT_DIR` (if set), otherwise default to `/app` (baked, as we can assume we are in the `development` container).
  - `/internal/util/(get_project_root_dir.go|get_project_root_dir_scripts.go)` is now introduced to automatically switch to the proper implementation based on the `// +build !scripts` or `// +build scripts` build tag, thus it's now consistent to import `util.GetProjectRootDir()`, especially while handler generation time (`make go-generate`).

## 2021-03-12
### Changed
- Upgrades to `golang@v1.16.2` (use `./docker-helper.sh --rebuild`).
- Silence resolve of `GO_MODULE_NAME` if `go` was not found in path (typically host env related).

## 2021-03-11
### Added
- `make build` (`make go-build`) now sets `internal/config.ModuleName`, `internal/config.Commit` and `internal/config.BuildDate` via `-ldflags`.
  - `/-/version` (mgmt key auth) endpoint is now available, prints the same as `app -v`.
  - `app -v` is now available and prints out buildDate and commit. Sample:
```bash
app -v
allaboutapps.dev/aw/go-starter @ 19c4cdd0da151df432cd5ab33c35c8987b594cac (2021-03-11T15:42:27+00:00)
```
### Changed
- Upgrades to `golang@v1.16.1` (use `./docker-helper.sh --rebuild`).
- Updates `google.golang.org/api@v0.41.0`, `github.com/gabriel-vasile/mimetype@v1.2.0` ([new supported formats](https://github.com/gabriel-vasile/mimetype/tree/v1.2.0)), `golang.org/x/sys`
- Removed `**/.git` from `.dockerignore` (`builder` stage) as we want the local git repo available while running `make go-build`.
- `app --help` now prominently includes the module name of the project. 
- Prominently recommend `make force-module-name` after running `make git-merge-go-starter` to fix all import paths.
## 2021-03-09
### Added
- Introduces `CHANGELOG.md`
### Changed
- `make git-merge-go-starter` now uses `--allow-unrelated-histories` by default.
  - `README.md` and FAQ now mention that it's recommended to execute `make git-merge-go-starter` during project setup (especially for single commit generated from template project project setups).
  - See [FAQ: I want to compare or update my project/fork to the latest go-starter master.](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-want-to-compare-or-update-my-projectfork-to-the-latest-go-starter-master)
- Various typos in `README.md` and `Makefile`.
- Upgrade to [`golangci-lint@v1.38.0`](https://github.com/golangci/golangci-lint/releases/tag/v1.38.0)

## 2021-03-08
### Added
- `allaboutapps/nullable` is now included by default. See [#58](https://github.com/allaboutapps/go-starter/pull/58), [FAQ: I need an optional Swagger payload property that is nullable!](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-need-an-optional-swagger-payload-property-that-is-nullable)
### Changed
- Upgrade to [`labstack/echo@v4.2.1`](https://github.com/labstack/echo/releases/tag/v4.2.1), [`lib/pq@v1.10.0`](https://github.com/lib/pq/releases/tag/v1.10.0)

## 2021-02-23
### Deprecated
- `util.BindAndValidate` is now marked as deprecated as [`labstack/echo@v4.2.0`](https://github.com/labstack/echo/releases/tag/v4.2.0) exposes a more granular binding through its `DefaultBinder`.
### Added
- The more specialized variants `util.BindAndValidatePathAndQueryParams` and `util.BindAndValidateBody` are now available. See [`/internal/util/http.go`](https://github.com/allaboutapps/go-starter/blob/master/internal/util/http.go#L87).

### Changed
- `golang@v1.16.0`
- [`labstack/echo@v4.2.0`](https://github.com/labstack/echo/releases/tag/v4.2.0)

## 2021-02-16
### Changed
- Upgrades to [`pgFormatter@v5.0.0`](https://github.com/darold/pgFormatter/releases) + forces VSCode to use that version within the devcontainer through it's extension.

## 2021-02-09
### Changed
- `golang@v1.15.8`, `go-swagger@v0.26.1`

## 2021-02-01
### Changed
```
- Dockerfile updates: 
  - golang@1.15.7
  - apt add icu-devtools (VSCode live sharing)
  - gotestsum@1.6.1
  - golangci-lint@v1.36.0
  - goswagger@v0.26.0
- go.mod:
  - sqlboiler@4.4.0
  - swag@0.19.3
  - strfmt@0.20.0
  - testify@1.7.0
  - go-openapi/runtime@v0.19.26
  - go-openapi/swag@v0.19.13
  - go-openapi/validate@v0.20.1
  - jordan-wright/email
  - rogpeppe/go-internal@v1.7.0
  - golang.org/x/crypto
  - golang.org/x/sys
  - google.golang.org/api@v0.38.0
```

### Fixed
- disabled goswagger generate server flag `--keep-spec-order` as relative resolution of its temporal created yml file is broken - see https://github.com/go-swagger/go-swagger/issues/2216

## 2020-11-04
### Added
- `make watch-swagger` and `make watch-sql`
### Changed
- sqlboiler@4.3.0

## 2020-11-02
### Added
- `make watch-tests`: Watches .go files and runs package tests on modifications.


## 2020-09-30
### Added
- `pprof` handlers, see [FAQ: I need to (remotely) pprof my running service!](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-need-to-remotely-pprof-my-running-service)


## 2020-09-24
### Added
- `make git-merge-go-starter`, see [FAQ: I want to compare or update my project/fork to the latest go-starter master.](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-want-to-compare-or-update-my-projectfork-to-the-latest-go-starter-master)

## 2020-09-22
### Added
- `app probe readiness` and `app probe liveness` sub-commands.
- `/-/ready` and `/-/healthy` handlers.


## 2020-09-16
### Changed
- Force VSCode to use our installed version of golang-cilint
- All `*.go` files in `/scripts` now use the build tag `scripts` so we can ensure they are not compiled into the final `app` binary.

### Added
- `go.not` file to ensure certain generation- / test-only dependencies don't end up in the final `app` binary. Automatically checked though `make` (sub-target `make check-embedded-modules-go-not`).

## 2020-09-11
- Switch to `distroless` as final app stage, see [FAQ: Should I use distroless/base or debian:buster-slim in the Dockerfile app stage?](https://github.com/allaboutapps/go-starter/wiki/FAQ#should-i-use-distrolessbase-or-debianbuster-slim-in-the-dockerfile-app-stage)
