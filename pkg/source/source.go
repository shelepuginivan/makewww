package source

import (
	"fmt"
	"path/filepath"
	"text/template"
)

const (
	ContentDir   = "content"
	TemplatesDir = "templates"
)

type Source struct {
	root string
}

func FromProjectRoot(root string) *Source {
	return &Source{root: root}
}

func (src *Source) GetTemplate(path string) (*template.Template, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("template path must be absolute")
	}

	tmpl, err := template.ParseFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return tmpl, nil
}
