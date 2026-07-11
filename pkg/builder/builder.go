package builder

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/dist"
	"github.com/shelepuginivan/makewww/pkg/source"
	"github.com/yuin/goldmark"
)

type Builder struct {
	src    *source.Source
	dist   *dist.Dist
	parser goldmark.Markdown
}

func New(cfg *config.Config) (*Builder, error) {
	src, err := source.FromProjectRoot(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	dist := dist.FromRoot(cfg.Output)
	parser := goldmark.New()

	return &Builder{
		src:    src,
		dist:   dist,
		parser: parser,
	}, nil
}

func (b *Builder) Build() error {
	documents, err := b.src.Documents()
	if err != nil {
		return fmt.Errorf("failed to get documents: %w", err)
	}

	rawFiles, err := b.src.RawFiles()
	if err != nil {
		return fmt.Errorf("failed to get raw files: %w", err)
	}

	global := &GlobalContext{
		Documents: documents,
	}

	for _, doc := range documents {
		if err := b.renderDocument(doc, global); err != nil {
			return fmt.Errorf("failed to render: %w", err)
		}
	}

	for _, raw := range rawFiles {
		if err := b.copyRawFile(raw); err != nil {
			return fmt.Errorf("failed to copy raw file: %w", err)
		}
	}

	return nil
}

func (b *Builder) renderDocument(doc source.Document, global *GlobalContext) error {
	var err error

	switch document := doc.(type) {
	case *source.TemplateDocument:
		err = b.renderTemplateDocument(document, global)
	case *source.MarkdownDocument:
		err = b.renderMarkdownDocument(document, global)
	}

	return err
}

func (b *Builder) renderTemplateDocument(doc *source.TemplateDocument, global *GlobalContext) error {
	content, err := doc.Content()
	if err != nil {
		return err
	}

	tmpl, err := template.New("page").Parse(content)
	if err != nil {
		return err
	}

	file, err := b.dist.CreateOutputFile(doc.Path().Relative())
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, newDocumentContext(doc, global))
}

func (b *Builder) renderMarkdownDocument(doc *source.MarkdownDocument, global *GlobalContext) error {
	content, err := doc.Content()
	if err != nil {
		return err
	}

	tmpl, err := template.New("page").Parse(content)
	if err != nil {
		return err
	}

	documentCtx := newDocumentContext(doc, global)

	markdown := new(bytes.Buffer)
	if err := tmpl.Execute(markdown, documentCtx); err != nil {
		return err
	}

	html := new(bytes.Buffer)
	if err := b.parser.Convert(markdown.Bytes(), html); err != nil {
		return err
	}

	tmpl, err = b.src.GetTemplate(doc.Metadata().Template)
	if err != nil {
		return err
	}

	f, err := b.dist.CreateOutputFile(doc.Path().Relative())
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, newTemplateContext(html.String(), documentCtx))
}

func (b *Builder) copyRawFile(raw *source.Raw) error {
	dest, err := raw.CanonicalPath(b.src.ContentDir())
	if err != nil {
		return err
	}

	f, err := b.dist.CreateOutputFile(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := raw.Render(f); err != nil {
		return err
	}

	return nil
}
