// Package config provides makewww configuration.
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type MarkdownExtensions struct {
	Definitions   bool
	Footnotes     bool
	Links         bool
	Strikethrough bool
	Tables        bool
	TaskList      bool
	Typography    bool
}

type MarkdownParser struct {
	Attributes    bool
	AutoHeadingID bool
}

type MarkdownRender struct {
	HardWraps bool
	Unsafe    bool
}

type Markdown struct {
	Extensions MarkdownExtensions
	Parser     MarkdownParser
	Render     MarkdownRender
}

type Config struct {
	Dir             string
	Output          string
	TransformDirs   bool
	TransformIgnore SliceFlags
	Markdown        Markdown
}

func Parse() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get cwd: %w", err)
	}

	var cfg Config

	flag.StringVar(&cfg.Dir, "src", cwd, "")
	flag.StringVar(&cfg.Output, "out", filepath.Join(cwd, "dist"), "")
	flag.BoolVar(&cfg.TransformDirs, "transform-dirs", false, "")
	flag.Var(&cfg.TransformIgnore, "transform-ignore", "")

	flag.BoolVar(&cfg.Markdown.Extensions.Definitions, "md-ext-definitions", false, "")
	flag.BoolVar(&cfg.Markdown.Extensions.Footnotes, "md-ext-footnotes", false, "")
	flag.BoolVar(&cfg.Markdown.Extensions.Links, "md-ext-links", false, "")
	flag.BoolVar(&cfg.Markdown.Extensions.Strikethrough, "md-ext-strikethrough", false, "")
	flag.BoolVar(&cfg.Markdown.Extensions.Tables, "md-ext-tables", false, "")
	flag.BoolVar(&cfg.Markdown.Extensions.TaskList, "md-ext-tasklist", false, "")
	flag.BoolVar(&cfg.Markdown.Extensions.Typography, "md-ext-typography", false, "")
	flag.BoolVar(&cfg.Markdown.Parser.Attributes, "md-parse-attrs", false, "")
	flag.BoolVar(&cfg.Markdown.Parser.AutoHeadingID, "md-parse-heading-id", false, "")
	flag.BoolVar(&cfg.Markdown.Render.HardWraps, "md-render-hardwraps", false, "")
	flag.BoolVar(&cfg.Markdown.Render.Unsafe, "md-render-unsafe", false, "")

	flag.Usage = usage
	flag.Parse()

	if len(cfg.TransformIgnore) == 0 {
		cfg.TransformIgnore = append(cfg.TransformIgnore, "index.html", "404.html")
	}

	return &cfg, nil
}

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), `makewww - extremely simple static site generator

General options:

  -src <dir>
        Root directory of the website. Defaults to current working directory.

  -out <dir>
        Output directory for the website. Defaults to ./dist

  -transform-dirs
        Whether to transform paths to directory-style structure. If set,
        paths like '/<name>.html' will be rewritten as '/<name>/index.html'.
        Disabled by default.

  -transform-ignore <base>
        When using -transform-dirs, don't transform paths matching '<base>'.
        Can be set multiple times. Defaults to ['index.html', '404.html'].


Markdown options:

  -md-ext-definitions
        Markdown extension: definition lists (part of PHP Markdown Extra).
        Disabled by default.

  -md-ext-footnotes
        Markdown extension: footnotes (part of PHP Markdown Extra).
        Disabled by default.

  -md-ext-links
        Markdown extension: automatic links (part of GitHub Flavored Markdown).
        Disabled by default.

  -md-ext-strikethrough
        Markdown extension: strikethrough (part of GitHub Flavored Markdown).
        Disabled by default.

  -md-ext-tables
        Markdown extension: tables (part of GitHub Flavored Markdown).
        Disabled by default.

  -md-ext-tasklist
        Markdown extension: task lists (part of GitHub Flavored Markdown).
        Disabled by default.

  -md-ext-typography
        Markdown extension: smart typography.
        Disabled by default.

  -md-parse-attrs
        Whether to enable custom attributes. Currently only headings support
        attributes. Disabled by default.

  -md-parse-heading-id
        Whether to automatically insert 'id' attribute in headings.
        Disabled by default.

  -md-render-hardwraps
        Whether to render line breaks as '<br>' HTML tags.
        Disabled by default.

  -md-render-unsafe
        Whether to enable raw (and potentially unsafe) HTML in Markdown.
        If this feature is disabled, raw HTML will be omitted from the
        resulting document, otherwise it is copied as is. Disabled by default.
`)
}
