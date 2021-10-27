package i18n_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestI18NGlobal(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	i18n.InitPackage(config.I18n)

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

func TestI18N(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	bundle, err := i18n.NewBundle(config.I18n)
	require.NoError(t, err)

	msg := bundle.T("Test.Welcome", language.German, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	msg = bundle.T("Test.Welcome", language.English, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = bundle.T("Test.Welcome", language.Spanish, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = bundle.T("Test.Welcome", language.English)
	assert.Equal(t, "Welcome <no value>", msg)

	msg = bundle.T("Test.Body", language.German)
	assert.Equal(t, "Das ist ein Test", msg)

	msg = bundle.T("Test.Body", language.English)
	assert.Equal(t, "This is a test", msg)

	msg = bundle.T("Test.Body", language.Spanish)
	assert.Equal(t, "This is a test", msg)
}

func TestParseAcceptLanguageGlobal(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	i18n.InitPackage(config.I18n)

	tag := i18n.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.German, tag)
}

func TestParseAcceptLanguage(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	matcher, err := i18n.NewMatcher(config.I18n)
	require.NoError(t, err)

	tag := matcher.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.German, tag)
}

func TestParseLangGlobal(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	i18n.InitPackage(config.I18n)

	tag := i18n.ParseLang("de")
	assert.Equal(t, language.German, tag)
}

func TestParseLang(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	matcher, err := i18n.NewMatcher(config.I18n)
	require.NoError(t, err)

	tag := matcher.ParseLang("de")
	assert.Equal(t, language.German, tag)
}
