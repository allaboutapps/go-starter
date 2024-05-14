package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/util"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/runtime"
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
	skips    []string
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

	s.save(t, dump)
}

// Save creates a dump of the given data.
// It will fail the test if the dump is different from the saved dump.
// It will also fail if it is the creation or an update of the snapshot.
// vastly inspired by https://github.com/bradleyjkemp/cupaloy
// main reason for self implementation is the replacer function and general flexibility
func (s snapshoter) SaveBytes(t TestingT, data []byte, fileExtensionOverride ...string) {
	t.Helper()
	err := os.MkdirAll(s.location, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	s.saveBytes(t, data, fileExtensionOverride...)
}

// SaveJSON creates a dump of the given data as JSON.
// It will fail the test if the dump is different from the saved dump.
// It will also fail if it is the creation or an update of the snapshot.
// vastly inspired by https://github.com/bradleyjkemp/cupaloy
// main reason for self implementation is the replacer function and general flexibility
func (s snapshoter) SaveJSON(t TestingT, data any) {
	t.Helper()

	// marshal data
	marshaled, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// indent data
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, marshaled, "", "\t"); err != nil {
		t.Fatal(err)
	}

	jsonS := s
	// set custom replacer for JSON compared to dumps
	jsonS.replacer = func(s string) string {
		skipString := strings.Join(jsonS.skips, "|")
		re, err := regexp.Compile(fmt.Sprintf(`"(?i)(%s)": .*`, skipString))
		if err != nil {
			panic(err)
		}

		// replace lines with property name + <redacted>
		return re.ReplaceAllString(s, `"$1": <redacted>,`)
	}

	jsonS.label += "JSON"
	jsonS.SaveString(t, prettyJSON.String())
}

// SaveUJSON is a short version for .Update(true).SaveJSON(...)
func (s snapshoter) SaveUJSON(t TestingT, data any) {
	t.Helper()
	s.Update(true).SaveJSON(t, data)
}

// SaveString creates a snapshot of the raw string.
// Used to snapshot payloads or mails as formatted data.
// It will fail the test if the dump is different from the saved dump.
// It will also fail if it is the creation or an update of the snapshot.
// vastly inspired by https://github.com/bradleyjkemp/cupaloy
// main reason for self implementation is the replacer function and general flexibility
func (s snapshoter) SaveString(t TestingT, data string) {
	t.Helper()
	err := os.MkdirAll(s.location, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	data = s.replacer(data)

	s.save(t, data)
}

func (s snapshoter) SaveUString(t TestingT, data string) {
	t.Helper()
	s.Update(true).save(t, data)
}

// SaveResponseAndValidate is used to create 2 snapshots for endpoint tests.
// One snapshot will save the raw JSON response as indented JSON.
// For the second snapshot the response will be parsed and validated using request helpers (helper_request.go)
// Afterwards a dump of the response will be saved.
// It will fail the test if the dump is different from the saved dump.
// It will also fail if it is the creation or an update of the snapshot.
func (s snapshoter) SaveResponseAndValidate(t TestingT, res *httptest.ResponseRecorder, v runtime.Validatable) {
	t.Helper()

	// snapshot prettyfied json first
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, res.Body.Bytes(), "", "\t"); err != nil {
		t.Fatal(err)
	}

	jsonS := s
	// set custom replacer for JSON compared to dumps
	jsonS.replacer = func(s string) string {
		skipString := strings.Join(jsonS.skips, "|")
		re, err := regexp.Compile(fmt.Sprintf(`"(?i)(%s)": .*`, skipString))
		if err != nil {
			panic(err)
		}

		// replace lines with property name + <redacted>
		return re.ReplaceAllString(s, `"$1": <redacted>,`)
	}

	jsonS.label += "JSON"
	jsonS.SaveString(t, prettyJSON.String())

	// bind and snapshot response type struct
	ParseResponseAndValidate(t, res, v)
	s.Save(t, v)
}

