package mime_test

import (
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/mime"
	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnownMIME(t *testing.T) {
	filePath := filepath.Join(util.GetProjectRootDir(), "test", "testdata", "example.jpg")

	var detectedMIME mime.MIME
	var err error
	detectedMIME, err = mimetype.DetectFile(filePath)
	require.NoError(t, err)

	var knownMIME mime.MIME = &mime.KnownMIME{
		MimeType:      "image/jpeg",
		FileExtension: ".jpg",
	}

	assert.Equal(t, detectedMIME.Extension(), knownMIME.Extension())
	assert.Equal(t, detectedMIME.String(), knownMIME.String())
	assert.True(t, knownMIME.Is(detectedMIME.String()))
}
