package document

import (
	"fmt"
	"io"
	"text/template"
)

type Markdown struct {
	path string
}

func MarkdownFromPath(path string) *Markdown {
	return &Markdown{path: path}
}

func (d *Markdown) Render(w io.Writer) error {
	tmpl, err := template.ParseFiles(d.path)
	if err != nil {
		return fmt.Errorf("failed to parse markdown document: %w", err)
	}

	if err := tmpl.Execute(w, struct{}{}); err != nil {
		return fmt.Errorf("failed to render markdown document: %w", err)
	}

	return nil
}
