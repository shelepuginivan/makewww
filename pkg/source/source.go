package source

import (
	"fmt"
	"io/fs"
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

func (src *Source) GetDocuments() ([]Document, error) {
	var documents []Document

	err := filepath.Walk(src.ContentDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var d Document

		switch {
		case strings.HasSuffix(path, ".html.tmpl"):
			d, err = htmlFromPath(path)
		case strings.HasSuffix(path, ".md.tmpl"):
			d, err = markdownFromPath(path)
		default:
			return nil
		}

		if err != nil {
			return err
		}

		documents = append(documents, d)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	return documents, nil
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
