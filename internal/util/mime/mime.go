package mime

import "github.com/gabriel-vasile/mimetype"

var _ MIME = (*mimetype.MIME)(nil)

// MIME interface enables to use either *mimetype.MIME or KnownMIME as mimetype.
type MIME interface {
	String() string
	Extension() string
	Is(expectedMIME string) bool
}

// KnownMIME implements the MIME interface to be able to pass a *mimetype.MIME
// compatible value if the mimetype is already known so mimetype detection is not
// needed. It is therefore possible to skip mimetype detection if the mimetype is known
// or it is not possible to use a readSeeker but a mimetype is required.
type KnownMIME struct {
	MimeType      string
	FileExtension string
}

func (m *KnownMIME) String() string {
	return m.MimeType
}

func (m *KnownMIME) Extension() string {
	return m.FileExtension
}

func (m *KnownMIME) Is(expectedMIME string) bool {
	return expectedMIME == m.MimeType
}
