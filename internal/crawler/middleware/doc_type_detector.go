package middleware

import (
	"DataCollector/internal/crawler"
	"DataCollector/internal/logger"
	"DataCollector/internal/models"
	"strings"
)

// DocumentTypeDetectorMiddleware creates a middleware that detects and sets the document type
// based on content analysis. It uses heuristic patterns to classify documents into:
// article, chat, code, instruction, or qa.
func DocumentTypeDetectorMiddleware() crawler.Middleware {
	return func(doc *models.Document) (bool, error) {
		if doc.Content == nil {
			return false, nil
		}

		content, ok := doc.Content.(string)
		if !ok || content == "" {
			return false, nil
		}

		docType := models.DocumentType(detectDocumentType(content))
		doc.Type = docType

		logger.Info("Detected document type: %s for document: %s", docType, doc.URL)

		return true, nil
	}
}

// detectDocumentType analyzes the content and returns the most likely document type.
func detectDocumentType(content string) string {
	lowerContent := strings.ToLower(content)

	// Count indicators for each type
	scores := map[string]int{
		"code":        countCodeIndicators(content),
		"chat":        countChatIndicators(lowerContent),
		"qa":          countQAIndicators(lowerContent),
		"instruction": countInstructionIndicators(lowerContent),
		"article":     0,
	}

	// Article is the default/fallback, but boost it if it has paragraph structure
	if isArticleLike(content) {
		scores["article"] = 1
	}

	// Find the type with the highest score
	maxType := "article"
	maxScore := scores[maxType]

	for docType, score := range scores {
		if score > maxScore {
			maxScore = score
			maxType = docType
		}
	}

	return maxType
}

// countCodeIndicators counts programming/code patterns in the content.
func countCodeIndicators(content string) int {
	count := 0

	// Code block markers
	codeMarkers := []string{"```", "~~~", "<code>", "</code>", "<pre>", "</pre>"}
	for _, marker := range codeMarkers {
		count += strings.Count(content, marker)
	}

	// Common programming keywords
	codeKeywords := []string{
		"function", "const ", "let ", "var ", "import ", "export ",
		"class ", "def ", "return ", "if (", "for (", "while (",
		"public ", "private ", "static ", "void ", "int ", "string ",
		"print(", "console.log", "printf(", "System.out",
		"#include", "package ", "namespace ",
		"SELECT ", "INSERT ", "UPDATE ", "DELETE ", "CREATE TABLE",
		"<html>", "<div", "<span", "<script", "<style",
	}
	for _, keyword := range codeKeywords {
		if strings.Contains(content, keyword) {
			count += 2
		}
	}

	// Brackets and braces typical in code
	brackets := []string{"{", "}", "(", ")", "[", "]"}
	for _, b := range brackets {
		count += strings.Count(content, b) / 10
	}

	// Semicolons (common in many languages)
	count += strings.Count(content, ";") / 5

	return count
}

// countChatIndicators counts dialogue/chat patterns in the content.
func countChatIndicators(content string) int {
	count := 0

	// Common chat/dialogue patterns
	chatPatterns := []string{
		"user:", "assistant:", "human:", "ai:", "bot:",
		"question:", "answer:", "q:", "a:",
		"you:", "me:", "customer:", "support:",
		"<user>", "</user>", "<assistant>", "</assistant>",
		"<human>", "</human>", "<ai>", "</ai>",
	}

	for _, pattern := range chatPatterns {
		count += strings.Count(content, pattern) * 3
	}

	// Dialogue markers
	dialogueMarkers := []string{"\"", "«", "»", "—", "–"}
	for _, marker := range dialogueMarkers {
		count += strings.Count(content, marker) / 5
	}

	return count
}

// countQAIndicators counts question-answer patterns in the content.
func countQAIndicators(content string) int {
	count := 0

	// Question markers
	questionMarkers := []string{
		"q:", "question:", "faq", "frequently asked",
		"how to", "what is", "why does", "when should",
	}
	for _, marker := range questionMarkers {
		if strings.Contains(content, marker) {
			count += 3
		}
	}

	// Count question marks (but not too many, as that might be chat)
	questionMarks := strings.Count(content, "?")
	if questionMarks > 2 && questionMarks < 20 {
		count += questionMarks
	}

	// Answer indicators
	answerMarkers := []string{
		"a:", "answer:", "solution:", "response:",
	}
	for _, marker := range answerMarkers {
		count += strings.Count(content, marker) * 2
	}

	return count
}

// countInstructionIndicators counts step-by-step instruction patterns.
func countInstructionIndicators(content string) int {
	count := 0

	// Numbered list patterns
	numberedPatterns := []string{
		"step 1", "step 2", "step 3",
		"1.", "2.", "3.", "4.", "5.",
		"first,", "second,", "third,", "finally,",
		"1)", "2)", "3)",
	}
	for _, pattern := range numberedPatterns {
		count += strings.Count(content, pattern) * 2
	}

	// Imperative verbs common in instructions
	imperativeVerbs := []string{
		"click", "select", "open", "close", "type", "enter",
		"press", "choose", "navigate", "install", "configure",
		"follow these", "make sure", "ensure that",
		"to begin", "to start", "next,", "then,",
	}
	for _, verb := range imperativeVerbs {
		if strings.Contains(content, verb) {
			count += 2
		}
	}

	return count
}

// isArticleLike checks if the content has typical article structure.
func isArticleLike(content string) bool {
	// Articles typically have multiple paragraphs
	paragraphs := strings.Split(content, "\n\n")
	if len(paragraphs) < 3 {
		return false
	}

	// Check for heading-like patterns
	headingPatterns := 0
	for _, para := range paragraphs {
		trimmed := strings.TrimSpace(para)
		if len(trimmed) > 0 && len(trimmed) < 100 {
			// Short lines that could be headings
			words := strings.Fields(trimmed)
			if len(words) <= 10 && !strings.HasSuffix(trimmed, ".") {
				headingPatterns++
			}
		}
	}

	return headingPatterns >= 2
}
