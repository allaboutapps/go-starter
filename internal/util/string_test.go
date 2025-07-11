package util_test

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandom(t *testing.T) {
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

	randString, err = util.GenerateRandomString(19, []util.CharRange{util.CharRangeAlphaLowerCase}, "/%$")
	require.NoError(t, err)
	assert.Len(t, randString, 19)
	for _, r := range randString {
		assert.True(t, (r >= 'a' && r <= 'z') || r == '/' || r == '%' || r == '$')
	}

	randString, err = util.GenerateRandomString(19, []util.CharRange{util.CharRangeAlphaUpperCase}, "^\"")
	require.NoError(t, err)
	assert.Len(t, randString, 19)
	for _, r := range randString {
		assert.True(t, (r >= 'A' && r <= 'Z') || r == '^' || r == '"')
	}

	randString, err = util.GenerateRandomString(19, []util.CharRange{util.CharRangeNumeric}, "")
	require.NoError(t, err)
	assert.Len(t, randString, 19)
	for _, r := range randString {
		assert.True(t, (r >= '0' && r <= '9'))
	}

	_, err = util.GenerateRandomString(1, nil, "")
	require.Error(t, err)

	randString, err = util.GenerateRandomString(8, nil, "a")
	require.NoError(t, err)
	assert.Len(t, randString, 8)
	assert.Equal(t, "aaaaaaaa", randString)
}

func TestNonEmptyOrNil(t *testing.T) {
	assert.Equal(t, "test", *util.NonEmptyOrNil("test"))
	assert.Equal(t, (*string)(nil), util.NonEmptyOrNil(""))
}

func TestContainsAll(t *testing.T) {
	assert.True(t, util.ContainsAll("Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "dolor"))
	assert.False(t, util.ContainsAll("Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "dolorx"))

	assert.True(t, util.ContainsAll("Lorem ipsum dolor sit amet, consectetur adipiscing elit.", ".", "sit", "elit", "ipsum", "Lorem ipsum"))
	assert.False(t, util.ContainsAll("Lorem ipsum dolor sit amet, consectetur adipiscing elit.", ".", "sit", "elit", "ipsum", " Lorem"))

	assert.True(t, util.ContainsAll("Lorem ipsum dolor sit amet,  ÄÜiö consectetur adipiscing elit.", "ÄÜiö c"))

	assert.False(t, util.ContainsAll("", "ÄÜiö c"))
	assert.True(t, util.ContainsAll("Lorem ipsum dolor sit amet,  ÄÜiö consectetur adipiscing elit.", ""))
}

func TestEmptyIfNil(t *testing.T) {
	s := "Lorem ipsum"
	assert.Equal(t, s, util.EmptyIfNil(&s))
	assert.Empty(t, util.EmptyIfNil(nil))
}
