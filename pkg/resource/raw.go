package resource

import (
	"fmt"
	"io"
	"os"
)

type Raw struct {
	path       string
	sourceFile string
	isTemplate bool
}

func NewRaw(path, sourceFile string, isTemplate bool) *Raw {
	return &Raw{
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
	}
}

func (doc *Raw) CopyTo(w io.Writer) error {
	file, err := os.Open(doc.sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open raw document: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("failed to write raw document: %w", err)
	}

	return nil
}

func (doc *Raw) Content() (string, error) {
	content, err := os.ReadFile(doc.sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", doc.sourceFile, err)
	}
	return string(content), nil
}

func (doc *Raw) IsTemplate() bool {
	return doc.isTemplate
}

func (doc *Raw) Path() *Path {
	return &Path{doc.path}
}