func (s snapshoter) SaveUResponseAndValidate(t TestingT, res *httptest.ResponseRecorder, v runtime.Validatable) {
	t.Helper()
	s.Update(true).SaveResponseAndValidate(t, res, v)
}

func (s snapshoter) save(t TestingT, dump string) {
	t.Helper()
	snapshotName := fmt.Sprintf("%s%s", strings.Replace(t.Name(), "/", "-", -1), s.label)
	snapshotAbsPath := filepath.Join(s.location, fmt.Sprintf("%s.golden", snapshotName))

	if s.update || UpdateGoldenGlobal {
		err := writeSnapshotString(snapshotAbsPath, dump)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("Updating snapshot: '%s'", snapshotName)
		return
	}

	prevSnapBytes, err := os.ReadFile(snapshotAbsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = writeSnapshotString(snapshotAbsPath, dump)
			if err != nil {
				t.Fatal(err)
			}

			t.Errorf("No snapshot exists for name: '%s'. Creating new snapshot", snapshotName)
			return
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

		t.Error(fmt.Sprintf("%s: %s", snapshotName, diff))
	}
}

func (s snapshoter) saveBytes(t TestingT, dump []byte, fileExtensionOverride ...string) {
	t.Helper()
	snapshotName := fmt.Sprintf("%s%s", strings.Replace(t.Name(), "/", "-", -1), s.label)

	fileExtension := "golden"
	if len(fileExtensionOverride) > 0 {
		fileExtension = fileExtensionOverride[0]
	}

	snapshotAbsPath := filepath.Join(s.location, fmt.Sprintf("%s.%s", snapshotName, fileExtension))

	if s.update || UpdateGoldenGlobal {
		err := writeSnapshot(snapshotAbsPath, dump)
		if err != nil {
			t.Fatal(err)
		}

		t.Errorf("Updating snapshot: '%s'", snapshotName)
		return
	}

	prevSnapBytes, err := os.ReadFile(snapshotAbsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = writeSnapshot(snapshotAbsPath, dump)
			if err != nil {
				t.Fatal(err)
			}

			t.Errorf("No snapshot exists for name: '%s'. Creating new snapshot", snapshotName)
			return
		}

		t.Fatal(err)
	}

	if !bytes.Equal(prevSnapBytes, dump) {
		t.Error(fmt.Sprintf("%s: Byte Snapshot Diff", snapshotName))
	}
}

// SaveU is a short version for .Update(true).Save(...)
func (s snapshoter) SaveU(t TestingT, data ...interface{}) {
	t.Helper()
	s.Update(true).Save(t, data...)
}

// Skip creates a custom replace function using a regex, this will replace any
// replacer function set in the Snapshoter.
// Each line of the formatted dump is matched against the property name defined in skip and
// the value will be replaced to deal with generated values that change each test.
func (s snapshoter) Skip(skip []string) snapshoter {
	s.skips = skip
	s.replacer = func(s string) string {
		skipString := fmt.Sprintf("\\s+%s", strings.Join(skip, "|\\s+"))

		re, err := regexp.Compile(fmt.Sprintf("(?m)(%s): .*[^{]$", skipString))
		if err != nil {
			panic(err)
		}

		reStruct, err := regexp.Compile(fmt.Sprintf("((%s): .*){\n([^}]|\n)*}", skipString))
		if err != nil {
			panic(err)
		}

		// replace lines with property name + <redacted>
		return reStruct.ReplaceAllString(re.ReplaceAllString(s, "$1: <redacted>,"), "$1 { <redacted> }")
	}

	return s
}

// Redact is a wrapper for Skip for easier usage with a variadic.
func (s snapshoter) Redact(skip ...string) snapshoter {
	return s.Skip(skip)
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

func writeSnapshotString(absPath string, dump string) error {
	return writeSnapshot(absPath, []byte(dump))
}

func writeSnapshot(absPath string, dump []byte) error {
	return os.WriteFile(absPath, dump, os.FileMode(0644))
}
