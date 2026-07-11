package source

import (
	"fmt"
	"os"
	"strings"
)

type TemplateDocument struct {
	path       string
	sourceFile string
}

func templateFromPath(path, sourceFile string) (*TemplateDocument, error) {
	return &TemplateDocument{
		path:       path,
		sourceFile: sourceFile,
	}, nil
}

func (doc *TemplateDocument) Metadata() *Metadata {
	return &Metadata{}
}

func (doc *TemplateDocument) Content() (string, error) {
	content, err := os.ReadFile(doc.sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", doc.sourceFile, err)
	}
	return string(content), nil
}

func (doc *TemplateDocument) Path() *Path {
	return &Path{strings.TrimSuffix(doc.path, ".tmpl")}
}
