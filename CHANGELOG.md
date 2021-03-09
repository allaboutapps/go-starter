# Changelog

- All notable changes to this project will be documented in this file.
- The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
- We **do not follow [semantic versioning](https://semver.org/)**.
- There are **no git tags**. 
- All changes are solely **tracked by date**. 
- The latest `master` is considered **stable** and should be periodically merged into our customer projects.
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