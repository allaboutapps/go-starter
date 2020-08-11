package util

import (
	"io/ioutil"
	"log"

	"github.com/rogpeppe/go-internal/modfile"
)

// https://stackoverflow.com/questions/53183356/api-to-get-the-module-name
// https://github.com/rogpeppe/go-internal
func GetModuleName(absolutePathToGoMod string) (string, error) {
	dat, err := ioutil.ReadFile(absolutePathToGoMod)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	mod := modfile.ModulePath(dat)

	return mod, nil
}
