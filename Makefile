### -----------------------
# --- Building
### -----------------------

# first is default task when running "make" without args
build: 
	@$(MAKE) --no-print-directory build-pre
	@$(MAKE) --no-print-directory go-format
	@$(MAKE) --no-print-directory build-post

# these recipies may execute in parallel
build-pre: sql-format sql-lint sql-check-migrations sqlboiler swagger go-generate 
build-post: go-build go-lint

go-format:
	go fmt

go-build: 
	go build -o bin/apiserver ./cmd/api

go-lint:
	golangci-lint run --fast

sqlboiler:
	sqlboiler --wipe --no-hooks psql

go-generate:
	go generate ./...

sql-format:
	@echo "make sql-format"
	@find ${PWD} -name ".*" -prune -o -type f -iname "*.sql" -print \
		| xargs -i pg_format {} -o {}

# check syntax via the real database
# https://stackoverflow.com/questions/8271606/postgresql-syntax-check-without-running-the-query
sql-lint:
	@echo "make sql-lint"
	@find ${PWD} -name ".*" -prune -o -type f -iname "*.sql" -print \
		| xargs -i sed '1s#^#DO $$SYNTAX_CHECK$$ BEGIN RETURN;#; $$aEND; $$SYNTAX_CHECK$$;' {} \
		| psql --quiet -v ON_ERROR_STOP=1

sql-check-migrations:
	@echo "make sql-check-migrations"
	@(grep -R " NULL" ./migrations/ | grep --invert "DEFAULT NULL" | grep --invert "NOT") && (echo "Unnecessary use of NULL keyword" && exit 1) || exit 0

swagger-gen-spec: 
	@echo "make swagger-gen-spec"
	@swagger generate spec \
		-i types/swagger/swagger.yml \
		-o types/swagger/swagger.json \
		--scan-models \
		-q

swagger-models:
	@echo "make swagger-models"
	@swagger generate model \
		--allow-template-override \
		--template-dir=types/swagger \
		--spec=types/swagger/swagger.json \
		--existing-models=allaboutapps.at/aw/go-mranftl-sample/types \
		--model-package=types \
		--all-definitions \
		-q

swagger-validate:
	@echo "make swagger-validate"
	@swagger validate types/swagger/swagger.json \
		--stop-on-error \
		-q

swagger-gen-server: swagger-validate swagger-models

swagger: 
	@$(MAKE) --no-print-directory swagger-gen-spec
	@$(MAKE) --no-print-directory swagger-gen-server

# https://github.com/golang/go/issues/24573
# w/o cache - see "go help testflag"
# use https://github.com/kyoh86/richgo to color
# note that these tests should not run verbose by default (e.g. use your IDE for this)
# TODO: add test shuffling/seeding when landed in go v1.15 (https://github.com/golang/go/issues/28592)
test:
	richgo test -cover -race -count=1 ./...

### -----------------------
# --- Initializing
### -----------------------

init:
	@$(MAKE) --no-print-directory modules
	@$(MAKE) --no-print-directory tools
	@$(MAKE) --no-print-directory tidy
	@go version

# cache go modules (locally into .pkg)
modules:
	go mod download

# https://marcofranssen.nl/manage-go-tools-via-go-modules/
tools:
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

tidy:
	go mod tidy

### -----------------------
# --- Helpers
### -----------------------

clean:
	rm -rf bin

reset:
	@echo "DROP & CREATE database:"
	@echo "  PGHOST=${PGHOST} PGDATABASE=${PGDATABASE}" PGUSER=${PGUSER}
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	psql -d postgres -c 'DROP DATABASE IF EXISTS "${PGDATABASE}";'
	psql -d postgres -c 'CREATE DATABASE "${PGDATABASE}" WITH OWNER ${PGUSER} TEMPLATE "template0"'

### -----------------------
# --- Special targets
### -----------------------

# https://www.gnu.org/software/make/manual/html_node/Special-Targets.html
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
# ignore matching file/make rule combinations in working-dir
.PHONY: test
