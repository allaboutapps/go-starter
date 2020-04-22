package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var (
	// TODO: env: how to get to the project root?
	handlersPackages = []string{
		"/app/api/handlers/auth",
		"/app/api/handlers/users",
	}

	// https://developer.mozilla.org/de/docs/Web/HTTP/Methods
	// any func that start with this name are considered
	// TODO: also check fn signature
	methodPrefixes = []string{
		"Get", "Head", "Patch", "Post", "Put", "Delete",
	}
)

// get all functions in above handler packages
// that match Get*, Put*, Post*, Patch*, Delete*
func main() {
	set := token.NewFileSet()

	for _, subPackage := range handlersPackages {

		packs, err := parser.ParseDir(set, subPackage, nil, 0)

		if err != nil {
			fmt.Println("Failed to parse package:", err)
			os.Exit(1)
		}

		funcs := []*ast.FuncDecl{}
		for _, pack := range packs {
			for _, f := range pack.Files {
				for _, d := range f.Decls {
					if fn, isFn := d.(*ast.FuncDecl); isFn {

						fnName := fn.Name.String()

						for _, prefix := range methodPrefixes {
							if strings.HasPrefix(fnName, prefix) {
								funcs = append(funcs, fn)
							}
						}
					}
				}
			}
		}

		// print out
		for _, function := range funcs {
			fmt.Println(subPackage, function.Name)
		}
	}
}
