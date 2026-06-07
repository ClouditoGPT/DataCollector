package middleware

import (
	"DataCollector/internal/crawler"
	"DataCollector/internal/lang"
	"DataCollector/internal/logger"
	"DataCollector/internal/models"
)

// LanguageDetectorMiddleware creates a middleware that detects and sets the document language
// using the lang library (language detection library).
// This middleware will detect the language and update the document's Language field.
func LanguageDetectorMiddleware() crawler.Middleware {
	return func(doc *models.Document) (bool, error) {
		if doc.Content == nil {
			return false, nil
		}

		content, ok := doc.Content.(string)
		if !ok || content == "" {
			return false, nil
		}

		// Detect language using lang
		result := lang.Detect(content)
		langCode := result.Lang
		doc.Language = langCode

		logger.Info("Detected language: %s (confidence: %.2f) for document: %s", langCode, result.Confidence, doc.URL)

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

	return func(doc *models.Document) (bool, error) {
		if doc.Content == nil {
			return false, nil
		}

		content, ok := doc.Content.(string)
		if !ok || content == "" {
			return false, nil
		}

		// Detect language using lang
		result := lang.Detect(content)
		langCode := result.Lang
		doc.Language = langCode

		// Check if language is allowed
		if !allowedSet[langCode] {
			logger.Info("Document language '%s' not in allowed list, skipping: %s", langCode, doc.URL)
			return false, nil
		}

		logger.Info("Detected and accepted language: %s (confidence: %.2f) for document: %s", langCode, result.Confidence, doc.URL)
		return true, nil
	}
}
