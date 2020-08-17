package test_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/mocks"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSnapshot(t *testing.T) {
	t.Parallel()

	a := struct {
		A string
		B int
		C bool
		D *string
	}{
		A: "foo",
		B: 1,
		C: true,
		D: swag.String("bar"),
	}

	b := "Hello World!"

	test.Snapshot(t, false, a, b)
}

func TestSnapshotWithReplacer(t *testing.T) {
	t.Parallel()

	randID, err := util.GenerateRandomBase64String(20)
	require.NoError(t, err)
	a := struct {
		ID string
		A  string
		B  int
		C  bool
		D  *string
	}{
		ID: randID,
		A:  "foo",
		B:  1,
		C:  true,
		D:  swag.String("bar"),
	}

	replacer := func(s string) string {
		re, err := regexp.Compile(`ID:.*"(.*)",`)
		require.NoError(t, err)
		return re.ReplaceAllString(s, "ID: <redacted>,")
	}
	test.SnapshotWithReplacer(t, false, replacer, a)
}

func TestSnapshotShouldFail(t *testing.T) {
	t.Parallel()

	a := struct {
		A string
		B int
		C bool
		D *string
	}{
		A: "fo",
		B: 1,
		C: true,
		D: swag.String("bar"),
	}

	b := "Hello World!"

	tMock := new(mocks.TestingT)
	tMock.On("Helper").Return()
	tMock.On("Name").Return("TestSnapshotShouldFail")
	tMock.On("Error", mock.Anything).Return()
	test.Snapshot(tMock, false, a, b)
	tMock.AssertNotCalled(t, "Fatal")
	tMock.AssertCalled(t, "Error", mock.Anything)
}

func TestSnapshotWithUpdate(t *testing.T) {
	t.Parallel()

	a := struct {
		A string
		B int
		C bool
		D *string
	}{
		A: "fo",
		B: 1,
		C: true,
		D: swag.String("bar"),
	}

	b := "Hello World!"

	tMock := new(mocks.TestingT)
	tMock.On("Helper").Return()
	tMock.On("Name").Return("TestSnapshotWithUpdate")
	tMock.On("Fatal", mock.Anything).Return()
	test.Snapshot(tMock, true, a, b)
	tMock.AssertNotCalled(t, "Error")
	tMock.AssertCalled(t, "Fatal", mock.Anything)
}

func TestSnapshotNotExists(t *testing.T) {
	t.Parallel()

	a := struct {
		A string
		B int
		C bool
		D *string
	}{
		A: "foo",
		B: 1,
		C: true,
		D: swag.String("bar"),
	}

	b := "Hello World!"

	defer func() {
		os.Remove(filepath.Join(test.SnapshotDirPathAbs, "TestSnapshotNotExists.dump"))
	}()

	tMock := new(mocks.TestingT)
	tMock.On("Helper").Return()
	tMock.On("Name").Return("TestSnapshotNotExists")
	tMock.On("Fatal", mock.Anything).Return()
	tMock.On("Error", mock.Anything).Return()
	test.Snapshot(tMock, false, a, b)
	tMock.AssertNotCalled(t, "Error")
	tMock.AssertCalled(t, "Fatal", mock.Anything)
}
