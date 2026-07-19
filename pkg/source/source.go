// Package source provides methods for working with website source code.
package source

import (
	"text/template"

	"github.com/shelepuginivan/makewww/pkg/resource"
)

type Source interface {
	Resources() ([]resource.Resource, error)
	Layouts() (map[string]*template.Template, error)
	Components() (*template.Template, error)
}
