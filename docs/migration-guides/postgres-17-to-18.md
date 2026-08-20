# Migration guide: Postgres 17 → 18

go-starter moved from `postgres:17.4-alpine` to `postgres:18.6-alpine`.

There are **two independent breaking changes** here. The second one is the dangerous
one, because nothing errors when you get it wrong.

## 1. The data directory is not compatible

A Postgres 18 server cannot read a Postgres 17 data directory. There is no in-place
upgrade — the server refuses to start.

## 2. The image moved `PGDATA` and its `VOLUME`

The official Postgres image changed its layout in 18:

| | 17 and earlier | 18+ |
| --- | --- | --- |
| `PGDATA` | `/var/lib/postgresql/data` | `/var/lib/postgresql/18/docker` |
| declared `VOLUME` | `/var/lib/postgresql/data` | `/var/lib/postgresql` |

go-starter's `docker-compose.yml` was updated accordingly:

```yaml
volumes:
  - pgvolume:/var/lib/postgresql        # was: pgvolume:/var/lib/postgresql/data
```

> **If you keep the old mount path, Postgres 18 starts up perfectly happily.**
> It just writes its data to an *anonymous* volume instead of yours, so the
> container comes up with an empty database and every write is discarded on the
> next `docker compose down`. There is no error message, no warning, and no failed
> health check — the only symptom is that your data is gone.

If your project changed the compose volume layout, or mounts the Postgres data
directory anywhere else (a bind mount, a second compose file, a CI service
definition), audit every one of those paths.

## Local development

```bash
docker compose down -v          # BEFORE rebuilding — drops the old PG17 volume
./docker-helper.sh --rebuild
docker compose up -d
make sql-reset && make sql-spec-reset && make sql-spec-migrate && make all
```

`down -v` deletes your local database. If it holds anything you care about, dump it
first with `pg_dump` while still on 17, then restore after the upgrade.

Verify the volume is actually wired up afterwards — write something, recreate the
containers, and check it survived:

```bash
docker compose exec postgres psql -U dbuser -d development \
  -c "create table probe(id int); insert into probe values (42);"
docker compose down && docker compose up -d postgres
docker compose exec postgres psql -U dbuser -d development -c "select * from probe;"
# must print 42 — if the table is missing, your volume mount path is wrong
docker compose exec postgres psql -U dbuser -d development -c "drop table probe;"
```

## Live / staging

Plan the major version upgrade explicitly — `pg_upgrade` or a dump/restore — and
schedule it with the deploy. Before merging:

- confirm your managed database provider actually offers Postgres 18
- confirm the `postgresql-client-<major>` in the `Dockerfile` matches the server
  major you will run against
- check any tooling that reads the data directory path directly (backup jobs,
  sidecars, volume snapshots)

## Staying on Postgres 17 for now

Supported, but keep the three settings consistent with each other:

- `postgres:17.11-alpine` in `docker-compose.yml`, `.drone.yml`,
  `.woodpecker/ci.yaml` and `.github/workflows/build-test.yml`
- `postgresql-client-17` in the `Dockerfile`
- the **old** `pgvolume:/var/lib/postgresql/data` mount

The mount path and the client version both follow the server major — do not mix
them.

## Release notes

Postgres 18 release notes: https://www.postgresql.org/docs/release/18.0/
