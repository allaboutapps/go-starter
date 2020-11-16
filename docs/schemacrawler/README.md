# `/docs/schemacrawler`

To locally (re-)generate a schemacrawler diagramm, execute any of the following commands from your **host** machine.

```bash
# Note that the project must be already running within docker-compose (and the "spec" database should already be migrated via "make sql" or "make all").
# First find out under which docker network the "allaboutapps.dev/aw/go-starter" project is available (as started via ./docker-helper.sh --up).
# Typically it's "<dir_name>_default".
docker network ls
# [...]
# go-starter_default

# Ensure you are within the /docs/schemacrawler directory
cd docs/schemacrawler
pwd
# [...]/docs/schemacrawler

# Generate a png (exchange --network="..." with your docker network before executing this command)
docker run --network=go-starter_default -v $(pwd):/home/schcrwlr/share -v $(pwd)/schemacrawler.config.properties:/opt/schemacrawler/config/schemacrawler.config.properties --entrypoint=/opt/schemacrawler/schemacrawler.sh schemacrawler/schemacrawler --server=postgresql --host=postgres --port=5432 --database=spec --schemas=public --user=dbuser --password=dbpass --info-level=standard --command=schema --portable-names --title "allaboutapps.dev/aw/go-starter" --output-format=png --output-file=/home/schcrwlr/share/schema.png

# Generate a pdf (exchange --network="..." with your docker network before executing this command)
docker run --network=go-starter_default -v $(pwd):/home/schcrwlr/share -v $(pwd)/schemacrawler.config.properties:/opt/schemacrawler/config/schemacrawler.config.properties --entrypoint=/opt/schemacrawler/schemacrawler.sh schemacrawler/schemacrawler --server=postgresql --host=postgres --port=5432 --database=spec --schemas=public --user=dbuser --password=dbpass --info-level=standard --command=schema --portable-names --title "allaboutapps.dev/aw/go-starter" --output-format=pdf --output-file=/home/schcrwlr/share/schema.pdf

# Feel free to override schemacrawler configuration settings in "./schemacrawler.config.properties".
```

For further information see:
- [SchemaCrawler Database Diagramming](https://www.schemacrawler.com/diagramming.html) (intro to most diagramming options)
- [Docker Image for SchemaCrawler](https://www.schemacrawler.com/docker-image.html) (about running schemacrawler in Docker)
- [DockerHub `schemacrawler/schemacrawler`](https://hub.docker.com/r/schemacrawler/schemacrawler/) (available version of this Docker image)