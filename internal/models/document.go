package models

type DocumentType string

const (
	ArticleDocument     DocumentType = "article"
	ChatDocument        DocumentType = "chat"
	CodeDocument        DocumentType = "code"
	InstructionDocument DocumentType = "instruction"
	QADocument          DocumentType = "qa"
)

type Document struct {
	ID       string         `json:"id"`
	Source   string         `json:"source"`
	Type     DocumentType   `json:"type"`
	Language string         `json:"language"`
	URL      string         `json:"url"`
	Title    string         `json:"title"`
	Content  any            `json:"content"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
