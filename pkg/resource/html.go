package resource

import (
	"fmt"
	"os"
	"strings"
)

type HTML struct {
	path       string
	sourceFile string
	isTemplate bool
}

func NewHTML(path, sourceFile string, isTemplate bool) (*HTML, error) {
	return &HTML{
		path:       path,
		sourceFile: sourceFile,
		isTemplate: isTemplate,
	}, nil
}

func (res *HTML) Content() ([]byte, error) {
	content, err := os.ReadFile(res.sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", res.sourceFile, err)
	}
	return content, nil
}

func (res *HTML) Path() *Path {
	p := res.path
	if res.isTemplate {
		p = strings.TrimSuffix(res.path, ".tmpl")
	}
	return &Path{p}
}

func (res *HTML) IsTemplate() bool {
	return res.isTemplate
}
