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

func TestServerProvidedMessages(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		// uses /app/web/messages by default
		// expect all messages file were loaded and the defaultLanguage matches the FIRST priority Tag.
		assert.Equal(t, s.Config.I18n.DefaultLanguage, s.Messages.Tags()[0])

		// expect all messages were loaded...
		files, err := os.ReadDir(s.Config.I18n.MessageFilesBaseDirAbs)
		require.NoError(t, err)

		msgFilesCount := 0

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
				continue
			}

			msgFilesCount++
		}

		if msgFilesCount == 0 {
			// all i18n messages were deleted, as the defaultLanguage is a tag itself, check for 1
			assert.Equal(t, len(s.Messages.Tags()), 1)
		} else {
			assert.Equal(t, len(s.Messages.Tags()), msgFilesCount)
		}

		msg := s.Messages.Translate("this.key.will.never.exist", s.Config.I18n.DefaultLanguage)
		assert.Equal(t, "this.key.will.never.exist", msg)
	})
}

// Note that all following tests use a special message directory within /internal/i18n/testdata.
// We do this to ensure we don't depend on your project specific i18n messages/configuration,
// that you would typically store within /web/messages.

func TestMessages(t *testing.T) {
	// Messages from /internal/i18n/testdata/messages
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})
	require.NoError(t, err)

	assert.Equal(t, language.English, messages.Tags()[0])
	assert.Equal(t, language.German, messages.Tags()[1])
	assert.Equal(t, 2, len(messages.Tags()))

	msg := messages.Translate("Test.Welcome", language.German, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	msg = messages.Translate("Test.Welcome", language.English, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = messages.Translate("Test.Welcome", language.Spanish, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	msg = messages.Translate("Test.Welcome", language.English, i18n.Data{"Name": "Franz"})
	assert.Equal(t, "Welcome Franz", msg)

	msg = messages.Translate("Test.Welcome", language.English)
	assert.Equal(t, "Welcome <no value>", msg)

	msg = messages.Translate("Test.Body", language.German)
	assert.Equal(t, "Das ist ein Test", msg)

	msg = messages.Translate("Test.Body", language.English)
	assert.Equal(t, "This is a test", msg)

	msg = messages.Translate("Test.Body", language.Spanish)
	assert.Equal(t, "This is a test", msg)

	msg = messages.Translate("Test.Invalid.Key.Does.Not.Exist", language.English)
	assert.Equal(t, "Test.Invalid.Key.Does.Not.Exist", msg)

	msg = messages.Translate("Test.Invalid.Key.Does.Not.Exist", language.German)
	assert.Equal(t, "Test.Invalid.Key.Does.Not.Exist", msg)

	msg = messages.Translate("Test.String.DE.only", language.English)
	assert.Equal(t, "Test.String.DE.only", msg)

	msg = messages.Translate("Test.String.DE.only", language.German)
	assert.Equal(t, "This key only exists in DE", msg)

	msg = messages.Translate("Test.String.EN.only", language.English)
	assert.Equal(t, "This key only exists in EN", msg)

	msg = messages.Translate("Test.String.EN.only", language.German)
	assert.Equal(t, "Test.String.EN.only", msg)

	// ensure language subvariants are supported
	deAt := messages.ParseLang("de-AT")
	assert.NotEqual(t, language.German, deAt)
	msg = messages.Translate("Test.Body", deAt)
	assert.Equal(t, "Das ist ein Test", msg)
}

func TestMessagesConcurrentUsage(t *testing.T) {
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(index int) {
			msg := messages.Translate("Test.Welcome", language.German, i18n.Data{"Name": fmt.Sprintf("%v", index)})
			assert.Equal(t, fmt.Sprintf("Guten Tag %v", index), msg)

			msg = messages.Translate("Test.Welcome", language.English, i18n.Data{"Name": fmt.Sprintf("%v", index)})
			assert.Equal(t, fmt.Sprintf("Welcome %v", index), msg)

			msg = messages.Translate("Test.Welcome", language.Spanish, i18n.Data{"Name": fmt.Sprintf("%v", index)})
			assert.Equal(t, fmt.Sprintf("Welcome %v", index), msg)

			msg = messages.Translate("Test.Welcome", language.English, i18n.Data{"Name": "Franz"})
			assert.Equal(t, "Welcome Franz", msg)

			msg = messages.Translate("Test.Welcome", language.English)
			assert.Equal(t, "Welcome <no value>", msg)

			msg = messages.Translate("Test.Body", language.German)
			assert.Equal(t, "Das ist ein Test", msg)

			msg = messages.Translate("Test.Body", language.English)
			assert.Equal(t, "This is a test", msg)

			msg = messages.Translate("Test.Body", language.Spanish)
			assert.Equal(t, "This is a test", msg)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestMessagesOtherDefault(t *testing.T) {
	// Messages from /internal/i18n/testdata/messages
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.German,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})
	require.NoError(t, err)

	assert.Equal(t, language.German, messages.Tags()[0])
	assert.Equal(t, language.English, messages.Tags()[1])
	assert.Equal(t, 2, len(messages.Tags()))
}

func TestMessagesInexistantDefault(t *testing.T) {
	// Messages from /internal/i18n/testdata/messages
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.Italian,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})
	require.NoError(t, err)

	assert.Equal(t, language.Italian, messages.Tags()[0])
	assert.Equal(t, language.German, messages.Tags()[1])
	assert.Equal(t, language.English, messages.Tags()[2])
	assert.Equal(t, 3, len(messages.Tags()))
}

