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

var (
	DefaultSnapshotDirPathAbs = filepath.Join(util.GetProjectRootDir(), "/test/testdata/snapshots")
	UpdateGoldenGlobal        = util.GetEnvAsBool("TEST_UPDATE_GOLDEN", false)
)

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

type snapshoter struct {
	update   bool
	label    string
	replacer func(s string) string
	location string
}

var Snapshoter = snapshoter{
	update:   false,
	label:    "",
	replacer: defaultReplacer,
	location: DefaultSnapshotDirPathAbs,
}

// Save creates a formatted dump of the given data.
// It will fail the test if the dump is different from the saved dump.
// It will also fail if it is the creation or an update of the snapshot.
// vastly inspired by https://github.com/bradleyjkemp/cupaloy
// main reason for self implementation is the replacer function and general flexibility
func (s snapshoter) Save(t TestingT, data ...interface{}) {
	t.Helper()
	err := os.MkdirAll(s.location, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	dump := s.replacer(spewConfig.Sdump(data...))
	snapshotName := fmt.Sprintf("%s%s", strings.Replace(t.Name(), "/", "-", -1), s.label)
	snapshotAbsPath := filepath.Join(s.location, fmt.Sprintf("%s.golden", snapshotName))

	if s.update || UpdateGoldenGlobal {
		err := writeSnapshot(snapshotAbsPath, dump)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("Updating snapshot: '%s'", snapshotName)
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

// SaveU is a short version for .Update(true).Save(...)
func (s snapshoter) SaveU(t TestingT, data ...interface{}) {
	s.Update(true).Save(t, data...)
}

// Skip creates a custom replace function using a regex, this will replace any
// replacer function set in the Snapshoter.
// Each line of the formatted dump is matched against the property name defined in skip and
// the value will be replaced to deal with generated values that change each test.
func (s snapshoter) Skip(skip []string) snapshoter {
	s.replacer = func(s string) string {
		skipString := strings.Join(skip, "|")
		re, err := regexp.Compile(fmt.Sprintf("(%s): .*", skipString))
		if err != nil {
			panic(err)
		}

		// replace lines with property name + <redacted>
		return re.ReplaceAllString(s, "$1: <redacted>,")
	}

	return s
}

// Upadte is used to force an update for the snapshot. Will fail the test.
func (s snapshoter) Update(update bool) snapshoter {
	s.update = update
	return s
}

// Label is used to add a suffix to the snapshots golden file.
func (s snapshoter) Label(label string) snapshoter {
	s.label = label
	return s
}

// Replacer is used to define a custom replace function in order to replace
// generated values (e.g. IDs).
func (s snapshoter) Replacer(replacer func(s string) string) snapshoter {
	s.replacer = replacer
	return s
}

// Location is used to save the golden file to a different location.
func (s snapshoter) Location(location string) snapshoter {
	s.location = location
	return s
}

func writeSnapshot(absPath string, dump string) error {
	return ioutil.WriteFile(absPath, []byte(dump), os.FileMode(0644))
}
