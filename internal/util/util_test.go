package util_test

import (
	"encoding/base64"
	"encoding/hex"
	"math"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogLevelFromString(t *testing.T) {
	t.Parallel()

	res := util.LogLevelFromString("panic")
	assert.Equal(t, zerolog.PanicLevel, res)

	res = util.LogLevelFromString("warn")
	assert.Equal(t, zerolog.WarnLevel, res)

	res = util.LogLevelFromString("foo")
	assert.Equal(t, zerolog.DebugLevel, res)
}

func TestMergeStringMap(t *testing.T) {
	t.Parallel()

	baseMap := map[string]string{
		"A": "a",
		"B": "b",
		"C": "c",
	}

	toMerge := map[string]string{
		"C": "1",
		"D": "2",
	}

	expected := map[string]string{
		"A": "a",
		"B": "b",
		"C": "c",
		"D": "2",
	}

	res := util.MergeStringMap(baseMap, toMerge)
	assert.Equal(t, expected, res)

	expected = map[string]string{
		"C": "1",
		"D": "2",
		"A": "a",
		"B": "b",
	}

	res = util.MergeStringMap(toMerge, baseMap)
	assert.Equal(t, expected, res)
}

func TestMinAndMapInt(t *testing.T) {
	t.Parallel()

	max := math.MaxInt32
	min := math.MinInt32
	assert.Equal(t, max, util.MaxInt(max, min))
	assert.Equal(t, max, util.MaxInt(min, max))
	assert.Equal(t, min, util.MinInt(max, min))
	assert.Equal(t, min, util.MinInt(min, max))
	assert.Equal(t, min, util.MaxInt(min, min))
	assert.Equal(t, max, util.MinInt(max, max))
}

func TestGenerateRandom(t *testing.T) {
	t.Parallel()

	res, err := util.GenerateRandomBytes(13)
	require.NoError(t, err)
	assert.Len(t, res, 13)

	randString, err := util.GenerateRandomBase64String(17)
	require.NoError(t, err)
	res, err = base64.StdEncoding.DecodeString(randString)
	require.NoError(t, err)
	assert.Len(t, res, 17)

	randString, err = util.GenerateRandomHexString(19)
	require.NoError(t, err)
	res, err = hex.DecodeString(randString)
	require.NoError(t, err)
	assert.Len(t, res, 19)
}
