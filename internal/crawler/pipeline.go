package crawler

import (
	"DataCollector/internal/logger"
)

// Middleware is a function that processes a document before it's saved.
// It can modify the document, skip it (by returning false), or pass it through.
type Middleware func(doc *Document) (bool, error)

// Pipeline manages a chain of middleware that process documents.
type Pipeline struct {
	middlewares []Middleware
}

// NewPipeline creates a new pipeline with the given middlewares.
func NewPipeline(middlewares ...Middleware) *Pipeline {
	return &Pipeline{
		middlewares: middlewares,
	}
}

// Use adds a middleware to the pipeline.
func (p *Pipeline) Use(middleware Middleware) {
	p.middlewares = append(p.middlewares, middleware)
}

// Process runs the document through all middleware in order.
// Returns true if the document should be saved, false if it should be skipped.
func (p *Pipeline) Process(doc *Document) (bool, error) {
	for i, middleware := range p.middlewares {
		shouldContinue, err := middleware(doc)
		if err != nil {
			logger.Error("Pipeline middleware %d error: %v", i, err)
			return false, err
		}
		if !shouldContinue {
			logger.Info("Document skipped by pipeline middleware %d: %s", i, doc.URL)
			return false, nil
		}
	}
	return true, nil
}

// ProcessSync is a synchronous version that processes and returns the result.
// It's useful for simple cases where you just need to check if a document should be saved.
func (p *Pipeline) ProcessSync(doc *Document) (bool, error) {
	return p.Process(doc)
}

// Built-in middleware functions

// FilterByLanguage creates a middleware that filters documents by language.
// If the document language doesn't match, it will be skipped.
func FilterByLanguage(allowedLanguages ...string) Middleware {
	return func(doc *Document) (bool, error) {
		for _, lang := range allowedLanguages {
			if doc.Language == lang {
				return true, nil
			}
		}
		return false, nil
	}
}

// FilterByMinLength creates a middleware that filters documents shorter than the minimum length.
func FilterByMinLength(minLength int) Middleware {
	return func(doc *Document) (bool, error) {
		if content, ok := doc.Content.(string); ok {
			if len(content) < minLength {
				return false, nil
			}
		}
		return true, nil
	}
}

// FilterByMaxLength creates a middleware that filters documents longer than the maximum length.
func FilterByMaxLength(maxLength int) Middleware {
	return func(doc *Document) (bool, error) {
		if content, ok := doc.Content.(string); ok {
			if len(content) > maxLength {
				return false, nil
			}
		}
		return true, nil
	}
}

// TransformContent creates a middleware that transforms the document content.
// The transform function receives the current content and returns the new content.
func TransformContent(transform func(content any) (any, error)) Middleware {
	return func(doc *Document) (bool, error) {
		newContent, err := transform(doc.Content)
		if err != nil {
			return false, err
		}
		doc.Content = newContent
		return true, nil
	}
}

// AddMetadata creates a middleware that adds metadata to the document.
func AddMetadata(key string, value any) Middleware {
	return func(doc *Document) (bool, error) {
		if doc.Metadata == nil {
			doc.Metadata = make(map[string]any)
		}
		doc.Metadata[key] = value
		return true, nil
	}
}

// LogDocument creates a middleware that logs document information.
func LogDocument() Middleware {
	return func(doc *Document) (bool, error) {
		logger.Info("Processing document: id=%s, url=%s, title=%s, lang=%s", doc.ID, doc.URL, doc.Title, doc.Language)
		return true, nil
	}
}

// ValidateDocument creates a middleware that validates the document.
// Returns false if the document is invalid (empty title, content, etc.).
func ValidateDocument() Middleware {
	return func(doc *Document) (bool, error) {
		if doc.Title == "" {
			return false, nil
		}
		if doc.Content == nil {
			return false, nil
		}
		if content, ok := doc.Content.(string); ok && content == "" {
			return false, nil
		}
		return true, nil
	}
}