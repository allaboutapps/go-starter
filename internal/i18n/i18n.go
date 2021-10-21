package i18n

import (
	"os"
	"path/filepath"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

var (
	bundle  *i18n.Bundle
	matcher language.Matcher
)

func init() {
	config := config.DefaultServiceConfigFromEnv()
	defaultLanguage, err := language.Parse(config.I18n.DefaultLanguage)
	if err != nil {
		log.Error().Err(err).Str("lang", config.I18n.DefaultLanguage).Msg("could not parse default language tag")
		panic(err)
	}

	bundle = i18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	tags := []language.Tag{}
	for _, lang := range config.I18n.AvailableLanguages {
		t, err := language.Parse(lang)
		if err != nil {
			log.Error().Err(err).Str("lang", lang).Msg("could not parse language tag")
			panic(err)
		}

		tags = append(tags, t)
	}

	matcher = language.NewMatcher(tags)

	files, err := os.ReadDir(config.I18n.MessageFilesBaseDirAbs)
	if err != nil {
		log.Error().Str("dir", config.I18n.MessageFilesBaseDirAbs).Err(err).Msg("Failed to read messages directory on init")
		return
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}
		_, err := bundle.LoadMessageFile(filepath.Join(config.I18n.MessageFilesBaseDirAbs, file.Name()))
		if err != nil {
			log.Error().Str("file", file.Name()).Err(err).Msg("Failed to load message file")
		}
	}
}

type Data map[string]string

// T makes a lookup for the key in the available messages in the current bundle with the specified language.
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
// T will not fail if a template value is missing "<no value>" will be inserted instead.
// T will also not fail if the key is not present. "" will be returned instead.
func T(key string, lang language.Tag, data ...Data) string {
	localizer := i18n.NewLocalizer(bundle, lang.String())

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

// ParseAcceptLanguage takes the value of the Accept-Language header and returns
// the best matched language using the matcher created with the default and available
// language settings defined in the server config.
func ParseAcceptLanguage(lang string) language.Tag {

	// we deliberately ignore the error returned here, as it will be nil and the matcher will simply pick the default language
	// this allows us to skip any malformed Accept-Language headers without returning 500 errors to the client
	// additionally, we don't really care about the q-factor weighting or confidence, the first match will be picked (with a fallback to English)
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Err(err).Msg("Failed to parse accept language")
	}
	matchedTag, _, _ := matcher.Match(tags...)

	return matchedTag
}

// ParseAcceptLanguage parses the string as language tag and returns
// the best matched language using the matcher created with the default and available
// language settings defined in the server config.
func ParseLang(lang string) language.Tag {
	t, err := language.Parse(lang)
	if err != nil {
		log.Err(err).Msg("Failed to parse language")
	}

	matchedTag, _, _ := matcher.Match(t)

	return matchedTag
}
