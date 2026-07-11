package builder

import "github.com/shelepuginivan/makewww/pkg/source"

type GlobalContext struct {
	Documents []source.Document
}

type DocumentContext struct {
	Global   *GlobalContext
	Metadata *source.Metadata
	Path     *source.Path
}

type TemplateContext struct {
	Document *DocumentContext
	Content  string
}

func newDocumentContext(doc source.Document, global *GlobalContext) *DocumentContext {
	return &DocumentContext{
		Global:   global,
		Metadata: doc.Metadata(),
		Path:     doc.Path(),
	}
}

func newTemplateContext(content string, document *DocumentContext) *TemplateContext {
	return &TemplateContext{
		Document: document,
		Content:  content,
	}
}
