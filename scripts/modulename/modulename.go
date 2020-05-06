// +build ignore

// This program generates handlers.go. It can be invoked by running go generate ./...
package main

import (
	"fmt"
	"log"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/scripts/scriptsutil"
)

// https://blog.carlmjohnson.net/post/2016-11-27-how-to-use-go-generate/

var (
	PROJECT_ROOT  = util.GetProjectRootDir()
	PATH_MOD_FILE = PROJECT_ROOT + "/go.mod"
)

type ResolvedFunction struct {
	PackageName  string
	FunctionName string
}

// get all functions in above handler packages
// that match Get*, Put*, Post*, Patch*, Delete*
func main() {
	baseModuleName, err := scriptsutil.GetModuleName(PATH_MOD_FILE)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(baseModuleName)
}
