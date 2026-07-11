package source

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templatesDir = "templates"
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

func (src *Source) ContentDir() string {
	return filepath.Join(src.root.Name(), "content")
}

func (src *Source) TemplatesDir() string {
	return filepath.Join(src.root.Name(), "templates")
}

func (src *Source) Documents() ([]Document, error) {
	var docs []Document

	err := filepath.Walk(src.ContentDir(), func(sourceFile string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		path, err := filepath.Rel(src.ContentDir(), sourceFile)
		if err != nil {
			return err
		}

		var doc Document

		switch {
		case strings.HasSuffix(sourceFile, ".md.tmpl"):
			doc, err = markdownFromPath(path, sourceFile)
		case strings.HasSuffix(sourceFile, ".tmpl"):
			doc, err = templateFromPath(path, sourceFile)
		default:
			return nil
		}

		if err != nil {
			return err
		}

		docs = append(docs, doc)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (src *Source) RawFiles() ([]*Raw, error) {
	var files []*Raw

	err := filepath.Walk(src.ContentDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(filepath.Rel(src.ContentDir(), path))

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		raw := rawFromPath(path)
		files = append(files, raw)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (src *Source) GetTemplate(path string) (*template.Template, error) {
	tmplPath := filepath.Join(templatesDir, path)
	content, err := src.root.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("template").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}
