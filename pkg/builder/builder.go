// Package builder provides [Builder] struct that builds the website.
package builder

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/source"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

type Builder struct {
	cfg *config.Config
	src *source.Source
	out *Output
	md  goldmark.Markdown
}

func New(cfg *config.Config) (*Builder, error) {
	src, err := source.FromProjectRoot(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	out, err := outputFromRoot(cfg.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to create output: %w", err)
	}

	md := markdownParserFromConfig(&cfg.Markdown)

	return &Builder{
		cfg: cfg,
		src: src,
		out: out,
		md:  md,
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
		Config:    b.cfg,
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
	case *source.HTMLDocument:
		err = b.renderHTMLDocument(document, global)
	case *source.MarkdownDocument:
		err = b.renderMarkdownDocument(document, global)
	case *source.TemplateDocument:
		err = b.renderTemplateDocument(document, global)
	}

	return err
}

func (b *Builder) renderHTMLDocument(doc *source.HTMLDocument, global *GlobalContext) error {
	file, err := b.out.CreateOutputFile(b.outputPath(doc.Path()))
	if err != nil {
		return err
	}
	defer file.Close()

	return b.execTemplateIfNeeded(file, doc, newDocumentContext(doc, global))
}

func (b *Builder) renderMarkdownDocument(doc *source.MarkdownDocument, global *GlobalContext) error {
	documentCtx := newDocumentContext(doc, global)
	markdown := new(bytes.Buffer)

	err := b.execTemplateIfNeeded(markdown, doc, documentCtx)
	if err != nil {
		return err
	}

	html := new(bytes.Buffer)
	if err := b.md.Convert(markdown.Bytes(), html); err != nil {
		return err
	}

	tmpl, err := b.src.GetTemplate(doc.Metadata().Template)
	if err != nil {
		return err
	}

	f, err := b.out.CreateOutputFile(b.outputPath(doc.Path()))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, newTemplateContext(html.String(), documentCtx))
}

func (b *Builder) renderTemplateDocument(doc *source.TemplateDocument, global *GlobalContext) error {
	file, err := b.out.CreateOutputFile(doc.Path().Relative())
	if err != nil {
		return err
	}
	defer file.Close()

	return b.execTemplateIfNeeded(file, doc, newDocumentContext(doc, global))
}

func (b *Builder) copyRawFile(raw *source.Raw) error {
	f, err := b.out.CreateOutputFile(raw.Path().Relative())
	if err != nil {
		return err
	}
	defer f.Close()

	if err := raw.CopyTo(f); err != nil {
		return err
	}

	return nil
}

func (b *Builder) outputPath(path *source.Path) string {
	if b.cfg.TransformDirs && path.Stem() != "index" {
		return filepath.Join(path.RelativeNormalized(), "index.html")
	} else {
		return path.Relative()
	}
}

func (b *Builder) execTemplateIfNeeded(w io.Writer, doc source.Document, data any) error {
	content, err := doc.Content()
	if err != nil {
		return err
	}

	if !doc.IsTemplate() {
		_, err = fmt.Fprint(w, content)
		return err
	}

	tmpl, err := template.New("page").Parse(content)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

func markdownParserFromConfig(cfg *config.Markdown) goldmark.Markdown {
	var extensions []goldmark.Extender
	var parserOpts []parser.Option
	var renderOpts []renderer.Option

	if cfg.Extensions.Definitions {
		extensions = append(extensions, extension.DefinitionList)
	}
	if cfg.Extensions.Footnotes {
		extensions = append(extensions, extension.Footnote)
	}
	if cfg.Extensions.GFM {
		extensions = append(extensions, extension.GFM)
	}

	if cfg.Parser.Attributes {
		parserOpts = append(parserOpts, parser.WithAttribute())
	}
	if cfg.Parser.AutoHeadingID {
		parserOpts = append(parserOpts, parser.WithAutoHeadingID())
	}

	if cfg.Render.HardWraps {
		renderOpts = append(renderOpts, html.WithHardWraps())
	}
	if cfg.Render.Unsafe {
		renderOpts = append(renderOpts, html.WithUnsafe())
	}

	return goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOpts...),
		goldmark.WithRendererOptions(renderOpts...),
	)
}
