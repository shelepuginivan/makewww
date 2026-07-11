package source

import (
	"fmt"
	"os"
	"strings"
)

type HTMLDocument struct {
	path       string
	sourceFile string
}

func htmlFromPath(path, sourceFile string) (*HTMLDocument, error) {
	return &HTMLDocument{
		path:       path,
		sourceFile: sourceFile,
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
	return &Path{strings.TrimSuffix(doc.path, ".tmpl")}
}
