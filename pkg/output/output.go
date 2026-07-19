// Package output provides primitives to output the built website.
package output

import "io"

type Output interface {
	WriterForPath(string) (io.WriteCloser, error)
}
