package test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
)

const snapshotDir string = "/snapshots"

var SnapshotDirPathAbs string
var defaultReplacer = func(s string) string {
	return s
}

var spewConfig = spew.ConfigState{
	Indent:                  "  ",
	SortKeys:                true, // maps should be spewed in a deterministic order
	DisablePointerAddresses: true, // don't spew the addresses of pointers
	DisableCapacities:       true, // don't spew capacities of collections
	SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
}

func init() {
	basePath := util.GetEnv("SERVER_PATHS_TEST_ASSETS_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/test/testdata"))
	SnapshotDirPathAbs = filepath.Join(basePath, snapshotDir)
}

// vastly inspired by https://github.com/bradleyjkemp/cupaloy
// main reason for self implementation is the replacer function, to replace generated values (IDs)
func SnapshotWithReplacer(t TestingT, update bool, replacer func(s string) string, data ...interface{}) {
	t.Helper()
	err := os.MkdirAll(SnapshotDirPathAbs, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	dump := replacer(spewConfig.Sdump(data...))
	snapshotName := strings.Replace(t.Name(), "/", "-", -1)
	snapshotAbsPath := filepath.Join(SnapshotDirPathAbs, snapshotName+".dump")

	if update {
		err := writeSnapshot(snapshotAbsPath, dump)
		if err != nil {
			t.Fatal(err)
		}

		t.Fatal(fmt.Errorf("Updating snapshot: '%s'", snapshotName))
	}

	prevSnapBytes, err := ioutil.ReadFile(snapshotAbsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = writeSnapshot(snapshotAbsPath, dump)
			if err != nil {
				t.Fatal(err)
			}

			t.Fatal(fmt.Errorf("No snapshot exists for name: '%s'. Creating new snapshot", snapshotName))
		}

		t.Fatal(err)
	}

	prevSnap := string(prevSnapBytes)
	if prevSnap != dump {
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(prevSnap),
			B:        difflib.SplitLines(dump),
			FromFile: "Previous",
			ToFile:   "Current",
			Context:  1,
		})
		if err != nil {
			t.Fatal(err)
		}

		t.Error(diff)
	}
}
func Snapshot(t TestingT, update bool, data ...interface{}) {
	t.Helper()
	SnapshotWithReplacer(t, update, defaultReplacer, data...)
}

func writeSnapshot(absPath string, dump string) error {
	return ioutil.WriteFile(absPath, []byte(dump), os.FileMode(0644))
}
