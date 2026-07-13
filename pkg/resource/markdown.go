package resource

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const frontmatterDelimiter = "---"

type MarkdownDocument struct {
	path       string
	sourceFile string
	isTemplate bool
	metadata   *Metadata
}

func NewMarkdown(path, sourceFile string, isTemplate bool) (*MarkdownDocument, error) {
	file, err := os.Open(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", sourceFile, err)
	}
	defer file.Close()

	metadata, err := parseFrontMatter(file)
	if err != nil {
		return nil, err
	}

	return &MarkdownDocument{
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
		metadata:   metadata,
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

	if bytes.HasPrefix(content, []byte("---\n")) {
		var ok bool

		_, content, ok = bytes.Cut(content[4:], []byte("---\n"))
		if !ok {
			return "", fmt.Errorf("failed to read: invalid metadata")
		}
	}

	return string(content), nil
}

func (doc *MarkdownDocument) Path() *Path {
	p := doc.path
	if doc.isTemplate {
		p = strings.TrimSuffix(doc.path, ".tmpl")
	}
	return &Path{strings.TrimSuffix(p, ".md") + ".html"}
}

func (doc *MarkdownDocument) IsTemplate() bool {
	return doc.isTemplate
}

func parseFrontMatter(r io.Reader) (*Metadata, error) {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return &Metadata{}, nil
	}

	firstLine := scanner.Text()
	if firstLine != frontmatterDelimiter {
		return &Metadata{}, nil
	}

	yamlBuffer := new(bytes.Buffer)
	for scanner.Scan() {
		line := scanner.Text()
		if line == frontmatterDelimiter {
			break
		}

		yamlBuffer.WriteString(line)
		yamlBuffer.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	meta, err := metadataFromYAML(yamlBuffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return meta, nil
}
