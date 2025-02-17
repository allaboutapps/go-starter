package util_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGetFileNameWithoutExtension(t *testing.T) {
	assert.Equal(t, "example", util.FileNameWithoutExtension("/a/b/c/d/example.jpg"))
	assert.Equal(t, "example", util.FileNameWithoutExtension("example.jpg"))
	assert.Equal(t, "example_test-check", util.FileNameWithoutExtension("example_test-check.jpg"))
	assert.Equal(t, "example", util.FileNameWithoutExtension("example"))
	assert.Equal(t, "", util.FileNameWithoutExtension(""))
	assert.Equal(t, "", util.FileNameWithoutExtension("."))
	assert.Equal(t, "", util.FileNameWithoutExtension("///"))
}

func TestFileNameAndExtension(t *testing.T) {
	name, ext := util.FileNameAndExtension("/a/b/c/d/example.jpg")
	assert.Equal(t, "example", name)
	assert.Equal(t, ".jpg", ext)
	name, ext = util.FileNameAndExtension("example.jpg")
	assert.Equal(t, "example", name)
	assert.Equal(t, ".jpg", ext)
	name, ext = util.FileNameAndExtension("example_test-check.jpg")
	assert.Equal(t, "example_test-check", name)
	assert.Equal(t, ".jpg", ext)
	name, ext = util.FileNameAndExtension("example")
	assert.Equal(t, "example", name)
	assert.Empty(t, ext)
	name, ext = util.FileNameAndExtension("")
	assert.Empty(t, name)
	assert.Empty(t, ext)
	name, ext = util.FileNameAndExtension(".")
	assert.Empty(t, name)
	assert.Empty(t, ext)
	name, ext = util.FileNameAndExtension("///")
	assert.Empty(t, name)
	assert.Empty(t, ext)
}
