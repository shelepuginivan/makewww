package builder

import (
	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/source"
)

type GlobalContext struct {
	Config    *config.Config
	Documents []source.Document
}

type DocumentContext struct {
	Global   *GlobalContext
	Metadata *source.Metadata
	Path     *source.Path
}

type TemplateDocumentContext struct {
	Metadata *source.Metadata
	Path     *source.Path
}

type TemplateContext struct {
	Global   *GlobalContext
	Document *TemplateDocumentContext
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
		Global: document.Global,
		Document: &TemplateDocumentContext{
			Metadata: document.Metadata,
			Path:     document.Path,
		},
		Content: content,
	}
}
