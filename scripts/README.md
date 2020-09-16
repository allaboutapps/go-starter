# `/scripts`

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple.

https://github.com/golang-standards/project-layout/tree/master/scripts

Examples:

* https://github.com/kubernetes/helm/tree/master/scripts
* https://github.com/cockroachdb/cockroach/tree/master/scripts
* https://github.com/hashicorp/terraform/tree/master/scripts

Please note that this scripts are not available in a final product. Use `/cmd` instead if you need to execute your script an a live environment.

### `/scripts/*.go`

`func`s may define shared logic used in `/scripts/**/*.go`.

### `// +build scripts`

Any `*.go` file in all subdirectories of `/scripts/**` should specify `// +build scripts` to signal that those files are not part of of our final product. To execute any script that has this build tag, you need to specify `-tags scripts`, otherwise you will run into an error like the following (also see our `Makefile` for a reference):

```bash

# Works
go run -tags scripts scripts/modulename/modulename.go
allaboutapps.dev/aw/go-starter

# Misses build tag "scripts"
go run scripts/modulename/modulename.go
package command-line-arguments
	imports allaboutapps.dev/aw/go-starter/scripts: build constraints exclude all Go files in /app/scripts
```
