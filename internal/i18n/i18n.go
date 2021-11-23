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

// TODO rename package to messages?
// TODO txt regarding assuption all i18n keys are available in all languages and no fallback if key in one language is not available.
// TODO txt regarding strict naming convention of messages files.

type Messages struct {
	bundle  *i18n.Bundle
	matcher language.Matcher
}

// New returns a new Messages struct holding bundle and matcher with the settings of the given config
func New(config config.I18n) (*Messages, error) {

	bundle := i18n.NewBundle(config.DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load all message files in each language...
	files, err := os.ReadDir(config.MessageFilesBaseDirAbs)
	if err != nil {
		log.Err(err).Str("dir", config.MessageFilesBaseDirAbs).Msg("Failed to read messages directory on init")
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
			return nil, fmt.Errorf("undetermined language at pos %v in message bundle", tagIndex)
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

// Tags returns the parsed and priority ordered []language.Tag (config.DefaultLanguage will be on position 0)
func (m *Messages) Tags() []language.Tag {
	return m.bundle.LanguageTags()
}
