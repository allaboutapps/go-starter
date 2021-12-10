package i18n_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestServerProvidedI18n(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		// uses /app/web/i18n by default
		// expect all i18n files were loaded and the defaultLanguage matches the FIRST priority Tag.
		assert.Equal(t, s.Config.I18n.DefaultLanguage, s.I18n.Tags()[0])

		// expect all i18ns were loaded...
		files, err := os.ReadDir(s.Config.I18n.BundleDirAbs)
		require.NoError(t, err)

		msgFilesCount := 0

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
				continue
			}

			msgFilesCount++
		}

		if msgFilesCount == 0 {
			// no i18n bundles were available, as the defaultLanguage is a tag itself, check for len 1
			assert.Equal(t, len(s.I18n.Tags()), 1)
		} else {
			assert.Equal(t, len(s.I18n.Tags()), msgFilesCount)
		}

		msg := s.I18n.Translate("this.key.will.never.exist", s.Config.I18n.DefaultLanguage)
		assert.Equal(t, "this.key.will.never.exist", msg)
	})
}

// Note that all following tests use a special message directory within /internal/i18n/testdata.
// We do this to ensure we don't depend on your project specific i18n bundle/configuration,
// that you would typically store within /web/i18n.

func TestI18n(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})
	require.NoError(t, err)

	assert.Equal(t, language.English, srv.Tags()[0])
	assert.Equal(t, language.German, srv.Tags()[1])
	assert.Equal(t, 2, len(srv.Tags()))

	msg := srv.Translate("Test.Welcome", language.German, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	msg = srv.Translate("Test.Welcome", language.English, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = srv.Translate("Test.Welcome", language.Spanish, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = srv.Translate("Test.Welcome", language.English, i18n.Data{"Name": "Franz"})
	assert.Equal(t, "Welcome Franz", msg)

	msg = srv.Translate("Test.Welcome", language.English)
	assert.Equal(t, "Welcome <no value>", msg)

	msg = srv.Translate("Test.Body", language.German)
	assert.Equal(t, "Das ist ein Test", msg)

	msg = srv.Translate("Test.Body", language.English)
	assert.Equal(t, "This is a test", msg)

	msg = srv.Translate("Test.Body", language.Spanish)
	assert.Equal(t, "This is a test", msg)

	msg = srv.Translate("Test.Invalid.Key.Does.Not.Exist", language.English)
	assert.Equal(t, "Test.Invalid.Key.Does.Not.Exist", msg)

	msg = srv.Translate("Test.Invalid.Key.Does.Not.Exist", language.German)
	assert.Equal(t, "Test.Invalid.Key.Does.Not.Exist", msg)

	msg = srv.Translate("Test.String.DE.only", language.English)
	assert.Equal(t, "Test.String.DE.only", msg)

	msg = srv.Translate("Test.String.DE.only", language.German)
	assert.Equal(t, "This key only exists in DE", msg)

	msg, err = srv.TranslateMaybe("Test.String.DE.only", language.English)
	assert.Error(t, err)
	assert.Equal(t, "", msg) // no fallback!

	msg = srv.Translate("Test.String.EN.only", language.English)
	assert.Equal(t, "This key only exists in EN", msg)

	msg, err = srv.TranslateMaybe("Test.String.EN.only", language.German)
	assert.Error(t, err)
	assert.Equal(t, "This key only exists in EN", msg) // fallback (but error)

	msg = srv.Translate("Test.String.EN.only", language.German)
	assert.Equal(t, "Test.String.EN.only", msg)

	msg = srv.Translate("", language.German) // empty key
	assert.Equal(t, "", msg)

	msg, err = srv.TranslateMaybe("", language.German) // empty key
	assert.Error(t, err)
	assert.Equal(t, "", msg)

	// ensure language subvariants are supported
	deAt := srv.ParseLang("de-AT")
	assert.NotEqual(t, language.German, deAt)
	msg = srv.Translate("Test.Body", deAt)
	assert.Equal(t, "Das ist ein Test", msg)
}

func TestI18nConcurrentUsage(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(index int) {
			msg := srv.Translate("Test.Welcome", language.German, i18n.Data{"Name": fmt.Sprintf("%v", index)})
			assert.Equal(t, fmt.Sprintf("Guten Tag %v", index), msg)

			msg = srv.Translate("Test.Welcome", language.English, i18n.Data{"Name": fmt.Sprintf("%v", index)})
			assert.Equal(t, fmt.Sprintf("Welcome %v", index), msg)

			msg = srv.Translate("Test.Welcome", language.Spanish, i18n.Data{"Name": fmt.Sprintf("%v", index)})
			assert.Equal(t, fmt.Sprintf("Welcome %v", index), msg)

			msg = srv.Translate("Test.Welcome", language.English, i18n.Data{"Name": "Franz"})
			assert.Equal(t, "Welcome Franz", msg)

			msg = srv.Translate("Test.Welcome", language.English)
			assert.Equal(t, "Welcome <no value>", msg)

			msg = srv.Translate("Test.Body", language.German)
			assert.Equal(t, "Das ist ein Test", msg)

			msg = srv.Translate("Test.Body", language.English)
			assert.Equal(t, "This is a test", msg)

			msg = srv.Translate("Test.Body", language.Spanish)
			assert.Equal(t, "This is a test", msg)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestI18nOtherDefault(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.German,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})
	require.NoError(t, err)

	assert.Equal(t, language.German, srv.Tags()[0])
	assert.Equal(t, language.English, srv.Tags()[1])
	assert.Equal(t, 2, len(srv.Tags()))
}

func TestI18nInexistantDefault(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.Italian,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})
	require.NoError(t, err)

	assert.Equal(t, language.Italian, srv.Tags()[0])
	assert.Equal(t, language.German, srv.Tags()[1])
	assert.Equal(t, language.English, srv.Tags()[2])
	assert.Equal(t, 3, len(srv.Tags()))
}

