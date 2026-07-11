package source

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

var frontmatterDelimiter = "---"

type MarkdownDocument struct {
	path     string
	metadata *Metadata
}

func markdownFromPath(path string) (*MarkdownDocument, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}

	buffer := bytes.NewBuffer(nil)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == frontmatterDelimiter {
			break
		}

		buffer.WriteString(line)
		buffer.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata Metadata
	if err := yaml.Unmarshal(buffer.Bytes(), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &MarkdownDocument{
		path:     path,
		metadata: &metadata,
	}, nil
}

func (doc *MarkdownDocument) Metadata() *Metadata {
	return doc.metadata
}

func (doc *MarkdownDocument) Content() (string, error) {
	content, err := os.ReadFile(doc.path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", doc.path, err)
	}

	_, text, ok := bytes.Cut(content, []byte("\n---\n"))
	if !ok {
		return "", fmt.Errorf("failed to read content: %w", err)
	}

	return string(text), nil
}

func (doc *MarkdownDocument) CanonicalPath(base string) (string, error) {
	return filepath.Rel(base, strings.TrimSuffix(doc.path, ".md.tmpl")+".html")
}
