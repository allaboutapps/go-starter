package test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

// SnapshotWithReplacer is similiar to Snapshot but with the addition of a custom replacer function,
// in order to replace generated values (e.g. IDs).
// vastly inspired by https://github.com/bradleyjkemp/cupaloy
// main reason for self implementation is the replacer function and general flexibility
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

		t.Fatalf("Updating snapshot: '%s'", snapshotName)
	}

	prevSnapBytes, err := ioutil.ReadFile(snapshotAbsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = writeSnapshot(snapshotAbsPath, dump)
			if err != nil {
				t.Fatal(err)
			}

			t.Fatalf("No snapshot exists for name: '%s'. Creating new snapshot", snapshotName)
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

// SnapshotWithSkipper creates a custom replace function using a regex.
// Each line of the formatted dump is matched against the property name defined in skip and
// the value will be replaced to deal with generated values that change each test.
// It will then call SnapshotWithReplacer with the replace function.
func SnapshotWithSkipper(t TestingT, update bool, skip []string, data ...interface{}) {
	t.Helper()
	replacer := func(s string) string {
		skipString := strings.Join(skip, "|")
		re, err := regexp.Compile(fmt.Sprintf("(%s): .*", skipString))
		if err != nil {
			t.Fatal("Could not compile regex")
		}

		// replace lines with property name + <redacted>
		return re.ReplaceAllString(s, "$1: <redacted>,")
	}

	SnapshotWithReplacer(t, update, replacer, data...)
}

// Snapshot creates a formatted dump of the given data.
// It will fail the test if the dump is different from the saved dump.
// It will also fail if it is the creation or an update of the snapshot.
func Snapshot(t TestingT, update bool, data ...interface{}) {
	t.Helper()
	SnapshotWithReplacer(t, update, defaultReplacer, data...)
}

func writeSnapshot(absPath string, dump string) error {
	return ioutil.WriteFile(absPath, []byte(dump), os.FileMode(0644))
}
