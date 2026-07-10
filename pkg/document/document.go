package document

import "io"

type Document interface {
	Render(w io.Writer) error
}
