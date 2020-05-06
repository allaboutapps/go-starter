# go-starter

> This project tries to adhere to the layout defined in [golang-standard/project-layout](https://github.com/golang-standards/project-layout)

### Development setup

Requires the following local setup for development:

- [Docker CE](https://docs.docker.com/install/) (19.03 or above)
- [Docker Compose](https://docs.docker.com/compose/install/) (1.25 or above)

The project makes use of the [devcontainer functionality](https://code.visualstudio.com/docs/remote/containers) provided by [Visual Studio Code](https://code.visualstudio.com/), thus a local installation of a Go compiler is *no longer* required when using this IDE.

### Development quickstart

> Requires docker and docker-compose installed locally

```bash

# $local
./docker-helper.sh --up

# You should now have a docker shell...
# development@XXXXXXXXX:/app$

# If you have forked this project, easily change the go project module name
# allaboutapps.dev/<GIT_PROJECT>/<GIT_REPO>
make set-module-name

# Shortcut for make init, make build, make info and make test
make all

# Init install/cache dependencies and install tools to bin
make init

# Rebuild only after changes to files (generate, format, build, vet)
make

# Execute all tests
make test

# Migrate up the development database
sql-migrate up

# Start the local-built server
apiserver

```

Regarding [Visual Studio Code](https://code.visualstudio.com/): Always develop *inside* the running `development` docker container, by attaching to this container.

Run CMD+SHIFT+P `Go: Install/Update Tools` after starting vscode to autoinstall all golang vscode dependencies, then **reload your window**.
