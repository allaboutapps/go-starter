package util

import (
	"path/filepath"
	"strings"
)

// FileNameWithoutExtension returns the name of the file referenced by the
// provided path without the file's extension.
// The function accepts a full (local) file path as well, only the latest
// element of the path will be considered as a name.
// If the provided path is empty or consists entirely of separators, an
// empty string will be returned.
func FileNameWithoutExtension(path string) string {
	base := filepath.Base(path)
	if base == "." {
		return ""
	} else if base == "/" {
		return ""
	}

	return strings.TrimSuffix(base, filepath.Ext(path))
}

// FileNameAndExtension returns the name of the file referenced by the
// provided path as well as its extension as separated strings.
// The function accepts a full (local) file path as well, only the latest
// element of the path will be considered as a name.
// If the provided path is empty or consists entirely of separators,
// empty strings will be returned.
func FileNameAndExtension(path string) (fileName string, extension string) {
	base := filepath.Base(path)
	if base == "." {
		return "", ""
	} else if base == "/" {
		return "", ""
	}

	extension = filepath.Ext(path)

	return strings.TrimSuffix(base, extension), extension
}
