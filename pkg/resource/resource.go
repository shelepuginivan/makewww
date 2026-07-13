// Package resource provides types for website resources.
package resource

import "strings"

type Resource interface {
	Path() *Path
	Content() ([]byte, error)
	IsTemplate() bool
}

func New(path, sourceFile string) (Resource, error) {
	var (
		res Resource
		err error
	)

	switch {
	case strings.HasSuffix(sourceFile, ".html.tmpl"):
		res, err = NewHTML(path, sourceFile, true)
	case strings.HasSuffix(sourceFile, ".html"):
		res, err = NewHTML(path, sourceFile, false)
	case strings.HasSuffix(sourceFile, ".md.tmpl"):
		res, err = NewMarkdown(path, sourceFile, true)
	case strings.HasSuffix(sourceFile, ".md"):
		res, err = NewMarkdown(path, sourceFile, false)
	case strings.HasSuffix(sourceFile, ".tmpl"):
		res = NewRaw(path, sourceFile, true)
	default:
		res = NewRaw(path, sourceFile, false)
	}

	return res, err
}
