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
	tags    []language.Tag
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

		// bundle.LoadMessageFile automatically guesses the language.Tag based on the filename it encounters
		_, err := bundle.LoadMessageFile(filepath.Join(config.MessageFilesBaseDirAbs, file.Name()))
		if err != nil {
			log.Err(err).Str("file", file.Name()).Msg("Failed to load message file")
			return nil, err
		}

	}

	// Build up the language priority array and check the set default language was loaded from the messages files.
	parsedTags := bundle.LanguageTags()

	tags := []language.Tag{config.DefaultLanguage}
	defaultLanguageLoaded := false

	for _, tag := range parsedTags {

		if tag == config.DefaultLanguage {
			defaultLanguageLoaded = true
			continue
		}

		// push next language priorities
		tags = append(tags, tag)
	}

	if !defaultLanguageLoaded {
		err = fmt.Errorf("default language '%v' is missing in messages", config.DefaultLanguage.String())

		log.Err(err).
			Str("DefaultLanguage", config.DefaultLanguage.String()).
			Str("MessageFilesBaseDirAbs", config.MessageFilesBaseDirAbs).
			Msg("Failed to find or load default language within MessageFilesBaseDirAbs dir")
		return nil, err
	}

	return &Messages{
		bundle:  bundle,
		tags:    tags,
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
