package middleware

import (
	"DataCollector/internal/crawler"
	"DataCollector/internal/logger"

	"github.com/abadojack/whatlanggo"
)

// LanguageDetectorMiddleware creates a middleware that detects and sets the document language
// using the whatlanggo library (popular Go language detection library).
// This middleware will detect the language and update the document's Language field.
func LanguageDetectorMiddleware() crawler.Middleware {
	return func(doc *crawler.Document) (bool, error) {
		if doc.Content == nil {
			return false, nil
		}

		content, ok := doc.Content.(string)
		if !ok || content == "" {
			return false, nil
		}

		// Detect language using whatlanggo
		lang, confidence, _ := whatlanggo.DetectLanguage(content)
		langCode := mapWhatlangToCode(lang)
		doc.Language = langCode

		logger.Info("Detected language: %s (confidence: %.2f) for document: %s", langCode, confidence, doc.URL)

		return true, nil
	}
}

// LanguageDetectorWithFilter creates a middleware that detects language and filters
// to only allow specific languages.
func LanguageDetectorWithFilter(allowedLanguages ...string) crawler.Middleware {
	allowedSet := make(map[string]bool)
	for _, lang := range allowedLanguages {
		allowedSet[lang] = true
	}

	return func(doc *crawler.Document) (bool, error) {
		if doc.Content == nil {
			return false, nil
		}

		content, ok := doc.Content.(string)
		if !ok || content == "" {
			return false, nil
		}

		// Detect language using whatlanggo
		lang, confidence, _ := whatlanggo.DetectLanguage(content)
		langCode := mapWhatlangToCode(lang)
		doc.Language = langCode

		// Check if language is allowed
		if !allowedSet[langCode] {
			logger.Info("Document language '%s' not in allowed list, skipping: %s", langCode, doc.URL)
			return false, nil
		}

		logger.Info("Detected and accepted language: %s (confidence: %.2f) for document: %s", langCode, confidence, doc.URL)
		return true, nil
	}
}

// mapWhatlangToCode maps whatlanggo language to simple language code
func mapWhatlangToCode(lang whatlanggo.Language) string {
	switch lang {
	case whatlanggo.English:
		return "en"
	case whatlanggo.Persian:
		return "fa"
	case whatlanggo.Arabic:
		return "ar"
	case whatlanggo.Spanish:
		return "es"
	case whatlanggo.French:
		return "fr"
	case whatlanggo.German:
		return "de"
	case whatlanggo.Italian:
		return "it"
	case whatlanggo.Portuguese:
		return "pt"
	case whatlanggo.Russian:
		return "ru"
	case whatlanggo.Chinese:
		return "zh"
	case whatlanggo.Japanese:
		return "ja"
	case whatlanggo.Korean:
		return "ko"
	case whatlanggo.Dutch:
		return "nl"
	case whatlanggo.Polish:
		return "pl"
	case whatlanggo.Turkish:
		return "tr"
	case whatlanggo.Swedish:
		return "sv"
	case whatlanggo.Norwegian:
		return "no"
	case whatlanggo.Danish:
		return "da"
	case whatlanggo.Finnish:
		return "fi"
	case whatlanggo.Greek:
		return "el"
	case whatlanggo.Hebrew:
		return "he"
	case whatlanggo.Thai:
		return "th"
	case whatlanggo.Vietnamese:
		return "vi"
	case whatlanggo.Indonesian:
		return "id"
	case whatlanggo.Malay:
		return "ms"
	case whatlanggo.Hindi:
		return "hi"
	case whatlanggo.Urdu:
		return "ur"
	case whatlanggo.Romanian:
		return "ro"
	case whatlanggo.Czech:
		return "cs"
	case whatlanggo.Hungarian:
		return "hu"
	case whatlanggo.Bulgarian:
		return "bg"
	case whatlanggo.Ukrainian:
		return "uk"
	case whatlanggo.Serbian:
		return "sr"
	case whatlanggo.Croatian:
		return "hr"
	case whatlanggo.Slovak:
		return "sk"
	case whatlanggo.Slovenian:
		return "sl"
	case whatlanggo.Lithuanian:
		return "lt"
	case whatlanggo.Latvian:
		return "lv"
	case whatlanggo.Estonian:
		return "et"
	case whatlanggo.Icelandic:
		return "is"
	case whatlanggo.Galician:
		return "gl"
	case whatlanggo.Catalan:
		return "ca"
	case whatlanggo.Welsh:
		return "cy"
	case whatlanggo.Arabic:
		return "ar"
	default:
		return lang.String()
	}
}