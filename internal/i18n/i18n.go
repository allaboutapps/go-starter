package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

// i18n implementation assumptions:
//   Your messages should live within /web/messages and are named according to their supported locale e.g. en.toml, de.toml or en-uk.toml, en-us.toml
//   All message files hold the same keys (there are no unique keys on a single message file)
//   The Messages object is created and owned by the api.Server (s.Messages), you typically don't want to create your own object.

// Messages is your convience object to call T (Translate) and match languages according to your loaded message bundle and its supported languages/locales.
type Messages struct {
	bundle  *i18n.Bundle
	matcher language.Matcher
}

// New returns a new Messages struct holding bundle and matcher with the settings of the given config
// Note that messages is typically created and owned by the api.Server (use it via s.Messages)
func New(config config.I18n) (*Messages, error) {

	bundle := i18n.NewBundle(config.DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load all message files in each language...
	files, err := os.ReadDir(config.MessageFilesBaseDirAbs)
	if err != nil {
		log.Err(err).Str("dir", config.MessageFilesBaseDirAbs).Msg("Failed to read messages directory")
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}

		// bundle.LoadMessageFile automatically guesses the language.Tag based on the filenames it encounters
		_, err := bundle.LoadMessageFile(filepath.Join(config.MessageFilesBaseDirAbs, file.Name()))
		if err != nil {
			log.Err(err).Str("file", file.Name()).Msg("Failed to load message file")
			return nil, err
		}

	}

	tags := bundle.LanguageTags()

	for tagIndex, tag := range tags {
		// Undetermined languages are disallowed in our bundle.
		if tag == language.Und {
			err := fmt.Errorf("undetermined language at index %v in message bundle: %v", tagIndex, tags)
			log.Err(err).Int("index", tagIndex).Str("tags", fmt.Sprintf("%v", tags))
			return nil, err
		}
	}

	return &Messages{
		bundle:  bundle,
		matcher: language.NewMatcher(tags),
	}, nil
}

type Data map[string]string

// T makes a lookup for the key in the available messages in the current bundle with the specified language.
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
// T will not fail if a template value is missing "<no value>" will be inserted instead.
// T will also not fail if the key is not present. key will be returned instead.
func (m *Messages) T(key string, lang language.Tag, data ...Data) string {

	localizer := i18n.NewLocalizer(m.bundle, lang.String())

	localizeConfig := &i18n.LocalizeConfig{
		MessageID: key,
	}
	if len(data) > 0 {
		localizeConfig.TemplateData = data[0]
	}

	msg, err := localizer.Localize(localizeConfig)
	if err != nil {
		log.Err(err).Str("key", key).Str("lang", lang.String()).Msg("Failed to localize message")
		return key
	}

	return msg
}

// ParseAcceptLanguage takes the value of the Accept-Language header and returns
// the best matched language using the matcher.
func (m *Messages) ParseAcceptLanguage(lang string) language.Tag {

	// we deliberately ignore the error returned here, as it will be nil and the matcher will simply pick the default language
	// this allows us to skip any malformed Accept-Language headers without returning 500 errors to the client
	// additionally, we don't really care about the q-factor weighting or confidence, the first match will be picked (with a fallback to config.DefaultLanguage)
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Err(err).Str("lang", lang).Msg("Failed to parse accept language")
	}
	matchedTag, _, _ := m.matcher.Match(tags...)

	return matchedTag
}

// ParseLang parses the string as language tag and returns
// the best matched language using the matcher.
func (m *Messages) ParseLang(lang string) language.Tag {
	t, err := language.Parse(lang)
	if err != nil {
		log.Err(err).Str("lang", lang).Msg("Failed to parse language")
	}

	matchedTag, _, _ := m.matcher.Match(t)

	return matchedTag
}

// Tags returns the parsed and priority ordered []language.Tag (your config.DefaultLanguage will be on position 0)
func (m *Messages) Tags() []language.Tag {
	return m.bundle.LanguageTags()
}
