package resource

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

func NewHTML(path, sourceFile string, isTemplate bool) (*HTMLDocument, error) {
	return &HTMLDocument{
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
	}, nil
}

func (doc *HTMLDocument) Content() ([]byte, error) {
	content, err := os.ReadFile(doc.sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", doc.sourceFile, err)
	}
	return content, nil
}

func (doc *HTMLDocument) Path() *Path {
	p := doc.path
	if doc.isTemplate {
		p = strings.TrimSuffix(doc.path, ".tmpl")
	}
	return &Path{p}
}

func (doc *HTMLDocument) IsTemplate() bool {
	return doc.isTemplate
}
