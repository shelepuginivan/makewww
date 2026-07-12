package source

import (
	"fmt"
	"os"
	"strings"
)

type HTMLDocument struct {
	path       string
	sourceFile string
	isTemplate bool
}

func htmlFromPath(path, sourceFile string, isTemplate bool) (*HTMLDocument, error) {
	return &HTMLDocument{
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
	}, nil
}

func (doc *HTMLDocument) Metadata() *Metadata {
	return &Metadata{}
}

func (doc *HTMLDocument) Content() (string, error) {
	content, err := os.ReadFile(doc.sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", doc.sourceFile, err)
	}
	return string(content), nil
}

func (doc *HTMLDocument) Path() *Path {
	path := doc.path
	if doc.isTemplate {
		path = strings.TrimSuffix(doc.path, ".tmpl")
	}
	return &Path{path}
}

func (doc *HTMLDocument) IsTemplate() bool {
	return doc.isTemplate
}
