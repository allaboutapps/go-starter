package i18n

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

var (
	_bundle   *Bundle
	_matcher  *Matcher
	_initOnce sync.Once
)

type Bundle struct {
	bundle *i18n.Bundle
}

type Matcher struct {
	matcher language.Matcher
}

// InitPackage initializes the global bundle and matcher with the default values from the environment.
func InitGlobalBundleMatcher(config config.I18n) {
	_initOnce.Do(func() {
		var err error
		_bundle, err = NewBundle(config)
		if err != nil {
			log.Err(err).Msg("Failed to initialize global i18n bundle")
		}

		_matcher, err = NewMatcher(config)
		if err != nil {
			log.Err(err).Msg("Failed to initialize global i18n matcher")
		}
	})
}

// NewBundle returns a new bundle with the settings of the given config.
func NewBundle(config config.I18n) (*Bundle, error) {
	bundle := i18n.NewBundle(config.DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	files, err := os.ReadDir(config.MessageFilesBaseDirAbs)
	if err != nil {
		log.Err(err).Str("dir", config.MessageFilesBaseDirAbs).Msg("Failed to read messages directory on init")
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}
		_, err := bundle.LoadMessageFile(filepath.Join(config.MessageFilesBaseDirAbs, file.Name()))
		if err != nil {
			log.Err(err).Str("file", file.Name()).Msg("Failed to load message file")
			return nil, err
		}
	}

	return &Bundle{
		bundle: bundle,
	}, nil
}

// NewBundle returns a new matcher with the settings of the given config.
func NewMatcher(config config.I18n) (*Matcher, error) {
	return &Matcher{
		matcher: language.NewMatcher(config.AvailableLanguages),
	}, nil
}

type Data map[string]string

// T makes a lookup for the key in the available messages in the current bundle with the specified language.
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
// T will not fail if a template value is missing "<no value>" will be inserted instead.
// T will also not fail if the key is not present. "" will be returned instead.
func (b *Bundle) T(key string, lang language.Tag, data ...Data) string {
	localizer := i18n.NewLocalizer(b.bundle, lang.String())

	localizeConfig := &i18n.LocalizeConfig{
		MessageID: key,
	}
	if len(data) > 0 {
		localizeConfig.TemplateData = data[0]
	}

	msg, err := localizer.Localize(localizeConfig)
	if err != nil {
		log.Err(err).Str("key", key).Msg("Failed to localize message")
		return ""
	}

	return msg
}

// T makes a lookup for the key in the available messages in the global bundle with the specified language.
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
// T will not fail if a template value is missing "<no value>" will be inserted instead.
// T will also not fail if the key is not present. "" will be returned instead.
func T(key string, lang language.Tag, data ...Data) string {
	return _bundle.T(key, lang, data...)
}

// ParseAcceptLanguage takes the value of the Accept-Language header and returns
// the best matched language using the matcher.
func (m *Matcher) ParseAcceptLanguage(lang string) language.Tag {

	// we deliberately ignore the error returned here, as it will be nil and the matcher will simply pick the default language
	// this allows us to skip any malformed Accept-Language headers without returning 500 errors to the client
	// additionally, we don't really care about the q-factor weighting or confidence, the first match will be picked (with a fallback to English)
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Err(err).Str("lang", lang).Msg("Failed to parse accept language")
	}
	matchedTag, _, _ := m.matcher.Match(tags...)

	return matchedTag
}

// ParseAcceptLanguage takes the value of the Accept-Language header and returns
// the best matched language using the global matcher created with the default and available
// language settings defined in the server config.
func ParseAcceptLanguage(lang string) language.Tag {
	return _matcher.ParseAcceptLanguage(lang)
}

// ParseLang parses the string as language tag and returns
// the best matched language using the matcher.
func (m *Matcher) ParseLang(lang string) language.Tag {
	t, err := language.Parse(lang)
	if err != nil {
		log.Err(err).Msg("Failed to parse language")
	}

	matchedTag, _, _ := m.matcher.Match(t)

	return matchedTag
}

// ParseLang parses the string as language tag and returns
// the best matched language using the matcher created with the default and available
// language settings defined in the server config.
func ParseLang(lang string) language.Tag {
	return _matcher.ParseLang(lang)
}
