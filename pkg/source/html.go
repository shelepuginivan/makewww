package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type HTMLDocument struct {
	path    string
	content string
}

func htmlFromPath(path string) (*HTMLDocument, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	return &HTMLDocument{
		path:    path,
		content: string(content),
	}, nil
}

func (doc *HTMLDocument) Metadata() *Metadata {
	return nil
}

func (doc *HTMLDocument) Content() string {
	return doc.content
}

func (doc *HTMLDocument) CanonicalPath(base string) (string, error) {
	return filepath.Rel(base, strings.TrimSuffix(doc.path, ".tmpl"))
}
