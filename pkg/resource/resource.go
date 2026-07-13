// Package resource provides types for website resources.
package resource

type Resource interface {
	Path() *Path
	Content() (string, error)
	IsTemplate() bool
}
