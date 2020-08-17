package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMgmtSecret(t *testing.T) {
	rs, err := GenerateRandomHexString(8)
	require.NoError(t, err)

	key := fmt.Sprintf("WE_WILL_NEVER_USE_THIS_MGMT_SECRET_%s", rs)
	expectedVal := fmt.Sprintf("SUPER_SECRET_%s", rs)

	err = os.Setenv(key, expectedVal)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		val := GetMgmtSecret(key)
		assert.Equal(t, expectedVal, val)
	}
}

func TestGetMgmtSecretRandom(t *testing.T) {
	expectedVal := GetMgmtSecret("DOES_NOT_EXIST_MGMT_SECRET")
	require.NotEmpty(t, expectedVal)

	for i := 0; i < 5; i++ {
		val := GetMgmtSecret("DOES_NOT_EXIST_MGMT_SECRET")
		assert.Equal(t, expectedVal, val)
	}
}