func TestI18nEmpty(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.Italian,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n-empty"),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(srv.Tags())) // the DefaultLanguage should still be set!
	assert.Equal(t, language.Italian, srv.Tags()[0])

	tag := srv.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.Italian, tag)

	msg := srv.Translate("no.test.key.exists", tag)
	assert.Equal(t, "no.test.key.exists", msg)

	msg, err = srv.TranslateMaybe("no.test.key.exists", tag)
	assert.Error(t, err)
	assert.Equal(t, "", msg)

	msg = srv.Translate("no.test.key.exists", language.Ukrainian)
	assert.Equal(t, "no.test.key.exists", msg)

	msg, err = srv.TranslateMaybe("no.test.key.exists", language.Ukrainian)
	assert.Error(t, err)
	assert.Equal(t, "", msg)
}

func TestI18nSpecialized(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.AmericanEnglish,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n-specialized"),
	})

	require.NoError(t, err)
	assert.Equal(t, 4, len(srv.Tags())) // specialized subvariant is default
	assert.Equal(t, language.AmericanEnglish, srv.Tags()[0])

	msg := srv.Translate("test.punchline", language.AmericanEnglish)
	assert.Equal(t, "I can has HUMOR?", msg)

	msg = srv.Translate("test.punchline", language.BritishEnglish)
	assert.Equal(t, "I can has HUMOUR?", msg)

	msg = srv.Translate("test.punchline", language.English)
	assert.Equal(t, "I can has HUMOR?", msg) // fall back to default

	msg = srv.Translate("test.punchline", language.German)
	assert.Equal(t, "Habe ich Humor?", msg) // jump to parsed Austrian German

	tag := srv.ParseAcceptLanguage("de-at,en-US;q=0.7,en;q=0.3") // explicit Austrian tag
	msg = srv.Translate("test.punchline", tag)
	assert.Equal(t, "Koan i Humor?", msg) // jump to parsed Austrian German
}

