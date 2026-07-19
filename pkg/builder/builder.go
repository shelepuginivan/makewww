// Package builder provides [Builder] struct that builds the website.
package builder

import (
	"log"
	"path/filepath"
	"slices"

	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/resource"
	"github.com/shelepuginivan/makewww/pkg/source"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

type Builder struct {
	cfg *config.Config
	md  goldmark.Markdown
}

func New(cfg *config.Config) (*Builder, error) {
	md := markdownParserFromConfig(&cfg.Markdown)

	return &Builder{
		cfg: cfg,
		md:  md,
	}, nil
}

func (b *Builder) Build(src *source.Source, out *Output) error {
	components, err := src.Components()
	if err != nil {
		return err
	}

	layouts, err := src.Layouts()
	if err != nil {
		return err
	}

	resources, err := src.Resources()
	if err != nil {
		return err
	}

	global := &GlobalContext{
		Config:    b.cfg,
		Resources: resources,
	}

	pipeline := &Pipeline{
		global:     global,
		components: components,
		layouts:    layouts,
		md:         b.md,
	}

	for _, res := range resources {
		log.Println(res.Path().Relative())

		file, err := out.CreateOutputFile(b.outputPath(res))
		if err != nil {
			return err
		}
		defer file.Close()

		err = pipeline.Process(res, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) outputPath(res resource.Resource) string {
	path := res.Path()

	switch res.(type) {
	case *resource.Raw:
		return path.Relative()
	}

	if b.cfg.TransformDirs && !slices.Contains(b.cfg.TransformIgnore, path.Base()) {
		return filepath.Join(path.RelativeNormalized(), "index.html")
	} else {
		return path.Relative()
	}
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
	if cfg.Extensions.Links {
		extensions = append(extensions, extension.Linkify)
	}
	if cfg.Extensions.Strikethrough {
		extensions = append(extensions, extension.Strikethrough)
	}
	if cfg.Extensions.TaskList {
		extensions = append(extensions, extension.Strikethrough)
	}
	if cfg.Extensions.Tables {
		extensions = append(extensions, extension.Table)
	}
	if cfg.Extensions.Typography {
		extensions = append(extensions, extension.Typographer)
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
