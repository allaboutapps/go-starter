build: prebuild format
	go build -o bin/app
	go vet

prebuild:
	go generate

format:
	go fmt
	find ${PWD} -name ".*" -prune -o -type f -iname "*.sql" -print | xargs -i pg_format {} -o {}

# https://github.com/golang/go/issues/24573
# w/o cache - see "go help testflag"
# use https://github.com/kyoh86/richgo to color
test:
	richgo test -cover -race -count=1 ./...

init: modules tools tidy
	@go version

# cache go modules (locally into .pkg)
modules:
	go mod download

# https://marcofranssen.nl/manage-go-tools-via-go-modules/
tools:
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

tidy:
	go mod tidy

clean:
	rm -rf bin