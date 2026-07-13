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

func (res *Raw) CopyTo(w io.Writer) error {
	file, err := os.Open(res.sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open raw document: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("failed to write raw document: %w", err)
	}

	return nil
}

func (res *Raw) Content() ([]byte, error) {
	content, err := os.ReadFile(res.sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", res.sourceFile, err)
	}
	return content, nil
}

func (res *Raw) IsTemplate() bool {
	return res.isTemplate
}

func (res *Raw) Path() *Path {
	return &Path{res.path}
}
