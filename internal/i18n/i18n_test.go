package i18n_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestI18N(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	msg := messages.T("Test.Welcome", language.German, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	msg = messages.T("Test.Welcome", language.English, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = messages.T("Test.Welcome", language.Spanish, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = messages.T("Test.Welcome", language.English)
	assert.Equal(t, "Welcome <no value>", msg)

	msg = messages.T("Test.Body", language.German)
	assert.Equal(t, "Das ist ein Test", msg)

	msg = messages.T("Test.Body", language.English)
	assert.Equal(t, "This is a test", msg)

	msg = messages.T("Test.Body", language.Spanish)
	assert.Equal(t, "This is a test", msg)

	msg = messages.T("Test.Invalid.Key.Does.Not.Exist", language.English)
	assert.Equal(t, "Test.Invalid.Key.Does.Not.Exist", msg)

	msg = messages.T("Test.Invalid.Key.Does.Not.Exist", language.German)
	assert.Equal(t, "Test.Invalid.Key.Does.Not.Exist", msg)

	msg = messages.T("Test.String.DE.only", language.English)
	assert.Equal(t, "Test.String.DE.only", msg)

	msg = messages.T("Test.String.DE.only", language.German)
	assert.Equal(t, "This key only exists in DE", msg)

	msg = messages.T("Test.String.EN.only", language.English)
	assert.Equal(t, "This key only exists in EN", msg)

	msg = messages.T("Test.String.EN.only", language.German)
	assert.Equal(t, "Test.String.EN.only", msg)

	// ensure language subvariants are supported
	deAt := messages.ParseLang("de-AT")
	assert.NotEqual(t, language.German, deAt)
	msg = messages.T("Test.Body", deAt)
	assert.Equal(t, "Das ist ein Test", msg)
}

func TestParseAcceptLanguageOverall(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.German, tag)
}

func TestParseAcceptLanguageSpecific(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseAcceptLanguage("de-AT,en-US;q=0.7,en;q=0.3")
	assert.NotEqual(t, language.German, tag) // actual: de-u-rg-atzzzz
}

func TestParseLang(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseLang("de")
	assert.Equal(t, language.German, tag)
}

func TestParseLangWellFormedUnknownLangTag(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseLang("xx")
	assert.Equal(t, config.I18n.DefaultLanguage, tag)

	msg := messages.T("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}

func TestParseLangInvalidLangTag(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseLang("§$%/%&/(/&%/)(")
	assert.Equal(t, config.I18n.DefaultLanguage, tag)

	msg := messages.T("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}

func TestParseAcceptLanguageWellFormedUnknownLangTag(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseAcceptLanguage("xx,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, config.I18n.DefaultLanguage, tag)

	msg := messages.T("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}
func TestParseAcceptLanguageInvalidLangTag(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	messages, err := i18n.New(config.I18n)
	require.NoError(t, err)

	tag := messages.ParseAcceptLanguage("§$%/%&/(/&%/)(")
	assert.Equal(t, config.I18n.DefaultLanguage, tag)

	msg := messages.T("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}
