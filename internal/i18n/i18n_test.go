package i18n_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestI18N(t *testing.T) {
	msg := i18n.T("Test.Welcome", language.German, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	msg = i18n.T("Test.Welcome", language.English, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = i18n.T("Test.Welcome", language.Spanish, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = i18n.T("Test.Welcome", language.English)
	assert.Equal(t, "Welcome <no value>", msg)

	msg = i18n.T("Test.Body", language.German)
	assert.Equal(t, "Das ist ein Test", msg)

	msg = i18n.T("Test.Body", language.English)
	assert.Equal(t, "This is a test", msg)

	msg = i18n.T("Test.Body", language.Spanish)
	assert.Equal(t, "This is a test", msg)
}