func TestMessagesEmpty(t *testing.T) {
	// Messages from /internal/i18n/testdata/messages
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.Italian,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages-empty"),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(messages.Tags())) // the DefaultLanguage should still be set!
	assert.Equal(t, language.Italian, messages.Tags()[0])

	tag := messages.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.Italian, tag)

	msg := messages.Translate("no.test.key.exists", tag)
	assert.Equal(t, "no.test.key.exists", msg)

	msg = messages.Translate("no.test.key.exists", language.Ukrainian)
	assert.Equal(t, "no.test.key.exists", msg)
}

func TestMessagesSpecialized(t *testing.T) {
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.AmericanEnglish,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages-specialized"),
	})

	require.NoError(t, err)
	assert.Equal(t, 4, len(messages.Tags())) // specialized subvariant is default
	assert.Equal(t, language.AmericanEnglish, messages.Tags()[0])

	msg := messages.Translate("test.punchline", language.AmericanEnglish)
	assert.Equal(t, "I can has HUMOR?", msg)

	msg = messages.Translate("test.punchline", language.BritishEnglish)
	assert.Equal(t, "I can has HUMOUR?", msg)

	msg = messages.Translate("test.punchline", language.English)
	assert.Equal(t, "I can has HUMOR?", msg) // fall back to default

	msg = messages.Translate("test.punchline", language.German)
	assert.Equal(t, "Habe ich Humor?", msg) // jump to parsed Austrian German

	tag := messages.ParseAcceptLanguage("de-at,en-US;q=0.7,en;q=0.3") // explicit Austrian tag
	msg = messages.Translate("test.punchline", tag)
	assert.Equal(t, "Koan i Humor?", msg) // jump to parsed Austrian German
}

func TestReservedKeywordsResolve(t *testing.T) {
	// "reserved" keys:
	// "id", "description", "hash", "leftdelim", "rightdelim", "zero", "one", "two", "few", "many", "other"
	// see https://github.com/nicksnyder/go-i18n/blob/2180cd9f35b3e125cfe3773a6bf3ea483347f060/v2/i18n/message.go#L181

	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages-reserved"),
	})

	require.NoError(t, err)
	assert.Equal(t, 2, len(messages.Tags()))

	// english: single reserved word
	msg := messages.Translate("reserved1.zero", language.English)
	assert.Equal(t, "Zero", msg)

	msg = messages.Translate("reserved2.one", language.English)
	assert.Equal(t, "One", msg)

	msg = messages.Translate("reserved3.two", language.English)
	assert.Equal(t, "Two", msg)

	msg = messages.Translate("reserved4.few", language.English)
	assert.Equal(t, "Few", msg)

	msg = messages.Translate("reserved5.many", language.English)
	assert.Equal(t, "Many", msg)

	msg = messages.Translate("reserved6.other", language.English)
	assert.Equal(t, "Other", msg)

	msg = messages.Translate("reserved7.id", language.English)
	assert.Equal(t, "id", msg)

	msg = messages.Translate("reserved8.description", language.English)
	assert.Equal(t, "Description", msg)

	// german: single reserved word
	msg = messages.Translate("reserved1.zero", language.German)
	assert.Equal(t, "Null", msg)

	msg = messages.Translate("reserved2.one", language.German)
	assert.Equal(t, "Eins", msg)

	msg = messages.Translate("reserved3.two", language.German)
	assert.Equal(t, "Zwei", msg)

	msg = messages.Translate("reserved4.few", language.German)
	assert.Equal(t, "Wenig", msg)

	msg = messages.Translate("reserved5.many", language.German)
	assert.Equal(t, "Mehr", msg)

	msg = messages.Translate("reserved6.other", language.German)
	assert.Equal(t, "Andere", msg)

	msg = messages.Translate("reserved7.id", language.German)
	assert.Equal(t, "ID", msg)

	msg = messages.Translate("reserved8.description", language.German)
	assert.Equal(t, "Beschreibung", msg)

	// combined toml map: all reserved words
	msg = messages.Translate("reservedMap.zero", language.English)
	assert.Equal(t, "Zero", msg)

	msg = messages.Translate("reservedMap.one", language.English)
	assert.Equal(t, "One", msg)

	msg = messages.Translate("reservedMap.two", language.English)
	assert.Equal(t, "Two", msg)

	msg = messages.Translate("reservedMap.few", language.English)
	assert.Equal(t, "Few", msg)

	msg = messages.Translate("reservedMap.many", language.English)
	assert.Equal(t, "Many", msg)

	msg = messages.Translate("reservedMap.other", language.English)
	assert.Equal(t, "Other", msg)

	// plain toml: all reserved words
	msg = messages.Translate("reserved.plain.zero", language.English)
	assert.Equal(t, "Zero", msg)

	msg = messages.Translate("reserved.plain.one", language.English)
	assert.Equal(t, "One", msg)

	msg = messages.Translate("reserved.plain.two", language.English)
	assert.Equal(t, "Two", msg)

	msg = messages.Translate("reserved.plain.few", language.English)
	assert.Equal(t, "Few", msg)

	msg = messages.Translate("reserved.plain.many", language.English)
	assert.Equal(t, "Many", msg)

	msg = messages.Translate("reserved.plain.other", language.English)
	assert.Equal(t, "Other", msg)

	msg = messages.Translate("reserved.plain2.id", language.English)
	assert.Equal(t, "id", msg)

	msg = messages.Translate("reserved.plain2.description", language.English)
	assert.Equal(t, "Description", msg)

}

