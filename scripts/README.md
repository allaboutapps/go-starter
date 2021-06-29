# `/scripts`

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple.

https://github.com/golang-standards/project-layout/tree/master/scripts

Examples:

* https://github.com/kubernetes/helm/tree/master/scripts
* https://github.com/cockroachdb/cockroach/tree/master/scripts
* https://github.com/hashicorp/terraform/tree/master/scripts

Please note that this scripts are not available in a final product. Head to `../cmd` if you need to execute your script in live environments.

The `gsdev` cli util executes this `scripts/main.go` file here and also describes all available commands available while developing a project locally. `gsdev` is made available during the `Dockerfile`'s development stage.

### `/scripts/cmd/*.go`

`func`s may define shared logic used in `/scripts/internal/**/*.go`, the actual usable commands are defined within `/scripts/internal`.

### `// +build scripts`

Any `*.go` file in all subdirectories of `/scripts/**` should specify `// +build scripts` to signal that those files are not part of of our final product. To execute any script that has this build tag, you need to specify `-tags scripts`, otherwise you will run into an error like the following (also see our `Makefile` for a reference):

```bash
# Works
go run -tags scripts scripts/main.go
# go-starter development scripts
# Utility commands while developing go-starter based projects.

# Works (same as above)
gsdev
# go-starter development scripts
# Utility commands while developing go-starter based projects.

# Misses build tag "scripts"
go run scripts/main.go
# package command-line-arguments
# 	imports allaboutapps.dev/aw/go-starter/scripts/cmd: build constraints exclude all Go files in /app/scripts/cmd
```
