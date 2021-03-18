package util_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_STRING"
	res := util.GetEnv(testVarKey, "noVal")
	assert.Equal(t, "noVal", res)

	os.Setenv(testVarKey, "string")
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

	os.Setenv(testVarKey, "mock")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvEnum(testVarKey, "smtp", []string{"mock", "smtp"})
	assert.Equal(t, "mock", res)

	os.Setenv(testVarKey, "foo")
	res = util.GetEnvEnum(testVarKey, "smtp", []string{"mock", "smtp"})
	assert.Equal(t, "smtp", res)
}

func TestGetEnvAsInt(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_INT"
	res := util.GetEnvAsInt(testVarKey, 1)
	assert.Equal(t, 1, res)

	os.Setenv(testVarKey, "2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsInt(testVarKey, 1)
	assert.Equal(t, 2, res)

	os.Setenv(testVarKey, "3x")
	res = util.GetEnvAsInt(testVarKey, 1)
	assert.Equal(t, 1, res)
}

func TestGetEnvAsUint32(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_UINT32"
	res := util.GetEnvAsUint32(testVarKey, 1)
	assert.Equal(t, uint32(1), res)

	os.Setenv(testVarKey, "2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsUint32(testVarKey, 1)
	assert.Equal(t, uint32(2), res)

	os.Setenv(testVarKey, "3x")
	res = util.GetEnvAsUint32(testVarKey, 1)
	assert.Equal(t, uint32(1), res)
}

func TestGetEnvAsUint8(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_UINT8"
	res := util.GetEnvAsUint8(testVarKey, 1)
	assert.Equal(t, uint8(1), res)

	os.Setenv(testVarKey, "2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsUint8(testVarKey, 1)
	assert.Equal(t, uint8(2), res)

	os.Setenv(testVarKey, "3x")
	res = util.GetEnvAsUint8(testVarKey, 1)
	assert.Equal(t, uint8(1), res)
}

func TestGetEnvAsBool(t *testing.T) {
	testVarKey := "TEST_ONLY_FOR_UNIT_TEST_BOOL"
	res := util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, true, res)

	os.Setenv(testVarKey, "f")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, false, res)

	os.Setenv(testVarKey, "0")
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, false, res)

	os.Setenv(testVarKey, "false")
	res = util.GetEnvAsBool(testVarKey, true)
	assert.Equal(t, false, res)

	os.Setenv(testVarKey, "3x")
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

	os.Setenv(testVarKey, "https://allaboutapps.at/")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsURL(testVarKey, "foo")
	assert.Equal(t, *testURL, *res)

	os.Setenv(testVarKey, "%")
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

	os.Setenv(testVarKey, "1,2")
	defer os.Unsetenv(testVarKey)
	res = util.GetEnvAsStringArr(testVarKey, testVal)
	assert.Equal(t, []string{"1", "2"}, res)

	os.Setenv(testVarKey, "")
	res = util.GetEnvAsStringArr(testVarKey, testVal)
	assert.Equal(t, testVal, res)
}

func TestGetMgmtSecret(t *testing.T) {
	rs, err := util.GenerateRandomHexString(8)
	require.NoError(t, err)

	key := fmt.Sprintf("WE_WILL_NEVER_USE_THIS_MGMT_SECRET_%s", rs)
	expectedVal := fmt.Sprintf("SUPER_SECRET_%s", rs)

	err = os.Setenv(key, expectedVal)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		val := util.GetMgmtSecret(key)
		assert.Equal(t, expectedVal, val)
	}
}

func TestGetMgmtSecretRandom(t *testing.T) {
	expectedVal := util.GetMgmtSecret("DOES_NOT_EXIST_MGMT_SECRET")
	require.NotEmpty(t, expectedVal)

	for i := 0; i < 5; i++ {
		val := util.GetMgmtSecret("DOES_NOT_EXIST_MGMT_SECRET")
		assert.Equal(t, expectedVal, val)
	}
}
