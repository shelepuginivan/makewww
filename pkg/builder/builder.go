package builder

import (
	"fmt"

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
	return nil
}

func (b *Builder) renderMarkdownDocument(doc *source.MarkdownDocument) error {
	return nil
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
