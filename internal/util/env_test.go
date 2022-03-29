package util_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestGetEnv(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_STRING"
	res := util.GetEnv(testVarKey, "noVal")
	assert.Equal(t, "noVal", res)

	t.Setenv(testVarKey, "string")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnv(testVarKey, "noVal")
	assert.Equal(t, "string", res)
}

func TestGetEnvEnum(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_ENUM"

	panicFunc := func() {
		_ = util.GetEnvEnum(testVarKey, "smtp", []string{"mock", "foo"})
	}
	assert.Panics(t, panicFunc)

	res := util.GetEnvEnum(testVarKey, "smtp", []string{"mock", "smtp"})
	assert.Equal(t, "smtp", res)

	t.Setenv(testVarKey, "mock")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvEnum(testVarKey, "smtp", []string{"mock", "smtp"})
	assert.Equal(t, "mock", res)

	t.Setenv(testVarKey, "foo")
	res = util.GetEnvEnum(testVarKey, "smtp", []string{"mock", "smtp"})
	assert.Equal(t, "smtp", res)
}

func TestGetEnvAsInt(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_INT"
	res := util.GetEnvAsInt(testVarKey, 1)
	assert.Equal(t, 1, res)

	t.Setenv(testVarKey, "2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsInt(testVarKey, 1)
	assert.Equal(t, 2, res)

	t.Setenv(testVarKey, "3x")
	res = util.GetEnvAsInt(testVarKey, 1)
	assert.Equal(t, 1, res)
}

func TestGetEnvAsUint32(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_UINT32"
	res := util.GetEnvAsUint32(testVarKey, 1)
	assert.Equal(t, uint32(1), res)

	t.Setenv(testVarKey, "2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsUint32(testVarKey, 1)
	assert.Equal(t, uint32(2), res)

	t.Setenv(testVarKey, "3x")
	res = util.GetEnvAsUint32(testVarKey, 1)
	assert.Equal(t, uint32(1), res)
}

func TestGetEnvAsUint8(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_UINT8"
	res := util.GetEnvAsUint8(testVarKey, 1)
	assert.Equal(t, uint8(1), res)

	t.Setenv(testVarKey, "2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsUint8(testVarKey, 1)
	assert.Equal(t, uint8(2), res)

	t.Setenv(testVarKey, "3x")
	res = util.GetEnvAsUint8(testVarKey, 1)
	assert.Equal(t, uint8(1), res)
}

func TestGetEnvAsBool(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_BOOL"
	res := util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, true, res)

	t.Setenv(testVarKey, "f")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, false, res)

	t.Setenv(testVarKey, "0")
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, false, res)

	t.Setenv(testVarKey, "false")
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, false, res)

	t.Setenv(testVarKey, "3x")
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, true, res)
}

func TestGetEnvAsURL(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_URL"
	testURL, err := url.Parse("https://allaboutapps.at/")
	require.NoError(t, err)

	panicFunc := func() {
		_ = util.GetEnvAsURL(testVarKey, "%")
	}
	assert.Panics(t, panicFunc)

	res := util.GetEnvAsURL(testVarKey, "https://allaboutapps.at/")
	assert.Equal(t, *testURL, *res)

	t.Setenv(testVarKey, "https://allaboutapps.at/")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsURL(testVarKey, "foo")
	assert.Equal(t, *testURL, *res)

	t.Setenv(testVarKey, "%")
	panicFunc = func() {
		_ = util.GetEnvAsURL(testVarKey, "https://allaboutapps.at/")
	}
	assert.Panics(t, panicFunc)
}