func TestReservedKeywordsResolve(t *testing.T) {
	// "reserved" keys:
	// "id", "description", "hash", "leftdelim", "rightdelim", "zero", "one", "two", "few", "many", "other"
	// see https://github.com/nicksnyder/go-i18n/blob/2180cd9f35b3e125cfe3773a6bf3ea483347f060/v2/i18n/message.go#L181

	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n-reserved"),
	})

	require.NoError(t, err)
	assert.Equal(t, 2, len(srv.Tags()))

	// single reserved word (not yet mapped, requires 2 keywords at least)
	msg := srv.Translate("reserved1.zero", language.English)
	assert.Equal(t, "Zero", msg)

	msg = srv.Translate("reserved2.one", language.English)
	assert.Equal(t, "One", msg)

	msg = srv.Translate("reserved3.two", language.English)
	assert.Equal(t, "Two", msg)

	msg = srv.Translate("reserved4.few", language.English)
	assert.Equal(t, "Few", msg)

	msg = srv.Translate("reserved5.many", language.English)
	assert.Equal(t, "Many", msg)

	msg = srv.Translate("reserved6.other", language.English)
	assert.Equal(t, "Other", msg)

	msg = srv.Translate("reserved7.id", language.English)
	assert.Equal(t, "id", msg)

	msg = srv.Translate("reserved8.description", language.English)
	assert.Equal(t, "Description", msg)

	// single parent is not directly resolveable
	msg = srv.Translate("reserved2", language.English)
	assert.Equal(t, "reserved2", msg)

	// german: single reserved word
	msg = srv.Translate("reserved1.zero", language.German)
	assert.Equal(t, "Null", msg)

	msg = srv.Translate("reserved2.one", language.German)
	assert.Equal(t, "Eins", msg)

	msg = srv.Translate("reserved3.two", language.German)
	assert.Equal(t, "Zwei", msg)

	msg = srv.Translate("reserved4.few", language.German)
	assert.Equal(t, "Wenig", msg)

	msg = srv.Translate("reserved5.many", language.German)
	assert.Equal(t, "Mehr", msg)

	msg = srv.Translate("reserved6.other", language.German)
	assert.Equal(t, "Andere", msg)

	msg = srv.Translate("reserved7.id", language.German)
	assert.Equal(t, "ID", msg)

	msg = srv.Translate("reserved8.description", language.German)
	assert.Equal(t, "Beschreibung", msg)

	// Combined toml map: all reserved words
	// This does not work as it's parsed as map and CLDR plural rules now apply!
	// The last key is interpreted as https://cldr.unicode.org/index/cldr-spec/plural-rules
	msg = srv.Translate("reservedMap.zero", language.English)
	assert.Equal(t, "reservedMap.zero", msg)

	msg = srv.Translate("reservedMap.one", language.English)
	assert.Equal(t, "reservedMap.one", msg)

	msg = srv.Translate("reservedMap.two", language.English)
	assert.Equal(t, "reservedMap.two", msg)

	msg = srv.Translate("reservedMap.few", language.English)
	assert.Equal(t, "reservedMap.few", msg)

	msg = srv.Translate("reservedMap.many", language.English)
	assert.Equal(t, "reservedMap.many", msg)

	msg = srv.Translate("reservedMap.other", language.English)
	assert.Equal(t, "reservedMap.other", msg)

	msg = srv.TranslatePlural("reservedMap", 0, language.English)
	assert.Equal(t, "Other", msg)

	msg = srv.TranslatePlural("reservedMap", 1, language.English)
	assert.Equal(t, "One", msg)

	msg = srv.TranslatePlural("reservedMap", 2, language.English)
	assert.Equal(t, "Other", msg)

	msg = srv.TranslatePlural("reservedMap", "asdfasdf", language.English)
	assert.Equal(t, "reservedMap (count=asdfasdf)", msg)

	// plain toml: all reserved words
	// This does not work as it's parsed as map and CLDR plural rules now apply!
	// The last key is interpreted as https://cldr.unicode.org/index/cldr-spec/plural-rules
	msg = srv.Translate("reserved.plain.zero", language.English)
	assert.Equal(t, "reserved.plain.zero", msg)

	msg = srv.Translate("reserved.plain.one", language.English)
	assert.Equal(t, "reserved.plain.one", msg)

	msg = srv.Translate("reserved.plain.two", language.English)
	assert.Equal(t, "reserved.plain.two", msg)

	msg = srv.Translate("reserved.plain.few", language.English)
	assert.Equal(t, "reserved.plain.few", msg)

	msg = srv.Translate("reserved.plain.many", language.English)
	assert.Equal(t, "reserved.plain.many", msg)

	msg = srv.Translate("reserved.plain.other", language.English)
	assert.Equal(t, "reserved.plain.other", msg)

	msg = srv.Translate("reserved.plain2.id", language.English)
	assert.Equal(t, "reserved.plain2.id", msg)

	msg = srv.Translate("reserved.plain2.description", language.English)
	assert.Equal(t, "reserved.plain2.description", msg)

	msg = srv.Translate("reserved.plain3.id", language.English)
	assert.Equal(t, "reserved.plain3.id", msg)

	msg = srv.Translate("reserved.plain3.description", language.English)
	assert.Equal(t, "reserved.plain3.description", msg)

	msg = srv.Translate("reserved.plain3.test", language.English)
	assert.Equal(t, "reserved.plain3.test", msg)

}

