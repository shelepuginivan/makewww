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

	flag.StringVar(&cfg.Dir, "dir", cwd, "directory to build the website from")
	flag.StringVar(&cfg.Output, "output", filepath.Join(cwd, "dist"), "output directory for the website")
	flag.BoolVar(&cfg.TransformDirs, "transform-dirs", false, "whether to generate directories with index.html")
	flag.Var(&cfg.TransformIgnore, "transform-ignore", "")

	flag.BoolVar(&cfg.Markdown.Extensions.Definitions, "md-ext-definitions", false, "whether to enable definition lists (PHP Markdown Extra)")
	flag.BoolVar(&cfg.Markdown.Extensions.Footnotes, "md-ext-footnotes", false, "whether to enable footnotes (PHP Markdown Extra)")
	flag.BoolVar(&cfg.Markdown.Extensions.Links, "md-ext-links", false, "whether to enable link detection extension")
	flag.BoolVar(&cfg.Markdown.Extensions.Strikethrough, "md-ext-strikethrough", false, "whether to enable strikethrough extension")
	flag.BoolVar(&cfg.Markdown.Extensions.TaskList, "md-ext-tasklist", false, "whether to enable tasklist extension")
	flag.BoolVar(&cfg.Markdown.Extensions.Tables, "md-ext-tables", false, "whether to enable tables extension")
	flag.BoolVar(&cfg.Markdown.Extensions.Typography, "md-ext-typography", false, "whether to enable smart typography extension")
	flag.BoolVar(&cfg.Markdown.Parser.Attributes, "md-parse-attrs", false, "whether to parse heading custom attributes")
	flag.BoolVar(&cfg.Markdown.Parser.AutoHeadingID, "md-parse-heading-id", false, "whether to enable auto heading IDs")
	flag.BoolVar(&cfg.Markdown.Render.HardWraps, "md-render-hardwraps", false, "whether to render newlines as <br>")
	flag.BoolVar(&cfg.Markdown.Render.Unsafe, "md-render-unsafe", false, "whether to render raw HTML")

	flag.Parse()

	if len(cfg.TransformIgnore) == 0 {
		cfg.TransformIgnore = append(cfg.TransformIgnore, "index.html", "404.html")
	}

	return &cfg, nil
}
