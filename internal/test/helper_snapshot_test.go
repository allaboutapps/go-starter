package test_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"allaboutapps.dev/aw/go-starter/internal/test/mocks"
	"allaboutapps.dev/aw/go-starter/internal/util"

	apitypes "allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	helloWorld = "Hello World!"
)

func TestSnapshot(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	data := struct {
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

	test.Snapshoter.Save(t, data, helloWorld)
}

func TestSnapshotWithReplacer(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	randID, err := util.GenerateRandomBase64String(20)
	require.NoError(t, err)
	data := struct {
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
	test.Snapshoter.Replacer(replacer).Save(t, data)
}

func TestSnapshotShouldFail(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	data := struct {
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

	tMock := new(mocks.TestingT)
	tMock.On("Helper").Return()
	tMock.On("Name").Return("TestSnapshotShouldFail")
	tMock.On("Error", mock.Anything).Return()
	test.Snapshoter.Save(tMock, data, helloWorld)
	tMock.AssertNotCalled(t, "Fatal")
	tMock.AssertNotCalled(t, "Fatalf")
	tMock.AssertCalled(t, "Error", mock.Anything)
}

func TestSnapshotWithUpdate(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	data := struct {
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

	tMock := new(mocks.TestingT)
	tMock.On("Helper").Return()
	tMock.On("Name").Return("TestSnapshotWithUpdate")
	tMock.On("Errorf", mock.Anything, mock.Anything).Return()
	test.Snapshoter.Update(true).Save(tMock, data, helloWorld)
	tMock.AssertNotCalled(t, "Error")
	tMock.AssertNotCalled(t, "Fatal")
	tMock.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestSnapshotNotExists(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	data := struct {
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

	defer func() {
		os.Remove(filepath.Join(test.DefaultSnapshotDirPathAbs, "TestSnapshotNotExists.golden"))
	}()

	tMock := new(mocks.TestingT)
	tMock.On("Helper").Return()
	tMock.On("Name").Return("TestSnapshotNotExists")
	tMock.On("Fatalf", mock.Anything, mock.Anything).Return()
	tMock.On("Fatal", mock.Anything).Return()
	tMock.On("Error", mock.Anything).Return()
	tMock.On("Errorf", mock.Anything, mock.Anything).Return()
	test.Snapshoter.Save(tMock, data, helloWorld)
	tMock.AssertNotCalled(t, "Error")
	tMock.AssertNotCalled(t, "Fatal")
	tMock.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestSnapshotSkipFields(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	randID, err := util.GenerateRandomBase64String(20)
	require.NoError(t, err)
	data := struct {
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

	test.Snapshoter.Skip([]string{"ID"}).Save(t, data)
}

func TestSnapshotSkipPrefixedFields(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}

	data := struct {
		ID            string
		OtherIDStr    string
		OtherIDInt    int
		OtherIDBool   bool
		OtherIDPTR    *string
		OtherIDStruct struct {
			ID string
		}
	}{
		ID:          "foo",
		OtherIDStr:  "id str",
		OtherIDInt:  4,
		OtherIDBool: true,
		OtherIDPTR:  swag.String("ID str ptr"),
		OtherIDStruct: struct{ ID string }{
			ID: "foo",
		},
	}

	test.Snapshoter.Skip([]string{"ID"}).Save(t, data)
}

func TestSnapshotSkipMultilineFields(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	randID, err := util.GenerateRandomBase64String(20)
	require.NoError(t, err)
	data := struct {
		ID string
		A  string
		B  int
		C  bool
		D  interface{}
		E  []string
		F  map[string]int
	}{
		ID: randID,
		A:  "foo",
		B:  1,
		C:  true,
		D: struct {
			Foo string
			Bar int
		}{
			Foo: "skip me",
			Bar: 3,
		},
		E: []string{"skip me", "skip me too"},
		F: map[string]int{
			"skip me":       1,
			"skip me too":   2,
			"skip me three": 3,
		},
	}

	test.Snapshoter.Skip([]string{"ID", "D", "E", "F"}).Save(t, data)
}

func TestSnapshotWithLabel(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	data := struct {
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

	test.Snapshoter.Label("_A").Save(t, data)
	test.Snapshoter.Label("_B").Save(t, helloWorld)
}

func TestSnapshotWithLocation(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}
	data := struct {
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

	location := filepath.Join(util.GetProjectRootDir(), "/internal/test/testdata")
	test.Snapshoter.Location(location).Save(t, data)
}

func TestSaveResponseAndValidate(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}

	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/auth/userinfo", nil, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response apitypes.GetUserInfoResponse
		test.Snapshoter.Redact("Email", "UpdatedAt", "updated_at").SaveResponseAndValidate(t, res, &response)
	})
}

func TestSnapshotJSON(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}

	randID, err := util.GenerateRandomBase64String(20)
	require.NoError(t, err)

	details := struct {
		ID string
		A  string
		B  int
		C  bool
		D  interface{}
		E  []string
		F  map[string]int
	}{
		ID: randID,
		A:  "foo",
		B:  1,
		C:  true,
		D: struct {
			Foo string
			Bar int
		}{
			Foo: "skip me",
			Bar: 3,
		},
		E: []string{"skip me", "skip me too"},
		F: map[string]int{
			"skip me":       1,
			"skip me too":   2,
			"skip me three": 3,
		},
	}

	marshaled, err := json.Marshal(details)
	require.NoError(t, err)

	test.Snapshoter.Redact("ID").SaveJSON(t, types.JSON(json.RawMessage(marshaled)))
}

func TestSnapshotSaveBytesImage(t *testing.T) {
	if test.UpdateGoldenGlobal {
		t.Skip()
	}

	filepath := filepath.Join(util.GetProjectRootDir(), "/test/testdata", "example.jpg")

	// read file and save bytes
	content, err := os.ReadFile(filepath)
	require.NoError(t, err)

	test.Snapshoter.SaveBytes(t, content, "jpg")
}