func TestGetEnvAsStringArr(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_STRING_ARR"
	testVal := []string{"a", "b", "c"}
	res := util.GetEnvAsStringArr(testVarKey, testVal)
	assert.Equal(t, testVal, res)

	t.Setenv(testVarKey, "1,2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsStringArr(testVarKey, testVal)
	assert.Equal(t, []string{"1", "2"}, res)

	t.Setenv(testVarKey, "")
	res = util.GetEnvAsStringArr(testVarKey, testVal)
	assert.Equal(t, testVal, res)

	t.Setenv(testVarKey, "a, b, c")
	res = util.GetEnvAsStringArr(testVarKey, testVal)
	assert.Equal(t, []string{"a", " b", " c"}, res)

	t.Setenv(testVarKey, "a|b|c")
	res = util.GetEnvAsStringArr(testVarKey, testVal, "|")
	assert.Equal(t, []string{"a", "b", "c"}, res)

	t.Setenv(testVarKey, "a,b,c")
	res = util.GetEnvAsStringArr(testVarKey, testVal, "|")
	assert.Equal(t, []string{"a,b,c"}, res)

	t.Setenv(testVarKey, "a||b||c")
	res = util.GetEnvAsStringArr(testVarKey, testVal, "||")
	assert.Equal(t, []string{"a", "b", "c"}, res)
}

func TestGetEnvAsStringArrTrimmed(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_STRING_ARR_TRIMMED"
	testVal := []string{"a", "b", "c"}

	t.Setenv(testVarKey, "a, b, c")
	defer os.Unsetenv(testVarKey)
	res := util.GetEnvAsStringArrTrimmed(testVarKey, testVal)
	assert.Equal(t, []string{"a", "b", "c"}, res)

	t.Setenv(testVarKey, "a,   b,c    ")
	res = util.GetEnvAsStringArrTrimmed(testVarKey, testVal)
	assert.Equal(t, []string{"a", "b", "c"}, res)

	t.Setenv(testVarKey, "  a || b  || c  ")
	res = util.GetEnvAsStringArrTrimmed(testVarKey, testVal, "||")
	assert.Equal(t, []string{"a", "b", "c"}, res)
}

func TestGetMgmtSecret(t *testing.T) {
	rs, err := util.GenerateRandomHexString(8)
	require.NoError(t, err)

	key := fmt.Sprintf("WE_WILL_NEVER_USE_THIS_MGMT_SECRET_%s", rs)
	expectedVal := fmt.Sprintf("SUPER_SECRET_%s", rs)

	t.Setenv(key, expectedVal)

	for i := 0; i < 5; i++ {
		val := util.GetMgmtSecret(key)
		assert.Equal(t, expectedVal, val)
	}
}

func TestGetEnvAsLanguageTag(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_LANG"
	res := util.GetEnvAsLanguageTag(testVarKey, language.German)
	assert.Equal(t, language.German, res)

	t.Setenv(testVarKey, "en")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsLanguageTag(testVarKey, language.German)
	assert.Equal(t, language.English, res)
}

func TestGetEnvAsLanguageTagArr(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_LANG_ARR"
	testVal := []language.Tag{language.German, language.English, language.Spanish}
	res := util.GetEnvAsLanguageTagArr(testVarKey, testVal)
	assert.Equal(t, testVal, res)

	t.Setenv(testVarKey, "de,en")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsLanguageTagArr(testVarKey, testVal)
	assert.Equal(t, []language.Tag{language.German, language.English}, res)

	t.Setenv(testVarKey, "")
	res = util.GetEnvAsLanguageTagArr(testVarKey, testVal)
	assert.Equal(t, testVal, res)

	t.Setenv(testVarKey, "en|es")
	res = util.GetEnvAsLanguageTagArr(testVarKey, testVal, "|")
	assert.Equal(t, []language.Tag{language.English, language.Spanish}, res)

	t.Setenv(testVarKey, "en||es")
	res = util.GetEnvAsLanguageTagArr(testVarKey, testVal, "||")
	assert.Equal(t, []language.Tag{language.English, language.Spanish}, res)
}

func TestGetMgmtSecretRandom(t *testing.T) {
	expectedVal := util.GetMgmtSecret("DOES_NOT_EXIST_MGMT_SECRET")
	require.NotEmpty(t, expectedVal)

	for i := 0; i < 5; i++ {
		val := util.GetMgmtSecret("DOES_NOT_EXIST_MGMT_SECRET")
		assert.Equal(t, expectedVal, val)
	}
}
