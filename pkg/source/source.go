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
	layoutsDir    = "layouts"
)

type Source struct {
	root *os.Root
}

func New(root string) (*Source, error) {
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

func (src *Source) Layouts() (map[string]*template.Template, error) {
	components, err := src.Components()
	if err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	pattern := filepath.Join(src.root.Name(), layoutsDir, "*.tmpl")
	layoutPaths, _ := filepath.Glob(pattern)
	layouts := make(map[string]*template.Template, len(layoutPaths))

	for _, path := range layoutPaths {
		base, err := components.Clone()
		if err != nil {
			return nil, fmt.Errorf("failed to clone components: %w", err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}

		layout, err := base.Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", path, err)
		}

		layoutName := filepath.Base(path)
		layouts[layoutName] = layout
	}

	return layouts, nil
}

func (src *Source) Components() (*template.Template, error) {
	pattern := filepath.Join(src.root.Name(), componentsDir, "*.tmpl")
	files, _ := filepath.Glob(pattern)
	tmpl := template.New("components").Funcs(tmplfn.FuncMap)

	if len(files) == 0 {
		return tmpl, nil
	}

	return tmpl.ParseFiles(files...)
}
