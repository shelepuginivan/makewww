// Package config provides makewww configuration.
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type MarkdownExtensions struct {
	Definitions bool
	Footnotes   bool
	GFM         bool
	Typography  bool
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
	Dir           string
	Output        string
	TransformDirs bool
	Markdown      Markdown
}

func Parse() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get cwd: %w", err)
	}

	var cfg Config

	flag.StringVar(&cfg.Dir, "dir", cwd, "directory to build the website from")
	flag.StringVar(&cfg.Output, "output", filepath.Join(cwd, "dist"), "output directory for the website")
	flag.BoolVar(&cfg.TransformDirs, "transform-dirs", false, "whether to generate directories with index.html")

	flag.BoolVar(&cfg.Markdown.Extensions.Definitions, "md-ext-definitions", false, "whether to enable definition lists (PHP Markdown Extra)")
	flag.BoolVar(&cfg.Markdown.Extensions.Footnotes, "md-ext-footnotes", false, "whether to enable footnotes (PHP Markdown Extra)")
	flag.BoolVar(&cfg.Markdown.Extensions.GFM, "md-ext-gfm", false, "whether to enable GFM (GitHub Flavored Markdown) extensions")
	flag.BoolVar(&cfg.Markdown.Extensions.Typography, "md-ext-typography", false, "whether to enable smart typography extension")
	flag.BoolVar(&cfg.Markdown.Parser.Attributes, "md-parse-attrs", false, "whether to parse heading custom attributes")
	flag.BoolVar(&cfg.Markdown.Parser.AutoHeadingID, "md-parse-heading-id", false, "whether to enable auto heading IDs")
	flag.BoolVar(&cfg.Markdown.Render.HardWraps, "md-render-hardwraps", false, "whether to render newlines as <br>")
	flag.BoolVar(&cfg.Markdown.Render.Unsafe, "md-render-unsafe", false, "whether to render raw HTML")

	flag.Parse()

	return &cfg, nil
}
