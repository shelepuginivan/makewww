package source

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

var frontmatterDelimiter = "---"

type MarkdownDocument struct {
	path       string
	sourceFile string
	isTemplate bool
	metadata   *Metadata
}

func markdownFromPath(path, sourceFile string, isTemplate bool) (*MarkdownDocument, error) {
	file, err := os.Open(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", sourceFile, err)
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
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
		metadata:   &metadata,
	}, nil
}

func (doc *MarkdownDocument) Metadata() *Metadata {
	return doc.metadata
}

func (doc *MarkdownDocument) Content() (string, error) {
	content, err := os.ReadFile(doc.sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", doc.sourceFile, err)
	}

	_, text, ok := bytes.Cut(content, []byte("\n---\n"))
	if !ok {
		return "", fmt.Errorf("failed to read content: %w", err)
	}

	return string(text), nil
}

func (doc *MarkdownDocument) Path() *Path {
	path := doc.path
	if doc.isTemplate {
		path = strings.TrimSuffix(doc.path, ".tmpl")
	}
	return &Path{strings.TrimSuffix(path, ".md") + ".html"}
}

func (doc *MarkdownDocument) IsTemplate() bool {
	return doc.isTemplate
}
