package document

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

type HTML struct {
	path string
}

func HTMLFromPath(path string) *HTML {
	return &HTML{path: path}
}

func (h *HTML) Render(w io.Writer) error {
	tmpl, err := template.ParseFiles(h.path)
	if err != nil {
		return fmt.Errorf("failed to parse html document: %w", err)
	}

	if err := tmpl.Execute(w, struct{}{}); err != nil {
		return fmt.Errorf("failed to render html document: %w", err)
	}

	return nil
}

func (h *HTML) CanonicalPath(base string) (string, error) {
	return filepath.Rel(base, strings.TrimSuffix(h.path, ".tmpl"))
}
