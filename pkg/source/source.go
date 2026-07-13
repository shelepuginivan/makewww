// Package source provides methods for working with website source code.
package source

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/shelepuginivan/makewww/pkg/resource"
	"github.com/shelepuginivan/makewww/pkg/tmplfn"
)

const (
	componentsDir = "components"
	contentDir    = "content"
	templatesDir  = "templates"
)

type Source struct {
	root *os.Root
}

func FromProjectRoot(root string) (*Source, error) {
	r, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}

	return &Source{root: r}, nil
}

func (src *Source) Resources() ([]resource.Resource, error) {
	var resources []resource.Resource

	content := filepath.Join(src.root.Name(), contentDir)

	err := filepath.Walk(content, func(sourceFile string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		path, err := filepath.Rel(content, sourceFile)
		if err != nil {
			return err
		}

		res, err := resource.New(path, sourceFile)
		if err != nil {
			return err
		}

		resources = append(resources, res)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (src *Source) GetTemplate(path string) (string, error) {
	path = filepath.Join(templatesDir, path)
	content, err := src.root.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}
	return string(content), nil
}

func (src *Source) GetComponents() (*template.Template, error) {
	pattern := filepath.Join(src.root.Name(), componentsDir, "*.tmpl")
	files, _ := filepath.Glob(pattern)
	tmpl := template.New("components").Funcs(tmplfn.FuncMap)

	if len(files) == 0 {
		return tmpl, nil
	}

	return tmpl.ParseFiles(files...)
}