func TestI18nPlural(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n-plural"),
	})

	require.NoError(t, err)
	assert.Equal(t, 2, len(srv.Tags()))

	msg := srv.TranslatePlural("cats", 0, language.AmericanEnglish)
	assert.Equal(t, "I've 0 cats.", msg) // zero is not supported for CLDR English, "I don't have a cat." is not possible!

	msg = srv.TranslatePlural("cats", 1, language.BritishEnglish)
	assert.Equal(t, "I've one cat.", msg)

	msg = srv.TranslatePlural("cats", 2, language.English)
	assert.Equal(t, "I've 2 cats.", msg)

	msg = srv.TranslatePlural("cats", 8, language.English)
	assert.Equal(t, "I've 8 cats.", msg)

	msg = srv.TranslatePlural("cats", -1, language.English) // negative and positive scales behave the same!
	assert.Equal(t, "I've one cat.", msg)

	msg = srv.TranslatePlural("cats", -2, language.English) // negative and positive scales behave the same!
	assert.Equal(t, "I've -2 cats.", msg)

	msg, err = srv.TranslatePluralMaybe("cats", -2, language.English)
	assert.NoError(t, err)
	assert.Equal(t, "I've -2 cats.", msg)

	msg = srv.TranslatePlural("cats", nil, language.English)
	assert.Equal(t, "I've <nil> cats.", msg) // invalid count

	msg, err = srv.TranslatePluralMaybe("cats", nil, language.English)
	assert.NoError(t, err)
	assert.Equal(t, "I've <nil> cats.", msg) // invalid count

	msg = srv.TranslatePlural("cats", "many", language.English)
	assert.Equal(t, "cats (count=many)", msg) // internal failed to translate plural!

	// overwrite Count
	msg = srv.TranslatePlural("cats", 8, language.English, i18n.Data{"Count": "too many"})
	assert.Equal(t, "I've too many cats.", msg)

	msg = srv.TranslatePlural("cats", 0, language.German)
	assert.Equal(t, "Ich habe 0 Katzen.", msg) // zero is not supported for CLDR German, "Ich habe keine Katze." is not possible!

	msg = srv.TranslatePlural("cats", 1, language.German)
	assert.Equal(t, "Ich habe eine Katze.", msg)

	msg = srv.TranslatePlural("cats", 2, language.German)
	assert.Equal(t, "Ich habe 2 Katzen.", msg)

	msg = srv.TranslatePlural("cats", 8, language.German)
	assert.Equal(t, "Ich habe 8 Katzen.", msg)

	msg = srv.TranslatePlural("cats", "viele", language.German)
	assert.Equal(t, "cats (count=viele)", msg) // internal failed to translate plural!

	msg, err = srv.TranslatePluralMaybe("cats", "viele", language.German)
	assert.Error(t, err)
	assert.Equal(t, "", msg) // empty string for errors!

	// overwrite Count
	msg = srv.TranslatePlural("cats", 8, language.German, i18n.Data{"Count": "zu viele"})
	assert.Equal(t, "Ich habe zu viele Katzen.", msg)

	// unknown language string
	tag := srv.ParseLang("xx")
	assert.Equal(t, language.English, tag)
	msg = srv.TranslatePlural("cats", 8, tag)
	assert.Equal(t, "I've 8 cats.", msg) // fall back to English

	// invalid specialized language string
	tag = srv.ParseLang("de-xx")
	msg = srv.TranslatePlural("cats", 8, tag)
	assert.Equal(t, "Ich habe 8 Katzen.", msg) // fall back to German

	// invalid language string
	tag = srv.ParseLang("§$%/%&/(/&%/)(")
	assert.Equal(t, language.English, tag)
	msg = srv.TranslatePlural("cats", 8, tag)
	assert.Equal(t, "I've 8 cats.", msg) // fall back to English

	// unknown
	msg = srv.TranslatePlural("this.key.will.never.exist", nil, language.English)
	assert.Equal(t, "this.key.will.never.exist (count=<nil>)", msg)

	msg, err = srv.TranslatePluralMaybe("this.key.will.never.exist", nil, language.English)
	assert.Error(t, err)
	assert.Equal(t, "", msg) // empty string for errors!
}

