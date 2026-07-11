package source

import (
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
	"strings"
	"text/template"
)

type Source struct {
	root string
}

func FromProjectRoot(root string) *Source {
	return &Source{root: root}
}

func (src *Source) ContentDir() string {
	return filepath.Join(src.root, "content")
}

func (src *Source) TemplatesDir() string {
	return filepath.Join(src.root, "templates")
}

func (src *Source) Documents() iter.Seq2[Document, error] {
	return func(yield func(Document, error) bool) {
		filepath.Walk(src.ContentDir(), func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				if !yield(nil, err) {
					return fs.SkipAll
				}
				return nil
			}

			if info.IsDir() {
				return nil
			}

			var doc Document

			switch {
			case strings.HasSuffix(path, ".html.tmpl"):
				doc, err = htmlFromPath(path)
			case strings.HasSuffix(path, ".md.tmpl"):
				doc, err = markdownFromPath(path)
			default:
				return nil
			}

			if err != nil {
				if !yield(nil, err) {
					return fs.SkipAll
				}
				return nil
			}

			if !yield(doc, nil) {
				return fs.SkipAll
			}

			return nil
		})
	}
}

func (src *Source) GetTemplate(path string) (*template.Template, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("template path must be absolute")
	}

	tmplPath := filepath.Join(src.TemplatesDir(), path)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return tmpl, nil
}
