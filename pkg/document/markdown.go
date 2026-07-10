package document

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

type Markdown struct {
	path string
}

func MarkdownFromPath(path string) *Markdown {
	return &Markdown{path: path}
}

func (m *Markdown) Render(w io.Writer) error {
	tmpl, err := template.ParseFiles(m.path)
	if err != nil {
		return fmt.Errorf("failed to parse markdown document: %w", err)
	}

	if err := tmpl.Execute(w, struct{}{}); err != nil {
		return fmt.Errorf("failed to render markdown document: %w", err)
	}

	return nil
}

func (m *Markdown) CanonicalPath(base string) (string, error) {
	return filepath.Rel(base, strings.TrimSuffix(m.path, ".md.tmpl")+".html")
}
