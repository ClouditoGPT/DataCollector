package lang

import (
	"github.com/pemistahl/lingua-go"
)

// Result holds the language detection result
type Result struct {
	Lang       string
	Confidence float64
}

// Detect detects the language of the given text
func Detect(text string) Result {
	detector := lingua.NewLanguageDetectorBuilder().FromAllLanguages().Build()
	lang, ok := detector.DetectLanguageOf(text)
	if !ok {
		return Result{Lang: "unknown", Confidence: 0}
	}
	return Result{
		Lang:       string(lang.IsoCode639_1()),
		Confidence: 1.0,
	}
}
