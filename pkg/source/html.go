package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type HTMLDocument struct {
	path string
}

func htmlFromPath(path string) (*HTMLDocument, error) {
	return &HTMLDocument{
		path: path,
	}, nil
}

func (doc *HTMLDocument) Metadata() *Metadata {
	return nil
}

func (doc *HTMLDocument) Content() (string, error) {
	content, err := os.ReadFile(doc.path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", doc.path, err)
	}
	return string(content), nil
}

func (doc *HTMLDocument) CanonicalPath(base string) (string, error) {
	return filepath.Rel(base, strings.TrimSuffix(doc.path, ".tmpl"))
}