// func TestMessagesConditionalCount(t *testing.T) {
// 	messages, err := i18n.New(config.I18n{
// 		DefaultLanguage:        language.English,
// 		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages-plural"),
// 	})

// 	require.NoError(t, err)
// 	assert.Equal(t, 2, len(messages.Tags()))

// 	msg := messages.TranslateConditionalCount("cats", 0, language.AmericanEnglish)
// 	assert.Equal(t, "I don't have a cat.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 1, language.BritishEnglish)
// 	assert.Equal(t, "I've one cat.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 2, language.English)
// 	assert.Equal(t, "I've 2 cats.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 8, language.English)
// 	assert.Equal(t, "I've 8 cats.", msg)

// 	// overwrite Count
// 	msg = messages.TranslateConditionalCount("cats", 8, language.English, i18n.Data{"Count": "too many"})
// 	assert.Equal(t, "I've too many cats.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 0, language.German)
// 	assert.Equal(t, "Ich habe keine Katze.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 1, language.German)
// 	assert.Equal(t, "Ich habe eine Katze.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 2, language.German)
// 	assert.Equal(t, "Ich habe 2 Katzen.", msg)

// 	msg = messages.TranslateConditionalCount("cats", 8, language.German)
// 	assert.Equal(t, "Ich habe 8 Katzen.", msg)

// 	// overwrite Count
// 	msg = messages.TranslateConditionalCount("cats", 8, language.German, i18n.Data{"Count": "zu viele"})
// 	assert.Equal(t, "Ich habe zu viele Katzen.", msg)
// }

func TestMessagesUndetermined(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages-undetermined"),
	})

	require.Error(t, err)
	assert.Equal(t, err, errors.New("undetermined language at index 1 in message bundle: [en und]"))
}

func TestMessagesUndeterminedDefaultLanguage(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage:        language.Und, // Undetermined is disallowed
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})

	require.Error(t, err)
	assert.Equal(t, err, errors.New("undetermined language at index 0 in message bundle: [und de en]"))
}

func TestMessagesInvalidToml(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages-invalid"),
	})

	require.Error(t, err)
}

func TestMessagesInexistantFolder(t *testing.T) {
	_, err := i18n.New(config.I18n{
		DefaultLanguage:        language.AmericanEnglish,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages/this/folder/does/not/exist"),
	})

	require.Error(t, err)
}

func TestParseAcceptLanguage(t *testing.T) {
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})

	require.NoError(t, err)

	tag := messages.ParseAcceptLanguage("de,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.German, tag)

	tag = messages.ParseAcceptLanguage("de-AT,en-US;q=0.7,en;q=0.3")
	assert.NotEqual(t, language.German, tag) // actual: de-u-rg-atzzzz

	// unknown language header
	tag = messages.ParseAcceptLanguage("xx,en-US;q=0.7,en;q=0.3")
	assert.Equal(t, language.English, tag)
	msg := messages.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	// invalid specialized language string
	tag = messages.ParseAcceptLanguage("de-xx,en-US;q=0.7,en;q=0.3")
	msg = messages.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg)

	// invalid language header
	tag = messages.ParseAcceptLanguage("§$%/%&/(/&%/)(")
	assert.Equal(t, language.English, tag)
	msg = messages.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}

func TestParseLanguage(t *testing.T) {
	messages, err := i18n.New(config.I18n{
		DefaultLanguage:        language.English,
		MessageFilesBaseDirAbs: filepath.Join(util.GetProjectRootDir(), "/internal/i18n/testdata/messages"),
	})

	require.NoError(t, err)

	tag := messages.ParseLang("de")
	assert.Equal(t, language.German, tag)

	// unknown language string
	tag = messages.ParseLang("xx")
	assert.Equal(t, language.English, tag)
	msg := messages.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)

	// invalid specialized language string
	tag = messages.ParseLang("de-xx")
	msg = messages.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Guten Tag Hans", msg) // fall back to German

	// invalid language string
	tag = messages.ParseLang("§$%/%&/(/&%/)(")
	assert.Equal(t, language.English, tag)
	msg = messages.Translate("Test.Welcome", tag, i18n.Data{"Name": "Hans"})
	assert.Equal(t, "Welcome Hans", msg)
}
