package source

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type MarkdownDocument struct {
	path     string
	content  string
	metadata *Metadata
}

func markdownFromPath(path string) (*MarkdownDocument, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	frontmatter, text, ok := bytes.Cut(content, []byte("\n---\n"))
	var metadata *Metadata
	if ok {
		metadata = &Metadata{}
		if err := yaml.Unmarshal(frontmatter, metadata); err != nil {
			return nil, fmt.Errorf("failed to parse metadata: %w", err)
		}
	}

	return &MarkdownDocument{
		path:     path,
		content:  string(text),
		metadata: metadata,
	}, nil
}

func (doc *MarkdownDocument) Metadata() *Metadata {
	return doc.metadata
}

func (doc *MarkdownDocument) Content() string {
	return doc.content
}

func (m *MarkdownDocument) CanonicalPath(base string) (string, error) {
	return filepath.Rel(base, strings.TrimSuffix(m.path, ".md.tmpl")+".html")
}
