# go-starter

**go-starter** is an opinionated *production-ready* RESTful JSON backend template written in [Go](https://golang.org/), highly integrated with [VSCode DevContainers](https://code.visualstudio.com/docs/remote/containers) by [allaboutapps](https://allaboutapps.at/).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/allaboutapps/go-starter/blob/master/LICENSE)
[![Build and Test](https://github.com/allaboutapps/go-starter/actions/workflows/build-test.yml/badge.svg)](https://github.com/allaboutapps/go-starter/actions)
[![codecov](https://codecov.io/gh/allaboutapps/go-starter/branch/master/graph/badge.svg?token=220E44857K)](https://codecov.io/gh/allaboutapps/go-starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/allaboutapps/go-starter)](https://goreportcard.com/report/github.com/allaboutapps/go-starter)
[![Swagger Validator](https://img.shields.io/swagger/valid/3.0?specUrl=https%3A%2F%2Fraw.githubusercontent.com%2Fallaboutapps%2Fgo-starter%2Fmaster%2Fapi%2Fswagger.yml)](https://go-starter.allaboutapps.at/documentation/)
![GitHub contributors](https://img.shields.io/github/contributors/allaboutapps/go-starter)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

![go starter overview](https://public.allaboutapps.at/go-starter-wiki/go-starter-main-overview.png)

Demo: **[https://go-starter.allaboutapps.at](https://go-starter.allaboutapps.at)**
FAQ: **[https://github.com/allaboutapps/go-starter/wiki/FAQ](https://github.com/allaboutapps/go-starter/wiki/FAQ)**

## Table of Contents

- [go-starter](#go-starter)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Usage](#usage)
    - [Demo](#demo)
    - [Requirements](#requirements)
    - [Quickstart](#quickstart)
    - [Merge with the go-starter template repository to get future updates](#merge-with-the-go-starter-template-repository-to-get-future-updates)
    - [Set project module name for your new project](#set-project-module-name-for-your-new-project)
    - [Visual Studio Code](#visual-studio-code)
    - [Building and testing](#building-and-testing)
    - [Running](#running)
    - [Uninstall](#uninstall)
  - [Additional resources](#additional-resources)
  - [Contributing](#contributing)
  - [Current maintainers](#current-maintainers)
    - [Previous maintainers](#previous-maintainers)
  - [License](#license)

## Features

- Full local golang service development environment using [Docker Compose](https://docs.docker.com/compose/install/) and [VSCode devcontainers](https://code.visualstudio.com/docs/remote/containers) that just works with Linux, MacOS and Windows.
- Adheres to the project layout defined in [golang-standard/project-layout](https://github.com/golang-standards/project-layout).
- Provides database migration ([sql-migrate](https://github.com/rubenv/sql-migrate)) and models generation ([SQLBoiler](https://github.com/volatiletech/sqlboiler)) workflows for [PostgreSQL](https://www.postgresql.org/) databases.
- Integrates [IntegreSQL](https://github.com/allaboutapps/integresql) for fast, concurrent and isolated integration testing with real PostgreSQL databases.
- Auto-installs our recommended VSCode extensions for golang development.
- Integrates [go-swagger](https://github.com/go-swagger/go-swagger) for compile-time generation of `swagger.yml`, structs and request/response validation functions.
- Integrates [MailHog](https://github.com/mailhog/MailHog) for easy SMTP-based email testing.
- Integrates [SwaggerUI](https://github.com/swagger-api/swagger-ui) for live-previewing your Swagger v2 schema.
- Integrates [pgFormatter](https://github.com/darold/pgFormatter) and [vscode-pgFormatter](https://marketplace.visualstudio.com/items?itemName=bradymholt.pgformatter) for SQL formatting.
- Comes with fully implemented `auth` package, an OAuth2 RESTful JSON API ready to be extended according to your requirements.
- Implements [OAuth 2.0 Bearer Tokens](https://tools.ietf.org/html/rfc6750) and password authentication using [argon2id](https://godoc.org/github.com/alexedwards/argon2id) hashes.
- Comes with a tested mock and [FCM](https://firebase.google.com/docs/cloud-messaging) provider for sending push notifications and storing push tokens.
- CLI layer provided by [spf13/cobra](https://github.com/spf13/cobra). It's exceptionally easy to [add additional sub-commands via `cobra-cli`](https://github.com/spf13/cobra-cli/blob/main/README.md#add-commands-to-a-project).
- Comes with an initial [PostgreSQL](https://www.postgresql.org/) database structure (see [/migrations](https://github.com/allaboutapps/go-starter/tree/master/migrations)), covering:
  - auth tokens (access-, refresh-, password-reset-tokens),
  - a generic auth-related `user` model
  - an app-specific bare-bones `app_user_profile` model,
  - push notification tokens and
  - a health check sequence (for performing writeable checks).
- API endpoints and CLI for liveness (`/-/healthy`) and readiness (`/-/ready`) probes
- Parallel jobs optimized `Makefile` and various convenience scripts (see all targets and its description via `make help`). A full rebuild only takes seconds.
- Multi-staged `Dockerfile` (`development` -> `builder` -> `app`).

## Usage

> Please find more detailed information regarding the history, usage and other *whys?* of this project in our **[FAQ](https://github.com/allaboutapps/go-starter/wiki/FAQ)**.

### Demo

A demo go-starter service is deployed at **[https://go-starter.allaboutapps.at](https://go-starter.allaboutapps.at)** for you to play around with.

Please visit our [FAQ](https://github.com/allaboutapps/go-starter/wiki/FAQ#what-are-the-limitations-of-your-demo-environment) to find out more about the limitations of this demo environment.

### Requirements

Requires the following local setup for development:

- [Docker CE](https://docs.docker.com/install/) (19.03 or above)
- [Docker Compose](https://docs.docker.com/compose/install/) (1.25 or above)
- [VSCode Extension: Remote - Containers](https://code.visualstudio.com/docs/remote/containers) (`ms-vscode-remote.remote-containers`)

This project makes use of the [Remote - Containers extension](https://code.visualstudio.com/docs/remote/containers) provided by [Visual Studio Code](https://code.visualstudio.com/). A local installation of the Go tool-chain is **no longer required** when using this setup.

Please refer to the [official installation guide](https://code.visualstudio.com/docs/remote/containers) how this works for your host OS and head to our [FAQ: How does our VSCode setup work?](https://github.com/allaboutapps/go-starter/wiki/FAQ#how-does-our-vscode-setup-work) if you encounter issues.

### Quickstart

Create a new git repository through the GitHub template repository feature ([use this template](https://github.com/allaboutapps/go-starter/generate)). You will then start with a **single initial commit** in your own repository.

```bash
# Clone your new repository, cd into it, then easily start the docker-compose dev environment through our helper
./docker-helper.sh --up
```

You should be inside the 'service' docker container with a bash shell.

```bash
development@94242c61cf2b:/app$ # inside your container...

# Shortcut for make init, make build, make info and make test
make all

# Print all available make targets
make help
```

### Merge with the go-starter template repository to get future updates

> These steps are **not** necessary if you have a *"real"* fork.

If your new project is generated from a template project (you have a **single commit**), you want to run the following command immediately and **before** applying any changes. Otherwise you won't be able to easily merge upstream go-starter changes into your own repository (see [GitHub Template Repositories](https://docs.github.com/en/free-pro-team@latest/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template), [Refusing to merge unrelated histories](https://www.educative.io/edpresso/the-fatal-refusing-to-merge-unrelated-histories-git-error) and [FAQ: I want to compare or update my project/fork to the latest go-starter master](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-want-to-compare-or-update-my-projectfork-to-the-latest-go-starter-master)).

```bash
make git-merge-go-starter
# Attempting to execute 'git merge --no-commit --no-ff go-starter/master' into your current HEAD.
# Are you sure? [y/N]y
# git merge --no-commit --no-ff --allow-unrelated-histories go-starter/master

git commit -m "Initial merge of unrelated go-starter template history"
```

### Set project module name for your new project

To replace all occurrences of `allaboutapps.dev/aw/go-stater` (our internal module name of this project) with your desired projects' module name, do the following:

```bash
development@94242c61cf2b:/app$ # inside your container...

# Set a new go project module name.
make set-module-name
# allaboutapps.dev/<GIT_PROJECT>/<GIT_REPO> (internal only)
# github.com/<USER>/<PROJECT>
# e.g. github.com/majodev/my-service
```

The above command writes your new go module name to `tmp/.modulename`, `go.mod`. It actually sets it everywhere in `**/*` - thus this step is typically only required **once**. If you need to merge changes from the upstream go-starter later, we may want to run `make force-module-name` to set your own go module name everywhere again (especially relevant for new files / import paths). See our [FAQ](https://github.com/allaboutapps/go-starter/wiki/FAQ#i-want-to-compare-or-update-my-projectfork-to-the-latest-go-starter-master) for more information about this update flow.

Optionally you may want to move the original `README.md` and `LICENSE` away:

```bash
development@94242c61cf2b:/app$ # inside your container...

# Optionally you may want to move our LICENSE and README.md away.
mv README.md README-go-starter.md
mv LICENSE LICENSE-go-starter

# Optionally create a new README.md for your project.
make get-module-name > README.md
```

### Visual Studio Code

> If you are new to VSCode Remote - Containers feature, see our [FAQ: How does our VSCode setup work?](https://github.com/allaboutapps/go-starter/wiki/FAQ#how-does-our-vscode-setup-work).

Run `CMD+SHIFT+P` `Go: Install/Update Tools` **after** attaching to the container with VSCode to auto-install all golang related vscode extensions.


### Building and testing

Other useful commands while developing your service:

```bash
development@94242c61cf2b:/app$ # inside your container...

# Print all available make targets
make help

# Shortcut for make init, make build, make info and make test
make all

# Init install/cache dependencies and install tools to bin
make init

# Rebuild only after changes to files (generate, format, build, lint)
make

# Execute all tests
make test
```

### Running

To run the service locally you may:

```bash
development@94242c61cf2b:/app$ # inside your development container...

# First ensure you have a fresh `app` executable available
make build

# Check if all requirements for becoming are met (db is available, mnt path is writeable)
app probe readiness -v

# Migrate up the database
app db migrate

# Seed the database (if you have any fixtures defined in `/internal/data/fixtures.go`)
app db seed

# Start the locally-built server
app server

# Now available at http://127.0.0.1:8080

# You may also run all the above commands in a single command
app server --probe --migrate --seed # or `app server -pms`
```

### Uninstall

Simply run `./docker-helper --destroy` in your working directory (on your host machine) to wipe all docker related traces of this project (and its volumes!).

## Additional resources

* **Please visit our [FAQ](https://github.com/allaboutapps/go-starter/wiki/FAQ)**.
* [Random Training Material](https://github.com/allaboutapps/go-starter/wiki/Random-training-material)

## Contributing

Pull requests are welcome. For major changes, please [open an issue](https://github.com/allaboutapps/go-starter/issues/new) first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Current maintainers

- [Anna - @anjankow](https://github.com/anjankow)
- [Michael Farkas - @farkmi](https://github.com/farkmi)
- [Mario Ranftl - @majodev](https://github.com/majodev)
- [Manuel Wieser - @mwieser](https://github.com/mwieser)

### Previous maintainers

- [Dominic Aschauer - @eldelto](https://github.com/eldelto)
- [Nick Müller - @MorpheusXAUT](https://github.com/MorpheusXAUT)

## License

[MIT](LICENSE)
