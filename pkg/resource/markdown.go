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

type Markdown struct {
	path       string
	sourceFile string
	isTemplate bool
	metadata   *Metadata
}

func NewMarkdown(path, sourceFile string, isTemplate bool) (*Markdown, error) {
	file, err := os.Open(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", sourceFile, err)
	}
	defer file.Close()

	metadata, err := parseFrontMatter(file)
	if err != nil {
		return nil, err
	}

	return &Markdown{
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
		metadata:   metadata,
	}, nil
}

func (res *Markdown) Metadata() *Metadata {
	return res.metadata
}

func (res *Markdown) Content() ([]byte, error) {
	content, err := os.ReadFile(res.sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", res.sourceFile, err)
	}

	if bytes.HasPrefix(content, []byte("---\n")) {
		var ok bool

		_, content, ok = bytes.Cut(content[4:], []byte("---\n"))
		if !ok {
			return nil, fmt.Errorf("failed to read: invalid metadata")
		}
	}

	return content, nil
}

func (res *Markdown) Path() *Path {
	p := res.path
	if res.isTemplate {
		p = strings.TrimSuffix(res.path, ".tmpl")
	}
	return &Path{strings.TrimSuffix(p, ".md") + ".html"}
}

func (res *Markdown) IsTemplate() bool {
	return res.isTemplate
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