func TestI18nUndetermined(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n-undetermined"),
	})

	require.Error(t, err)
	assert.Equal(t, err, errors.New("undetermined language at index 1 in i18n message bundle: [en und]"))
}

func TestI18nUndeterminedDefaultLanguage(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage: language.Und, // Undetermined is disallowed
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})

	require.Error(t, err)
	assert.Equal(t, err, errors.New("undetermined language at index 0 in i18n message bundle: [und de en]"))
}

func TestI18nInvalidToml(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n-invalid"),
	})

	require.Error(t, err)
}

func TestI18nInexistantFolder(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage: language.AmericanEnglish,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n/this/folder/does/not/exist"),
	})

	require.Error(t, err)
}

func TestParseAcceptLanguage(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})

	require.NoError(t, err)

	tag := srv.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.German, tag)

	tag = srv.ParseAcceptLanguage("de-AT,en-US;q=0.7,en;q=0.3")
	assert.NotEqual(t, language.German, tag) // actual: de-u-rg-atzzzz

	// unknown language header
	tag = srv.ParseAcceptLanguage("xx,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.English, tag)
	msg := srv.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	// invalid specialized language string
	tag = srv.ParseAcceptLanguage("de-xx,en-US;q=0.7,en;q=0.3")
	msg = srv.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	// invalid language header
	tag = srv.ParseAcceptLanguage("§$%/%&/(/&%/)(")
	assert.Equal(t, language.English, tag)
	msg = srv.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}

func TestParseLanguage(t *testing.T) {
	srv, err := i18n.New(config.I18n{
		DefaultLanguage: language.English,
		BundleDirAbs:    filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/i18n"),
	})

	require.NoError(t, err)

	tag := srv.ParseLang("de")
	assert.Equal(t, language.German, tag)

	// unknown language string
	tag = srv.ParseLang("xx")
	assert.Equal(t, language.English, tag)
	msg := srv.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	// invalid specialized language string
	tag = srv.ParseLang("de-xx")
	msg = srv.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg) // fall back to German

	// invalid language string
	tag = srv.ParseLang("§$%/%&/(/&%/)(")
	assert.Equal(t, language.English, tag)
	msg = srv.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}
