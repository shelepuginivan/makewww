package builder

import (
	"fmt"
	"strings"
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

func New(cfg *config.Config) *Builder {
	src := source.FromProjectRoot(cfg.Dir)
	dist := dist.FromRoot(cfg.Output)
	parser := goldmark.New()

	return &Builder{
		src:    src,
		dist:   dist,
		parser: parser,
	}
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

	for _, doc := range documents {
		if err := b.renderDocument(doc); err != nil {
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

func (b *Builder) renderDocument(doc source.Document) error {
	var err error

	switch document := doc.(type) {
	case *source.HTMLDocument:
		err = b.renderHTMLDocument(document)
	case *source.MarkdownDocument:
		err = b.renderMarkdownDocument(document)
	}

	return err
}

func (b *Builder) renderHTMLDocument(doc *source.HTMLDocument) error {
	content, err := doc.Content()
	if err != nil {
		return err
	}

	path, err := doc.CanonicalPath(b.src.ContentDir())
	if err != nil {
		return err
	}

	tmpl, err := template.New("page").Parse(content)
	if err != nil {
		return err
	}

	file, err := b.dist.CreateOutputFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, struct{}{})
}

func (b *Builder) renderMarkdownDocument(doc *source.MarkdownDocument) error {
	content, err := doc.Content()
	if err != nil {
		return err
	}

	dest, err := doc.CanonicalPath(b.src.ContentDir())
	if err != nil {
		return err
	}

	tmpl, err := template.New("page").Parse(content)
	if err != nil {
		return err
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, struct{}{}); err != nil {
		return err
	}

	tmpl, err = b.src.GetTemplate(doc.Metadata().Template)
	if err != nil {
		return err
	}

	f, err := b.dist.CreateOutputFile(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, struct{}{})
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
