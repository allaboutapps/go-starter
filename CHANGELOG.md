# Changelog

- All notable changes to this project will be documented in this file.
- The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
- We **do not follow [semantic versioning](https://semver.org/)**.
- There are **no git tags**. 
- All changes are solely **tracked by date**. 
- The latest `master` is considered **stable** and should be periodically merged into our customer projects.

## 2021-03-15
### Changed
- Upgrades `go.mod`:
  - `github.com/volatiletech/sqlboiler/v4@v4.5.0`
  - `github.com/rogpeppe/go-internal@v1.8.0`
  - `golang.org/x/crypto@v0.0.0-20210314154223-e6e6c4f2bb5b`
  - ~`golang.org/x/sys@v0.0.0-20210314195730-07df6a141424`~
  - `golang.org/x/sys@v0.0.0-20210315160823-c6e025ad8005`
  - `google.golang.org/api@v0.42.0`
- `make help` no longer reports `(opt)` flagged targets, use `make help-all` instead.
- `make tools` now executes `go install {}` in parallel
- `make info` now fetches infomation in parallel
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