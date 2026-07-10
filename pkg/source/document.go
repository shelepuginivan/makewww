package source

import "io"

type Document interface {
	CanonicalPath(base string) (string, error)
	Render(w io.Writer) error
}
