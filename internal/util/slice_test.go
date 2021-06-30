package util_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	test := []string{"a", "b", "d"}
	assert.True(t, util.ContainsString(test, "a"))
	assert.True(t, util.ContainsString(test, "b"))
	assert.False(t, util.ContainsString(test, "c"))
	assert.True(t, util.ContainsString(test, "d"))
}

func TestContainsAllString(t *testing.T) {
	test := []string{"a", "b", "d"}
	assert.True(t, util.ContainsAllString(test, "a"))
	assert.True(t, util.ContainsAllString(test, "b"))
	assert.False(t, util.ContainsAllString(test, "c"))
	assert.True(t, util.ContainsAllString(test, "d"))
	assert.True(t, util.ContainsAllString(test, "a", "b"))
	assert.True(t, util.ContainsAllString(test, "a", "d"))
	assert.True(t, util.ContainsAllString(test, "b", "d"))
	assert.False(t, util.ContainsAllString(test, "a", "c"))
	assert.False(t, util.ContainsAllString(test, "b", "c"))
	assert.False(t, util.ContainsAllString(test, "c", "d"))
	assert.True(t, util.ContainsAllString(test, "a", "b", "d"))
	assert.False(t, util.ContainsAllString(test, "a", "b", "c"))
	assert.False(t, util.ContainsAllString(test, "a", "b", "c", "d"))
	assert.True(t, util.ContainsAllString(test))
}

func TestUniqueString(t *testing.T) {
	test := []string{"a", "b", "d", "d", "a", "d"}
	assert.Equal(t, []string{"a", "b", "d"}, util.UniqueString(test))
}
