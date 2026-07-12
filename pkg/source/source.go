// Package source provides methods for working with website source code.
package source

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
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

func (src *Source) Documents() ([]Document, error) {
	var docs []Document

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

		var doc Document

		switch {
		case strings.HasSuffix(sourceFile, ".html"):
			doc, err = htmlFromPath(path, sourceFile, false)
		case strings.HasSuffix(sourceFile, ".html.tmpl"):
			doc, err = htmlFromPath(path, sourceFile, true)
		case strings.HasSuffix(sourceFile, ".md"):
			doc, err = markdownFromPath(path, sourceFile, false)
		case strings.HasSuffix(sourceFile, ".md.tmpl"):
			doc, err = markdownFromPath(path, sourceFile, true)
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

	skipExt := []string{".md", ".html", ".tmpl"}
	content := filepath.Join(src.root.Name(), contentDir)

	err := filepath.Walk(content, func(sourceFile string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if slices.Contains(skipExt, filepath.Ext(sourceFile)) {
			return nil
		}

		path, err := filepath.Rel(content, sourceFile)
		if err != nil {
			return err
		}

		raw := rawFromPath(path, sourceFile)
		files = append(files, raw)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
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
	tmpl := template.New("components")

	if len(files) == 0 {
		return tmpl, nil
	}

	return tmpl.ParseFiles(files...)
}
