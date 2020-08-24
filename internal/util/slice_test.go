package util_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	t.Parallel()

	test := []string{"a", "b", "d"}
	assert.True(t, util.ContainsString(test, "a"))
	assert.True(t, util.ContainsString(test, "b"))
	assert.False(t, util.ContainsString(test, "c"))
	assert.True(t, util.ContainsString(test, "d"))
}
