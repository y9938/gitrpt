// Package i18n provides internationalization support for gitrpt.
// translations are loaded from files for other languages.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed messages/*.json
var messagesFS embed.FS

// Default language
const defaultLanguage = "en"

// supportedLanguages maps language codes to their language.Tag
var supportedLanguages = map[string]language.Tag{
	"en": language.English,
	"ru": language.Russian,
}

// I18n manages internationalization state
type I18n struct {
	lang      string
	localizer *i18n.Localizer
	bundle    *i18n.Bundle
}

// Messages - all translatable messages with English as default.
// These are extracted by goi18n extract command.
// NOTE: Version is not here - it's defined in main.go as a build-time variable
var Messages = struct {
	Name             *i18n.Message
	Description      *i18n.Message
	Usage            *i18n.Message
	FlagLang         *i18n.Message
	FlagHelp         *i18n.Message
	FlagVersion      *i18n.Message
	ErrorUnknownLang *i18n.Message
	LanguagePriority *i18n.Message
}{
	Name: &i18n.Message{
		ID:          "name",
		Description: "Tool name",
		Other:       "gitrpt",
	},
	Description: &i18n.Message{
		ID:          "description",
		Description: "Tool description",
		Other:       "Git activity reporter",
	},
	Usage: &i18n.Message{
		ID:          "usage",
		Description: "Usage syntax",
		Other:       "Usage: {{.Name}} [options]",
	},
	FlagLang: &i18n.Message{
		ID:          "flag_lang",
		Description: "Language flag description",
		Other:       "Language (en, ru)",
	},
	FlagHelp: &i18n.Message{
		ID:          "flag_help",
		Description: "Help flag description",
		Other:       "Show help message",
	},
	FlagVersion: &i18n.Message{
		ID:          "flag_version",
		Description: "Version flag description",
		Other:       "Show version",
	},
	ErrorUnknownLang: &i18n.Message{
		ID:          "error_unknown_lang",
		Description: "Unknown language error",
		Other:       "Unknown language: {{.Lang}}",
	},
	LanguagePriority: &i18n.Message{
		ID:          "language_priority",
		Description: "Language priority explanation",
		Other:       "Language priority: 1) --lang flag, 2) GITRPT_LANG env var, 3) default (en)",
	},
}

// New creates a new I18n instance with the specified language.
// - English (default) is embedded in code via DefaultMessage
// - Only load translation files for non-English languages
// - goi18n extract will generate active.en.json from Messages struct
func New(lang string) (*I18n, error) {
	if _, ok := supportedLanguages[lang]; !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	// Create bundle with English as default
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translation files for non-default languages only
	// English is provided by DefaultMessage in code
	if lang != defaultLanguage {
		if _, err := bundle.LoadMessageFileFS(messagesFS, fmt.Sprintf("messages/active.%s.json", lang)); err != nil {
			// Translation file not found is OK - we'll use English fallback
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to load %s messages: %w", lang, err)
			}
		}
	}

	// Create localizer with fallback chain:
	// 1. Requested language
	// 2. Default language (English - via DefaultMessage)
	localizer := i18n.NewLocalizer(bundle, lang, defaultLanguage)

	return &I18n{
		lang:      lang,
		localizer: localizer,
		bundle:    bundle,
	}, nil
}

// MustNew creates a new I18n instance and panics on error.
func MustNew(lang string) *I18n {
	i, err := New(lang)
	if err != nil {
		panic(err)
	}
	return i
}

// DetectLanguage determines the language from flag or environment.
// Priority: 1) flag value, 2) GITRPT_LANG env var, 3) default (en)
func DetectLanguage(flagLang string) string {
	if flagLang != "" {
		if _, ok := supportedLanguages[flagLang]; ok {
			return flagLang
		}
	}

	if envLang := os.Getenv("GITRPT_LANG"); envLang != "" {
		if _, ok := supportedLanguages[envLang]; ok {
			return envLang
		}
	}

	return defaultLanguage
}

// Lang returns the current language code.
func (i *I18n) Lang() string {
	return i.lang
}

// localize is the core localization method using DefaultMessage pattern.
// This ensures English always works even without translation files.
func (i *I18n) localize(msg *i18n.Message, templateData map[string]interface{}) string {
	result, err := i.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      msg.ID,
		DefaultMessage: msg,
		TemplateData:   templateData,
	})
	if err != nil {
		// Fallback to English text from DefaultMessage
		return msg.Other
	}
	return result
}

// Name returns the tool name.
func (i *I18n) Name() string {
	return i.localize(Messages.Name, nil)
}

// Description returns the tool description.
func (i *I18n) Description() string {
	return i.localize(Messages.Description, nil)
}

// Usage returns the usage string.
func (i *I18n) Usage() string {
	return i.localize(Messages.Usage, map[string]interface{}{"Name": i.localize(Messages.Name, nil)})
}

// FlagLang returns the lang flag description.
func (i *I18n) FlagLang() string {
	return i.localize(Messages.FlagLang, nil)
}

// FlagHelp returns the help flag description.
func (i *I18n) FlagHelp() string {
	return i.localize(Messages.FlagHelp, nil)
}

// FlagVersion returns the version flag description.
func (i *I18n) FlagVersion() string {
	return i.localize(Messages.FlagVersion, nil)
}

// ErrorUnknownLang returns the unknown language error message.
func (i *I18n) ErrorUnknownLang(lang string) string {
	return i.localize(Messages.ErrorUnknownLang, map[string]interface{}{"Lang": lang})
}

// LanguagePriority returns the language priority explanation.
func (i *I18n) LanguagePriority() string {
	return i.localize(Messages.LanguagePriority, nil)
}

// AddLanguage adds a new supported language at runtime.
// Useful for adding languages without modifying the package.
func AddLanguage(code string, tag language.Tag) {
	supportedLanguages[code] = tag
}

// SupportedLanguages returns a list of supported language codes.
func SupportedLanguages() []string {
	langs := make([]string, 0, len(supportedLanguages))
	for lang := range supportedLanguages {
		langs = append(langs, lang)
	}
	return langs
}
