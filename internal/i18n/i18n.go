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

// i18n, we expect the following:
// Your translation toml files should live within /web/i18n and are named according to their supported locale e.g. en.toml, de.toml or en-uk.toml, en-us.toml.
//
// All translation files should hold the same keys (there are no unique keys on a single bundle locale, apart from pluralization).
//
// The Service object is created and owned by the api.Server (s.I18n), you typically don't want to create your own object.
//
// Pluralization obeys to CLDR rules (https://cldr.unicode.org/index/cldr-spec/plural-rules).
// Some keywords are reserved for CLDR behaviour, templating and documentation: id, description, hash, leftdelim, rightdelim, zero, one, two, few, many, other
// See

// Service is your convenience object to call Translate/TranslatePlural and match languages according to your loaded translation bundle and its supported languages/locales.
type Service struct {
	bundle  *i18n.Bundle
	matcher language.Matcher
}

// Data should be used to pass your template data
type Data map[string]string

// New returns a new Service struct holding bundle and matcher with the settings of the given config
//
// Note that Service is typically created and owned by the api.Server (use it via s.I18n)
func New(config config.I18n) (*Service, error) {

	bundle := i18n.NewBundle(config.DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load all translation files in each language...
	files, err := os.ReadDir(config.BundleDirAbs)
	if err != nil {
		log.Err(err).Str("dir", config.BundleDirAbs).Msg("Failed to read i18n bundle directory")
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}

		// bundle.LoadMessageFile automatically guesses the language.Tag based on the filenames it encounters
		_, err := bundle.LoadMessageFile(filepath.Join(config.BundleDirAbs, file.Name()))
		if err != nil {
			log.Err(err).Str("file", file.Name()).Msg("Failed to load i18n message file")
			return nil, err
		}

	}

	tags := bundle.LanguageTags()

	for tagIndex, tag := range tags {
		// Undetermined languages are disallowed in our bundle.
		if tag == language.Und {
			err := fmt.Errorf("undetermined language at index %v in i18n message bundle: %v", tagIndex, tags)
			log.Err(err).Int("index", tagIndex).Str("tags", fmt.Sprintf("%v", tags)).Msg("Invalid i18n message bundle or default language.")
			return nil, err
		}
	}

	return &Service{
		bundle:  bundle,
		matcher: language.NewMatcher(tags),
	}, nil
}

// Translate your key into a localized string.
//
// Translate makes a lookup for the key in the current bundle with the specified language.
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
//
// Translate will not fail if a template value is missing "<no value>" will be inserted instead.
// Translate will also not fail if the key is not present. "{{key}}" will be returned instead.
func (m *Service) Translate(key string, lang language.Tag, data ...Data) string {
	msg, err := m.TranslateMaybe(key, lang, data...)

	if err != nil {
		log.Debug().Err(err).Str("key", key).Str("lang", lang.String()).Msg("Failed to translate")
		return key
	}

	return msg
}

// TranslateMaybe has the same sematics as Translate with the following exceptions:
// It exposes encountered errors (does not automatically log this error) and encountered errors may result in an empty "" string!
//
// This method may be useful for conditional translation rendering (if key is available, use that, else...).
func (m *Service) TranslateMaybe(key string, lang language.Tag, data ...Data) (string, error) {

	localizeConfig := &i18n.LocalizeConfig{
		MessageID: key,
	}

	if len(data) > 0 {
		localizeConfig.TemplateData = data[0]
	}

	return m.translateConfigurable(lang, localizeConfig)
}

// TranslatePlural translates a pluralized cldrKey into a localized string.
//
// TranslatePlural makes a lookup for the cldrKey (a base key holding CLDR keys like "one" and "other") in the current bundle with the specified language.
// This function should be used to conditionally show the pluralized form, controlled by the count param and according to the CLDR rules.
//
// Note that English and German only support .one and .other CLDR plural rules.
// See https://cldr.unicode.org/index/cldr-spec/plural-rules and https://www.unicode.org/cldr/cldr-aux/charts/28/supplemental/language_plural_rules.html
//
// If a language translation is not available the default language will be used.
// Additional data for templated strings can be passed as key value pairs with by passing an optional data map.
// The count param is automatically injected into this data map as stringified {{.Count}} and may be overwritten.
//
// TranslatePlural will not fail if a template value is missing "<no value>" will be inserted instead.
// TranslatePlural will also not fail if the cldrKey is not present. "{{cldrKey}} (count={{count}})" will be returned instead.
func (m *Service) TranslatePlural(cldrKey string, count interface{}, lang language.Tag, data ...Data) string {
	msg, err := m.TranslatePluralMaybe(cldrKey, count, lang, data...)
	if err != nil {
		log.Debug().Err(err).Str("count", fmt.Sprintf("%v", count)).Str("cldrKey", cldrKey).Str("lang", lang.String()).Msg("Failed to translate plural")
		return fmt.Sprintf("%s (count=%v)", cldrKey, count)
	}

	return msg
}

// TranslatePluralMaybe uses the same sematics as TranslatePlural with the following exceptions:
// It exposes encountered errors (does not automatically log this error) and encountered errors may result in an empty "" string!
//
// This method may be useful for conditional plural translation rendering (if key is available, use that, else...).
func (m *Service) TranslatePluralMaybe(cldrKey string, count interface{}, lang language.Tag, data ...Data) (string, error) {

	localizeConfig := &i18n.LocalizeConfig{
		MessageID:   cldrKey,
		PluralCount: count,
	}

	// We inject Count by default into our template data (for rare usecases you may overwrite it)
	templateData := make(Data)
	templateData["Count"] = fmt.Sprintf("%v", count)

	// If optional data was provided, merge them into the templateData map
	if len(data) > 0 {
		for k, v := range data[0] {
			templateData[k] = v
		}
	}

	localizeConfig.TemplateData = templateData

	return m.translateConfigurable(lang, localizeConfig)
}

// translateConfigurable is used internally for fully configurable translations according to our configured language precedence semantics (new Localizer per call).
func (m *Service) translateConfigurable(lang language.Tag, localizeConfig *i18n.LocalizeConfig) (string, error) {

	// We benchmarked precaching all known []i18n.NewLocalizer during initialization,
	// but it doesn't make a significant difference even with 10000 concurrent * 8 .Translate calls.
	// Thus we take the easy route and initialize a new localizer with each .Translate or .TranslatePlural call.
	localizer := i18n.NewLocalizer(m.bundle, lang.String())
	return localizer.Localize(localizeConfig)
}

// ParseAcceptLanguage takes the value of the Accept-Language header and returns
// the best matched language using the matcher.
func (m *Service) ParseAcceptLanguage(lang string) language.Tag {

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
func (m *Service) ParseLang(lang string) language.Tag {
	t, err := language.Parse(lang)
	if err != nil {
		log.Err(err).Str("lang", lang).Msg("Failed to parse language")
	}

	matchedTag, _, _ := m.matcher.Match(t)

	return matchedTag
}

// Tags returns the parsed and priority ordered []language.Tag (your config.DefaultLanguage will be on position 0)
func (m *Service) Tags() []language.Tag {
	return m.bundle.LanguageTags()
}
