### Quickstart

> Requires docker and docker-compose install locally

```bash

./docker-helper.sh --up

# You should now have a docker shell...

# Init install/cache dependencies and install tools to bin
make init

# Migrate up your local database
sql-migrate up

# Building (generate, format, build, vet)
make

# Execute tests
make test

```