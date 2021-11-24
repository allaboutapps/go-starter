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

// Messages is your convenience object to call T (Translate) and match languages according to your loaded message bundle and its supported languages/locales.
type Messages struct {
	bundle  *i18n.Bundle
	matcher language.Matcher
}

// New returns a new Messages struct holding bundle and matcher with the settings of the given config
//
// Note that Messages is typically created and owned by the api.Server (use it via s.Messages)
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
			log.Err(err).Int("index", tagIndex).Str("tags", fmt.Sprintf("%v", tags)).Msg("Invalid message bundle or default language.")
			return nil, err
		}
	}

	return &Messages{
		bundle:  bundle,
		matcher: language.NewMatcher(tags),
	}, nil
}

type Data map[string]string

// Translate your key into a localized string.
//
// Translate makes a lookup for the key in the current bundle with the specified language.
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
//
// Translate will not fail if a template value is missing "<no value>" will be inserted instead.
// Translate will also not fail if the key is not present. key will be returned instead.
func (m *Messages) Translate(key string, lang language.Tag, data ...Data) string {
	localizer := m.getLocalizer(lang)

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

// // TranslateConditionalCount translates a pluralized zeroOneOtherParentKey into a localized string.
// //
// // TranslateConditionalCount makes a lookup for the oneOtherParentKey (a base key holding the child keys "zero", "one" and "other") in the current bundle with the specified language.
// // This function should be used to conditionally show the pluralized form, controlled by the pluralCount param.
// // If a language translation is not available the default language will be used.
// // Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
// //
// // TranslateConditionalCount will not fail if a template value is missing "<no value>" will be inserted instead.
// // TranslateConditionalCount will also not fail if the key is not present. key will be returned instead.
// func (m *Messages) TranslateConditionalCount(zeroOneOtherParentKey string, count uint, lang language.Tag, data ...Data) string {
// 	localizer := m.getLocalizer(lang)
// 	localizeConfig := &i18n.LocalizeConfig{}

// 	switch count {
// 	case 0:
// 		localizeConfig.MessageID = fmt.Sprintf("%s.zero", zeroOneOtherParentKey)

// 	case 1:
// 		localizeConfig.MessageID = fmt.Sprintf("%s.one", zeroOneOtherParentKey)

// 	default:
// 		localizeConfig.MessageID = fmt.Sprintf("%s.other", zeroOneOtherParentKey)
// 	}

// 	fmt.Println(localizeConfig.MessageID)

// 	// We inject Count by default into our template data (for rare usecases you may overwrite it)
// 	templateData := make(Data)
// 	templateData["Count"] = fmt.Sprintf("%d", count)

// 	// If optional data was provided, merge them into the templateData map
// 	if len(data) > 0 {
// 		for k, v := range data[0] {
// 			templateData[k] = v
// 		}
// 	}

// 	localizeConfig.TemplateData = templateData

// 	msg, err := localizer.Localize(localizeConfig)
// 	if err != nil {
// 		log.Err(err).Uint("count", count).Str("key", localizeConfig.MessageID).Str("lang", lang.String()).Msg("Failed to localize conditional count message")
// 		return localizeConfig.MessageID
// 	}

// 	return msg
// }

// TranslateConfigurable exposes the real i18n.LocalizeConfig used internally
// and allows for fully configurable translations according to its semantics.
func (m *Messages) TranslateConfigurable(lang language.Tag, localizeConfig *i18n.LocalizeConfig) string {
	localizer := m.getLocalizer(lang)

	msg, err := localizer.Localize(localizeConfig)
	if err != nil {
		log.Err(err).Str("localizeConfig.MessageID", localizeConfig.MessageID).Str("lang", lang.String()).Msg("Failed to localize configured message")
		return localizeConfig.MessageID
	}

	return msg
}

// getLocalizer is a helper to return a new localizer for a potentially unknown language tag (best match from m.bundle.LanguageTags())
func (m *Messages) getLocalizer(lang language.Tag) *i18n.Localizer {

	// We benchmarked precaching i18n.NewLocalizer during initialization,
	// but it doesn't make a significant difference even with 10000 concurrent * 8 .Translate calls.
	// Thus we take the easy route and initialize a new localizer with each .Translate or .TranslateConditionalCount call.
	return i18n.NewLocalizer(m.bundle, lang.String())
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
