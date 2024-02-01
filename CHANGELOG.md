# Changelog

- All notable changes to this project will be documented in this file.
- The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
- We do not follow [semantic versioning](https://semver.org/).
- All changes are solely **tracked by date** and have a **git tag** available (from 2021-10-19 onwards):
  - Git tags are formatted like `go-starter-YYYY-MM-DD`. See [GitHub tags](https://github.com/allaboutapps/go-starter/tags) for all available go-starter git tags.
  - The latest `master` is considered **stable** and should be periodically merged into our customer projects.
- Please follow the update process in *[I just want to update / upgrade my project!](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-just-want-to-update--upgrade-my-project)*.

## Unreleased

## 2024-02-01
- [Persist bash history in development container](https://code.visualstudio.com/remote/advancedcontainers/persist-bash-history) (requires `./docker-helper.sh --rebuild`).
  - Your commands are now persisted between your development container restarts / rebuilds, making it easier to re-run specific commands you've previously executed (e.g. that one go command you cannot remember).
- Hotfix [types.NullDecimal error](https://github.com/volatiletech/sqlboiler/issues/1234) by  downgrading indirect `github.com/ericlagergren/decimal@v0.0.0-20190420051523-6335edbaa640`.
  - Note that we do not pin it in direct dependencies, as this downgrade is already in [SQLBoilers master](https://github.com/volatiletech/sqlboiler/commit/bc59c158590800f7810cce241a12d572e898014f) anyways.

## 2024-01-31
- Migration to Docker Compose V2 ([Docker Compose Docs](https://docs.docker.com/compose/reference/)), thx [@eklatzer](https://github.com/eklatzer)
- Upgrade to [IntegreSQL v1.1.0](https://github.com/allaboutapps/integresql/blob/v1.1.0/CHANGELOG.md#v110)
- Switch [from Go 1.20.3 Go 1.21.6](https://go.dev/doc/devel/release#go1.21.0) (requires `./docker-helper.sh --rebuild`).
- Fix premature optimization in `make swagger` -> `make swagger-generate` (rm `rsync` with `--size-only`), thx [@eklatzer](https://github.com/eklatzer)
- Dockerfile deps upgrade:
  - Upgrade pgFormatter [from v5.3 to v5.5](https://github.com/darold/pgFormatter/releases/tag/v5.5)
  - Upgrade gotestsum [from 1.9.0 to 1.11.0](https://github.com/gotestyourself/gotestsum/releases/tag/v1.11.0)
  - Upgrade golangci-lint [from 1.52.2 to 1.55.2](https://github.com/golangci/golangci-lint/releases/tag/v1.55.2)
  - Upgrade watchexec [from 1.20.6 to 1.25.1](https://github.com/watchexec/watchexec/releases/tag/v1.25.1)
- `go.mod` upgrades
  - Minor: [Bump github.com/BurntSushi/toml from v1.2.1 to v1.3.2](https://github.com/BurntSushi/toml)
  - Minor: [Bump github.com/davecgh/go-spew from v1.1.1 to v1.1.2-0.20180830191138-d8f796af33cc](https://github.com/davecgh/go-spew/commit/d8f796af33cc11cb798c1aaeb27a4ebc5099927d)
  - Minor: [Bump github.com/gabriel-vasile/mimetype from v1.4.1 to v1.4.3](https://github.com/gabriel-vasile/mimetype)
  - Minor: [Bump github.com/go-openapi/errors from v0.20.3 to v0.21.0](https://github.com/go-openapi/errors)
  - Minor: [Bump github.com/go-openapi/runtime from v0.25.0 to v0.27.1](https://github.com/go-openapi/runtime)
  - Minor: [Bump github.com/go-openapi/strfmt from v0.21.3 to v0.22.0](https://github.com/go-openapi/strfmt)
  - Minor: [Bump github.com/go-openapi/swag from v0.22.3 to v0.22.9](https://github.com/go-openapi/swag)
  - Minor: [Bump github.com/go-openapi/validate from v0.22.0 to v0.22.6](https://github.com/go-openapi/validate)
  - Minor: [Bump github.com/labstack/echo/v4 from v4.9.1 to v4.11.4](https://github.com/labstack/echo)
  - Minor: [Bump github.com/lib/pq from v1.10.7 to v1.10.9](https://github.com/lib/pq)
  - Minor: [Bump github.com/nicksnyder/go-i18n/v2 from v2.2.1 to v2.4.0](https://github.com/nicksnyder/go-i18n)
  - Minor: [Bump github.com/pmezard/go-difflib from v1.0.0 to v1.0.1-0.20181226105442-5d4384ee4fb2](https://github.com/pmezard/go-difflib) (deprecated)
  - Minor: [Bump github.com/rs/zerolog from v1.28.0 to v1.31.0](https://github.com/rs/zerolog)
  - Minor: [Bump github.com/rubenv/sql-migrate from v1.2.0 to v1.6.1](https://github.com/rubenv/sql-migrate)
  - Minor: [Bump github.com/spf13/cobra from v1.6.1 to v1.8.0](https://github.com/spf13/cobra)
  - Minor: [Bump github.com/spf13/viper from v1.14.0 to v1.18.2](https://github.com/spf13/viper)
  - Minor: [Bump github.com/stretchr/testify from v1.8.1 to v1.8.4](https://github.com/stretchr/testify)
  - Minor: [Bump github.com/subosito/gotenv from v1.4.1 to v1.6.0](https://github.com/subosito/gotenv)
  - Minor: [Bump github.com/volatiletech/sqlboiler/v4 from v4.13.0 to v4.16.1](https://github.com/volatiletech/sqlboiler)
  - Minor: [Bump github.com/volatiletech/strmangle from v0.0.4 to v0.0.6](https://github.com/volatiletech/strmangle)
  - Minor: [Bump golang.org/x/crypto from v0.3.0 to v0.18.0](https://golang.org/x/crypto)
  - Minor: [Bump golang.org/x/sys from v0.5.0 to v0.16.0](https://golang.org/x/sys)
  - Minor: [Bump golang.org/x/text from v0.7.0 to v0.14.0](https://golang.org/x/text)
  - Minor: [Bump google.golang.org/api from v0.103.0 to v0.161.0](https://google.golang.org/api)
  - Minor: [Bump xxxx from yyy to zzz](https://xxxx)
  - Replace: [github.com/rogpeppe/go-internal v1.9.0](https://github.com/rogpeppe/go-internal) with [golang.org/x/mod v0.14.0](https://pkg.go.dev/golang.org/x/mod)

## 2023-05-03
- Switch [from Go 1.19.3 to Go 1.20.3](https://go.dev/doc/devel/release#go1.20) (requires `./docker-helper.sh --rebuild`).
- Add new log configuration:
  - optional `output` param of `LoggerWithConfig` to redirect the log output
  - optional caller info switched on with `SERVER_LOGGER_LOG_CALLER`
- Minor: rename unused function parameters to fix linter errors
- Minor: update devcontainer.json syntax to remove deprecation warning
- Minor: add `GetFieldsImplementing` to utils and use it to easier add new fixture fields.
- `go.mod` changes:
  - Minor: [Bump github.com/golangci/golangci-lint from 1.50.1 to 1.52.2](https://github.com/golangci/golangci-lint/releases/tag/v1.52.2)
  - Minor: [Bump golang.org/x/net from 0.2.0 to 0.7.0](https://cs.opensource.google/go/x/net) (Fixing CVE-2022-41723)

## 2023-03-03
- Switch [from Go 1.17.9 to Go 1.19.3](https://go.dev/doc/devel/release#go1.19) (requires `./docker-helper.sh --rebuild`).
  - Major: Update base docker image from debian buster to bullseye
  - Minor: [Bump github.com/darold/pgFormatter from 5.2 to 5.3](https://github.com/darold/pgFormatter/releases/tag/v5.3)
  - Minor: [Bump github.com/gotestyourself/gotestsum from 1.8.0 to 1.9.0](https://github.com/gotestyourself/gotestsum/releases/tag/v1.9.0)
  - Minor: [Bump github.com/golangci/golangci-lint from 1.45.2 to 1.50.1](https://github.com/golangci/golangci-lint/releases/tag/v1.50.1)
  - Minor: [Bump github.com/uw-labs/lichen from 0.1.5 to 0.1.7](https://github.com/uw-labs/lichen/releases/tag/v0.1.7)
  - Minor: [Bump github.com/watchexec/watchexec from 1.18.11 to 1.20.6](https://github.com/watchexec/watchexec/releases/tag/v1.20.6)
  - Minor: [Bump github.com/mikefarah/yq from 4.24.2 to 4.30.5](https://github.com/mikefarah/yq/releases/tag/v4.30.5)
- Major: Upgrade distroless app image from base-debian10 to base-debian11
- Major: Dockerfile is now build to support amd64 and arm64 architecture
- Improve speed of `make swagger` when dealing with many files in `/api` by generating to a docker volume instead of the host filesystem, rsyncing only to changes into `/internal/types`. Furthermore split our swagger type generation and validation into two separate make targets, that can run concurrently (requires `./docker-helper.sh --rebuild`).
  - Note that `/app/api/tmp`, `/app/tmp` and `/app/bin` are now baked by proper docker volumes when using our `docker-compose.yml`/`./docker-helper.sh --up`. You **cannot** remove these directories directly inside the container (but its contents) and you can also no longer see its files on your host machine directly!
- Fix `make check-gen-dirs` false positives hidden files.
- Allow to trace/benchmark `Makefile` targets execution by using a custom shell wrapper for make execution. See `SHELL` and `.SHELLFLAGS` within `Makefile` and the custom `rksh` script in the root working directory. Usage: `MAKE_TRACE_TIME=true make <target>`
- `go.mod` changes:
  - Minor: [Bump github.com/BurntSushi/toml from 1.1.0 to 1.2.1](https://github.com/BurntSushi/toml/releases/tag/v1.2.1)
  - Minor: [Bump github.com/gabriel-vasile/mimetype from 1.4.0 to 1.4.1](https://github.com/gabriel-vasile/mimetype/releases/tag/v1.4.1)
  - Minor: [Bump github.com/go-openapi/errors from 0.20.2 to 0.20.3](https://github.com/go-openapi/errors/releases/tag/v0.20.3)
  - Minor: [Bump github.com/go-openapi/runtime from 0.23.3 to 0.25.0](https://github.com/go-openapi/runtime)
  - Minor: [Bump github.com/go-openapi/strfmt from 0.21.2 to 0.21.3](https://github.com/go-openapi/strfmt/releases/tag/v0.21.3)
  - Minor: [Bump github.com/go-openapi/swag from 0.21.1 to 0.22.3](https://github.com/go-openapi/swag/releases/tag/v0.22.3)
  - Minor: [Bump github.com/go-openapi/validate from 0.21.0 to 0.22.0](https://github.com/go-openapi/validate/releases/tag/v0.22.0)
  - Minor: [Bump github.com/labstack/echo/v4 from 4.7.2 to 4.9.1](https://github.com/labstack/echo/releases/tag/v4.9.1) (Fixing CVE-2022-40083)
  - Minor: [Bump github.com/lib/pq from 1.10.5 to 1.10.7](https://github.com/lib/pq/releases/tag/v1.10.7)
  - Minor: [Bump github.com/nicksnyder/go-i18n/v2 from 2.2.0 to 2.2.1](https://github.com/nicksnyder/go-i18n/releases/tag/v2.2.1)
  - Minor: [Bump github.com/rogpeppe/go-internal from 1.8.1 to 1.9.0](https://github.com/rogpeppe/go-internal/releases/tag/v1.9.0)
  - Minor: [Bump github.com/rs/zerolog from 1.26.1 to 1.28.0](https://github.com/rs/zerolog/releases/tag/v1.28.0)
  - Minor: [Bump github.com/rubenv/sql-migrate from 1.1.1 to 1.2.0](https://github.com/rubenv/sql-migrate/releases/tag/v1.2.0)
  - Minor: [Bump github.com/spf13/cobra from 1.4.0 to 1.6.1](https://github.com/spf13/cobra/releases/tag/v1.6.1)
  - Minor: [Bump github.com/spf13/viper from 1.10.1 to 1.14.0](https://github.com/spf13/viper/releases/tag/v1.14.0)
  - Minor: [Bump github.com/stretchr/testify from 1.7.1 to 1.8.1](https://github.com/stretchr/testify/releases/tag/v1.8.1)
  - Minor: [Bump github.com/subosito/gotenv from 1.2.0 to 1.4.1](https://github.com/subosito/gotenv/releases/tag/v1.4.1)
  - Minor: [Bump github.com/volatiletech/sqlboiler/v4 from 4.9.2 to v4.13.0](https://github.com/volatiletech/sqlboiler/blob/master/CHANGELOG.md#v4130---2022-08-28)
  - Minor: [Bump github.com/volatiletech/strmangle from 0.0.2 to 0.0.4](https://github.com/volatiletech/strmangle/releases/tag/v0.0.4) (changes in enum generation might require manual changes, minor changes)
  - Minor: [Bump golang.org/x/crypto from v0.0.0-20220411220226-7b82a4e95df4 to 0.3.0](https://cs.opensource.google/go/x/crypto)
  - Minor: [Bump golang.org/x/sys from v0.0.0-20220412211240-33da011f77ad to 0.2.0](https://cs.opensource.google/go/x/sys)
  - Minor: [Bump golang.org/x/text from 0.3.7 to 0.4.0](https://cs.opensource.google/go/x/text) (Fixing CVE-2022-32149)
  - Minor: [Bump google.golang.org/api from 0.74.0 to 0.103.0](https://github.com/googleapis/google-api-go-client/compare/v0.80.0...v0.103.0)

## 2022-09-13
- Hotfix: Previously there was a chance of recursive error wrapping within our [`internal/api/router/error_handler.go`](https://github.com/allaboutapps/go-starter/blob/master/internal/api/router/error_handler.go) in combination with `*echo.HTTPError`. We currently disable this wrapping (as not used anyways) and will schedule a cleaner update regarding this error augmentation approach.

## 2022-04-15
- Switch [from Go 1.17.1 to Go 1.17.9](https://go.dev/doc/devel/release#go1.17.minor) (requires `./docker-helper.sh --rebuild`).
- **BREAKING** Add [`tenv`](https://github.com/sivchari/tenv) and [`errorlint`](https://github.com/polyfloyd/go-errorlint) linter to our default `.golangci.yml` configuration.
  - We switch from `os.Setenv` to [`t.Setenv`](https://pkg.go.dev/testing#T.Setenv) within our own test code.
  - **NOTE**: If you have used `os.Setenv` within your `*_test.go` code previously, simply replace those calls by `t.Setenv`.
  - **NOTE**: The go-starter base code now properly uses `errors.Is` and `errors.As` for comparisons (and `%w` wrapping where really needed). For a good overview regarding error handling see [Effective Error Handling in Golang](https://earthly.dev/blog/golang-errors/). For example, if you receive linting errors, you'll need to change your code like this:
    - Wrong: `if err == sql.ErrNoRows {`
      - Valid: `if errors.Is(err, sql.ErrNoRows) {`
    - Wrong: `if err != sql.ErrConnDone {`
      - Valid:  `if !errors.Is(err, sql.ErrConnDone) {`
    - Wrong: `gErr := err.(*googleapi.Error)`, Valid:
      - `var gErr *googleapi.Error`
      - `ok := errors.As(err, &gErr)`
- `Dockerfile` development stage changes (requires `./docker-helper.sh --rebuild`):
  - Bump [golang](https://hub.docker.com/_/golang) base image from `golang:1.17.1-buster` to **`golang:1.17.8-buster`**.
  - Bump [pgFormatter](https://github.com/darold/pgFormatter) from v5.0 to [v5.2](https://github.com/darold/pgFormatter/releases/tag/v5.2)
  - Bump [golangci-lint](https://github.com/golangci/golangci-lint) from v1.42.1 to [v1.45.2](https://github.com/golangci/golangci-lint/blob/master/CHANGELOG.md#v1452)
  - Bump [lichen](https://github.com/uw-labs/lichen) from v0.1.4 to [v0.1.5](https://github.com/uw-labs/lichen/compare/v0.1.4...v0.1.5)
  - Bump [watchexec](https://github.com/watchexec/watchexec) from v1.17.0 to [v1.18.11](https://github.com/watchexec/watchexec/releases/tag/cli-v1.18.11) (+ switch from gnu to musl)
  - Bump [yq](https://github.com/mikefarah/yq) from v4.16.2 to [v4.24.2](https://github.com/mikefarah/yq/releases/tag/v4.24.2)
  - Bump [gotestsum](https://github.com/gotestyourself/gotestsum) from v1.7.0 to [v1.8.0](https://github.com/gotestyourself/gotestsum/releases/tag/v1.8.0)
  - Adds [tmux](https://github.com/tmux/tmux) (debian apt managed)
- `go.mod` changes:
  - Major: [Bump `github.com/rubenv/sql-migrate` from v0.0.0-20210614095031-55d5740dbbcc to v1.1.1](https://github.com/rubenv/sql-migrate/compare/55d5740dbbccbaa4934009263b37ba52d837241f...v1.1.1) (though this should not lead to any major changes)
  - Minor: [Bump github.com/volatiletech/sqlboiler/v4 from 4.6.0 to v4.9.2](https://github.com/volatiletech/sqlboiler/blob/v4.9.2/CHANGELOG.md#v492---2022-04-11) (your generated model might slightly change, minor changes).
    - Note that v5 will prefer wrapping errors (e.g. `sql.ErrNoRows`) to retain the stack trace, thus it's about time for us to start to enforce proper `errors.Is` checks in our codebase (see above).
  - Minor: [#178: Bump github.com/labstack/echo/v4 from 4.6.1 to 4.7.2](https://github.com/allaboutapps/go-starter/pull/178) (support for HEAD method query params binding, minor changes).
  - Minor: [#160: Bump github.com/rs/zerolog from 1.25.0 to 1.26.1](https://github.com/allaboutapps/go-starter/pull/160) (minor changes).
  - Minor: [#179: Bump github.com/nicksnyder/go-i18n/v2 from 2.1.2 to 2.2.0](https://github.com/allaboutapps/go-starter/pull/179) (minor changes).
  - Minor: [Bump `github.com/gabriel-vasile/mimetype` from v1.3.1 to v1.4.0](https://github.com/gabriel-vasile/mimetype/releases/tag/v1.4.0)
  - Minor: [Bump `github.com/go-openapi/runtime` from v0.22.0 to v0.23.3](https://github.com/go-openapi/runtime/compare/v0.22.0...v0.23.3)
  - Patch: [Bump `github.com/go-openapi/strfmt` from v0.21.1 to v0.21.2](https://github.com/go-openapi/strfmt/compare/v0.21.1...v0.21.2)
  - Patch: [Bump `github.com/go-openapi/validate` from v0.20.3 to v0.21.0](https://github.com/go-openapi/validate/compare/v0.20.3...v0.21.0)
  - Patch: [Bump `github.com/lib/pq` from v1.10.3 to v1.10.5](https://github.com/lib/pq/compare/v1.10.3...v1.10.5)
  - Patch: [Bump `github.com/rogpeppe/go-internal` from v1.8.0 to v1.8.1](https://github.com/rogpeppe/go-internal/releases/tag/v1.8.1)
  - Patch: [Bump `github.com/stretchr/testify` from v1.7.0 to v1.7.1](https://github.com/stretchr/testify/compare/v1.7.0...v1.7.1)
  - Patch: [Bump `github.com/volatiletech/strmangle` from v0.0.1 to v0.0.2](https://github.com/volatiletech/strmangle/compare/v0.0.1...v0.0.2)
  - Minor: [Bump `google.golang.org/api` from v0.63.0 to v0.74.0](https://github.com/googleapis/google-api-go-client/compare/v0.63.0...v0.74.0)
  - Minor: [Bump `github.com/BurntSushi/toml` from v1.0.0 to v1.1.0](https://github.com/BurntSushi/toml/releases/tag/v1.1.0)
  - Bump `golang.org/x/crypto` from v0.0.0-20211215165025-cf75a172585e to v0.0.0-20220411220226-7b82a4e95df4
  - Bump `golang.org/x/sys` from v0.0.0-20211210111614-af8b64212486 to v0.0.0-20220412211240-33da011f77ad
- We now support overriding `ENV` variables during **local** development through a `.env.local` dotenv file.
  - This does not require a development container restart.
  - We override the env within the app process through `config.DefaultServiceConfigFromEnv()`, so this does not mess with the actual container ENV.
  - See `.env.local.sample` for further instructions to use this.
  - Note that `.env.local` is **NEVER automatically** applied during **test runs**. If you really need that, use the specialized `test.DotEnvLoadLocalOrSkipTest` helper before loading up your server within that very test! This ensures that this test is automatically skipped if the `.env.local` file is no longer available.
- VSCode windows closes now explicitly stop Docker containers via [`shutdownAction: "stopCompose"`](https://code.visualstudio.com/docs/remote/devcontainerjson-reference) within `.devcontainer.json`.
  - Use `./docker-helper --halt` or other `docker` or `docker-compose` management commands to do this explicitly instead.
- Drone CI specific (minor): Fix multiline ENV variables were messing up our `.hostenv` for `docker run` command testing of the final image.

## 2022-03-28

- Merged [#165: Allow use of db.join* methods more than once](https://github.com/allaboutapps/go-starter/pull/165), thx [danut007ro](https://github.com/danut007ro).
- Merged [#169: Switch to standalone cobra-cli dependency](https://github.com/allaboutapps/go-starter/pull/169), thx [liggitt](https://github.com/liggitt) (requires `./docker-helper.sh --rebuild`).
  - [`github.com/spf13/cobra@v1.4.0`](https://github.com/spf13/cobra/releases/tag/v1.4.0) split into `cobra` (the lib) and [`github.com/spf13/cobra-cli`](https://github.com/spf13/cobra-cli/releases) (the generator / scaffolding tool)
  - We'll now depend on `cobra-cli` directly in our `Dockerfile`, while the core `cobra` dependency stays unchanged within our `go.mod`.
  - Bumps [`github.com/spf13/cobra`](https://github.com/spf13/cobra) from v1.3.0 to [v1.4.0](https://github.com/spf13/cobra/releases/tag/v1.4.0)
- Fixed `test.ApplyMigrations` when combined with the import SQL dump mechanics in the testing context.
  - Previously, we did still use the default [sql-migrate](https://github.com/rubenv/sql-migrate) `gorp_migrations` table to track applied migrations in our test databases, not our typical `migrations` table used everywhere else.
  - This especially lead to problems when importing (production / live) SQL dumps via `test.WithTestDatabaseFromDump*`, `test.WithTestServerFromDump*` or `test.WithTestServerConfigurableFromDump` as our implementation tried to apply **all migrations** every time, regardless if a partial migration set was already applied previously (as the already applied migrations were not tracked within the `migrations` table (but within `gorp_migrations`) we did not notice).
  - We now initialize this pipeline correctly in the test context (similar to our usage within `cmd/db_migrate.go` or `app db migrate`) and explicitly set these globals through `config.DatabaseMigrationTable` and `config.DatabaseMigrationFolder`.
  - If you encounter problems after the upgrade, please execute `make sql-drop-all` in your local environment to reset the IntegreSQL test databases, then run `make sql-reset && make sql-spec-reset && make sql-spec-migrate && make all` to rebuild and test.


## 2022-02-28

### Changed

- **BREAKING** Username format change in auth handlers
  - Added the `util.ToUsernameFormat` helper function, which will **lowercase** and **trim whitespaces**. We use it to format usernames in the login, register, and forgot-password handlers.
    - This prevents user duplication (e.g. two accounts registered with the same email address with different casing) and
    - cases where users would inadvertently register with specific casing or a trailing whitespace after their username, and subsequently struggle to log into their account.
  - **This effectively locks existing users whose username contains uppercase characters and/or whitespaces out of their accounts.**
    - Before rolling out this change, check whether any existing users are affected and migrate their usernames to a format that is compatible with this change.
    - Be aware that this may cause conflicts in regard to the uniqueness constraint of usernames and therefore need to be resolved manually, which is why we are not including a database migration to automatically migrate existing usernames to the new format.
  - For more information and a possible manual database migration flow please see this special WIKI page: https://github.com/allaboutapps/go-starter/wiki/2022-02-28

## 2022-02-03

### Changed

- Changed order of make targets in the `make swagger` pipeline. `make swagger-lint-ref-siblings` will now run after `make swagger-concat`, always linting the current version of our swagger file. This helps avoid errors regarding an invalid `swagger.yml` when resolving merge conflicts as those are often resolved by running `make swagger` and generating a fresh `swagger.yml`.

## 2022-02-02

### Changed

- Upgrades to [go-swagger](https://github.com/go-swagger/go-swagger) from to v0.26.1 to [v0.29.0](https://github.com/go-swagger/go-swagger/releases/tag/v0.29.0) (development stage only, requires `./docker-helper.sh --rebuild`). Includes the following `go.mod` upgrades:
  - [github.com/go-openapi/runtime](https://github.com/go-openapi/runtime) from v0.19.31 to v0.22.0
  - [github.com/go-openapi/strfmt](https://github.com/go-openapi/strfmt) from v0.20.2 to v0.21.1
  - [github.com/go-openapi/validate](https://github.com/go-openapi/validate) from v0.20.2 to v0.20.3
  - [github.com/go-openapi/errors](https://github.com/go-openapi/errors) from v0.20.1 to v0.20.2
  - [github.com/go-openapi/swag](https://github.com/go-openapi/swag) from v0.19.15 to v0.21.1
- Adds `yq` ([yq: a lightweight and portable command-line YAML processor](https://github.com/mikefarah/yq)) to our `Dockerfile` (development stage only, requires `./docker-helper.sh --rebuild`).
- Adds `make swagger-lint-ref-siblings` which is now executed as part of the `make build` (and `make swagger`) pipeline.
  - Any sibling elements of a Swagger `$ref` are ignored.
  - We have seen several misuses of `$ref` in our projects causing weird merge/flatten behaviors, thus we now lint for this case explicitly.
  - Having `$ref` and sibling elements (e.g. `required`, `example`, ...) is unsupported by [OpenAPI v2: $ref and Sibling Elements](https://swagger.io/docs/specification/using-ref/) itself and the [JSON Reference specification](https://datatracker.ietf.org/doc/html/rfc3986) itself.
  - To mitigate these errors, either expand the referenced element (fully remove `$ref`) or create a new element including your custom siblings elements and `$ref` this new one.
- Fix schema visualization generation guide in `docs/schemacrawler/README.md`

## 2021-12-14

### Changed
- Add i18n service wrapping `go-i18n` package by nicksnyder.
  - Allows parsing of Accept-Language header and language string.
  - Support for templating using go templating language in message values.
  - Support for [CLDR plural keys](https://cldr.unicode.org/index/cldr-spec/plural-rules)
  - Added environment variables to configure i18n service
    - `SERVER_I18N_DEFAULT_LANGUAGE` - set default language for i18n service
    - `SERVER_I18N_BUNDLE_DIR_ABS` - set directory of i81n messages, available languages are automatically configured by the files present in the folder

## 2021-11-29

### Changed

- The `integresql` service previously bound its port (`5000`) to the host machine. As this conflicts with newer macOS releases and is not necessary for the development workflow, the port is now only exposed to the linked services.

## 2021-10-22

### Changed

- Fixes minor `Makefile` typos.
- New go-starter releases are now git tagged (starting from the previous release `go-starter-2021-10-19` onwards). See [FAQ: What's the process of a new go-starter release?](https://github.com/allaboutapps/go-starter/wiki/FAQ#whats-the-process-of-a-new-go-starter-release)
- You may now specify a **specific** tag/branch/commit from the upstream [go-starter](https://github.com/allaboutapps/go-starter) project while running `make git-fetch-go-starter`, `make git-compare-go-starter` and `make git-merge-go-starter`. This will especially come in handy if you want to do a multi-phased merge (for projects that haven't been updated in a long time):
  - Merge with the latest: `make git-merge-go-starter`
  - Merge with a specific tag, e.g. the tag [`go-starter-2021-10-19`](https://github.com/allaboutapps/go-starter/releases/tag/go-starter-2021-10-19): `GIT_GO_STARTER_TARGET=go-starter-2021-10-19 make git-merge-go-starter`
  - Merge with a specific branch, e.g. the branch [`mr/housekeeping`](https://github.com/allaboutapps/go-starter/tree/mr/housekeeping): `GIT_GO_STARTER_TARGET=go-starter/mr/housekeeping make git-merge-go-starter` (heads up! it's `go-starter/<branchname>`)
  - Merge with a specific commit, e.g. the commit [`e85bedb94c3562602bc23d2bfd09fca3b13d1e02`](https://github.com/allaboutapps/go-starter/commit/e85bedb94c3562602bc23d2bfd09fca3b13d1e02): `GIT_GO_STARTER_TARGET=e85bedb94c3562602bc23d2bfd09fca3b13d1e02 make git-merge-go-starter`
- The primary GitHub Action pipeline `.github/workflows/build-test.yml` has been synced to include most validation tasks from our internal `.drone.yml` pipeline. Furthermore:
  - Avoid `Build & Test` GitHub Action running twice (on `push` and on `pull_request`).
  - Add trivy scan to our base Build & Test pipeline (as we know also build and test the `app` target docker image).
  - Our GitHub Action pipeline will no longer attempt to cache the previously built Docker images by other pipelines, as extracting/restoring from cache (docker buildx) typically takes **longer** than fully rebuilding the whole image. We will reinvestigate caching mechanisms in the future if GitHub Actions provides a speedier and official integration for Docker images.

## 2021-10-19

### Changed

- **BREAKING** Upgrades to [Go 1.17.1](https://golang.org/doc/go1.17) `golang:1.17.1-buster`
  - Switch to `//go:build <tag>` from `// +build <tag>`.
  - Migrates `go.mod` via `go mod tidy -go=1.17` (pruned module graphs).
  - Do the following to upgrade:
    1. `make git-merge-go-starter`
    2. `./docker-helper --rebuild`
    3. Manually remove the new **second** `require` block (with all the `// indirect` modules) within your `go.mod`
    4. Execute `go mod tidy -go=1.17` once so the **second** `require` block appears again.
    5. Find `// +build <tag>` and replace it with `//go:build <tag>`.
    6. `make all`.
    7. Recheck your `go.mod` that the newly added `// indirect` transitive dependencies are the proper version as you were previously using (e.g. via the output from `make get-licenses` and `make get-embedded-modules`). Feel free to move any `// indirect` tagged dependencies in your **first** `require` block to the **second** block. This is where they should live.
- **BREAKING** You now need to take special care when it comes to parsing **semicolons** (`;`) in **query strings** via `net/url` and `net/http` from Go >1.17!
  - Anything before the semicolon will now be stripped. e.g. `example?a=1;b=2&c=3` would have returned `map[a:[1] b:[2] c:[3]]`, while now it returns `map[c:[3]]`
  - See [Go 1.17 URL query parsing](https://golang.org/doc/go1.17#semicolons).
  - You may need to manually migrate your handlers/tests regarding this new default handling.

## 2021-09-27

### Changed

- Added `make test-update-golden` for easily refreshing **all** golden files / snapshot tests (`y + ENTER` confirmation).
- Upgrades [golangci-lint](https://github.com/golangci/golangci-lint) from `v1.41.1` to [`v1.42.1`](https://github.com/golangci/golangci-lint/releases/tag/v1.42.1) (for reference [`v1.42.0`](https://github.com/golangci/golangci-lint/releases/tag/v1.42.0)).
- Bump github.com/go-openapi/strfmt from [0.20.1 to 0.20.2](https://github.com/go-openapi/strfmt/compare/v0.20.1...v0.20.2)
- Bump github.com/go-openapi/errors from [0.20.0 to 0.20.1](https://github.com/go-openapi/errors/compare/v0.20.0...v0.20.1)
- Bump github.com/go-openapi/runtime from [0.19.29 to 0.19.31](https://github.com/go-openapi/runtime/compare/v0.19.29...v0.19.31)
- Bump github.com/rs/zerolog from [1.23.0 to 1.25.0](https://github.com/rs/zerolog/compare/v1.23.0...v1.25.0)
- Bump google.golang.org/api from [0.52.0 to 0.57.0](https://github.com/allaboutapps/go-starter/pull/124)
- Bump github.com/lib/pq from [v1.10.2 to v1.10.3](https://github.com/lib/pq/releases/tag/v1.10.3)
- Bump github.com/spf13/viper from [1.8.1 to v1.9.0](https://github.com/spf13/viper/releases/tag/v1.9.0)
- Bump github.com/labstack/echo from [4.5.0 to v4.6.1](https://github.com/labstack/echo/compare/v4.5.0...v4.6.1)
- Update golang.org/x/crypto and golang.org/x/sys

## 2021-08-17

### Changed

- **Hotfix**: We will pin the `Dockerfile` development and builder stage to `golang:1.16.7-buster` (+ `-buster`) for now, as currently the [new debian bullseye release within the go official docker images](https://github.com/docker-library/golang/commit/48a7371ed6055a97a10adb0b75756192ad5f1c97) breaks some tooling. The upgrade to debian bullseye and Go 1.17 will happen ~simultaneously~ **separately** within go-starter in the following weeks.

## 2021-08-16

### Changed

- remove ioutil (https://golang.org/doc/go1.16#ioutil)

## 2021-08-06

### Changed

- Bump golang from 1.16.6 to [1.16.7](https://github.com/golang/go/issues?q=milestone%3AGo1.16.7+label%3ACherryPickApproved) (requires `./docker-helper.sh --rebuild`).
- Adds `util.GetEnvAsStringArrTrimmed` and minor `util` test coverage upgrades.

## 2021-08-04

### Changed

- `README.md` badges for go-starter.
- Fix some misspellings of English words within `internal/test/*.go` comments.
- Upgrades
  - Bump `github.com/labstack/echo/v4` from 4.4.0 to [4.5.0](https://github.com/labstack/echo/blob/master/CHANGELOG.md#v450---2021-08-01):
    - Switch from `github.com/dgrijalva/jwt-go` to [`github.com/golang-jwt/jwt`](https://github.com/golang-jwt/jwt) to mitigate [CVE-2020-26160](https://nvd.nist.gov/vuln/detail/CVE-2020-26160).
    - Note that it might take some time until the former dep fully leaves our dependency graph, as it is also a transitive dependency of various versions of [`github.com/spf13/viper`](https://github.com/spf13/viper/issues/997).
    - However, even though this functionality was never used by go-starter, this change fixes an important part: The original `github.com/dgrijalva/jwt-go` is no longer included in the **final `app` binary**, it is fully replaced by `github.com/golang-jwt/jwt`.
    - Our `.trivyignore` still excludes [CVE-2020-26160](https://nvd.nist.gov/vuln/detail/CVE-2020-26160) as trivy cannot skip checking transitive dependencies.
    - **Breaking**: If you have actually directly depended upon `github.com/dgrijalva/jwt-go`, please switch to `github.com/golang-jwt/jwt` via the following command: `find -type f -name "*.go" -exec sed -i "s/dgrijalva\/jwt-go/golang-jwt\/jwt/g" {} \;`

## 2021-07-30

### Changed

- Upgrades:
  - Bump golang from 1.16.5 to [1.16.6](https://groups.google.com/g/golang-announce/c/n9FxMelZGAQ)
  - Bump github.com/labstack/echo/v4 from 4.3.0 to [4.4.0](https://github.com/labstack/echo/blob/master/CHANGELOG.md) (adds `binder.BindHeaders` support, not affecting our goswagger `runtime.Validatable` bind helpers)
  - Bump github.com/gabriel-vasile/mimetype from 1.3.0 to [1.3.1](https://github.com/gabriel-vasile/mimetype/releases/tag/v1.3.1)
  - Bump github.com/spf13/cobra from 1.1.3 to [1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.1) (and see all the big completion upgrades in [1.2.0](https://github.com/spf13/cobra/releases/tag/v1.2.0))
  - Bump google.golang.org/api from 0.49.0 to [0.52.0](https://github.com/allaboutapps/go-starter/pull/106)
  - Bump gotestsum to [1.7.0](https://github.com/gotestyourself/gotestsum/releases/tag/v1.7.0) (adds handy keybindings while you are in `make watch-tests` mode, see [While in watch mode, pressing some keys will perform an action](https://github.com/gotestyourself/gotestsum#run-tests-when-a-file-is-saved))
  - Bump watchexec to [1.17.0](https://github.com/watchexec/watchexec/releases/tag/cli-v1.17.0)
  - Bump golang.org/x/crypto to `v0.0.0-20210711020723-a769d52b0f97`

## 2021-07-29

### Changed

- Fixed `Makefile` has disregarded `pipefail`s in executed targets (e.g. `make sql-spec-migrate` previously returned exit code `0` even if there were migration errors as its output was piped internally). We now set `-cEeuo pipefail` for make's shell args, preventing these issues.

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
  - Pin to `actions/checkout@v2.3.4`.
  - Remove unnecessary `git checkout HEAD^2` in CodeQL step (Code Scanning recommends analyzing the merge commit for best results).
  - Limit trivy and codeQL actions to `push` against `master` and `pull_request` against `master` to overcome read-only access workflow errors.

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

- Moved `/api/main.yml` to `/api/config/main.yml` to overcome path resolve issues (`../definitions`) with the VSCode [42crunch.vscode-openapi](https://github.com/42Crunch/vscode-openapi) extension (auto-included in our devContainer) and our go-swagger concat behaviour.
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
  - `make lint`: Runs golangci-lint and make check-\*.
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
- Seeding: Switch to `db|dbUtil.WithTransaction` instead of manually managing the db transaction. _Note_: We will enforce using `WithTransaction` instead of manually managing the life-cycle of db transactions through a custom linter in an upcoming change. It's way safer and manually managing db transactions only makes sense in very very special cases (where you will be able to opt-out via linter excludes). Also see [What's `WithTransaction`, shouldn't I use `db.BeginTx`, `db.Commit`, and `db.Rollback`?](https://github.com/allaboutapps/go-starter/wiki/FAQ#whats-withtransaction-shouldnt-i-use-dbbegintx-dbcommit-and-dbrollback).

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
