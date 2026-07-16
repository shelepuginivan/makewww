// Package tmplfn provides template functions.
package tmplfn

import (
	"html/template"
	"strings"
	"time"
)

var FuncMap = template.FuncMap{
	// Resource manipulation.
	"draft":             draft,
	"in_path":           inPath,
	"latest":            latest,
	"not_draft":         notDraft,
	"sort_latest_first": sortLatestFirst,
	"sort_by_order":     sortByOrder,

	// Strings.
	"contains":    pipeline(strings.Contains),
	"has_prefix":  pipeline(strings.HasPrefix),
	"has_suffix":  pipeline(strings.HasSuffix),
	"trim":        strings.TrimSpace,
	"trim_prefix": pipeline(strings.TrimPrefix),
	"trim_suffix": pipeline(strings.TrimSuffix),
	"lower":       strings.ToLower,
	"title":       strings.ToTitle,
	"upper":       strings.ToUpper,
	"index":       pipeline(strings.Index),
	"replace":     replace,

	// Misc.
	"props": props,
	"now":   time.Now,
	"seq":   seq,
}

// pipeline accepts a function with 2 arguments and returns a wrapper with
// these arguments swapped. This can convert many functions to a pipeline form
// for templates. For example
//
//	strings.Contains(summary, "golang")
//
// wrapped in pipeline, this can be used as follows:
//
//	{{ .Summary | contains "golang" }}
func pipeline[T, U, R any](fn func(T, U) R) func(U, T) R {
	return func(u U, t T) R {
		return fn(t, u)
	}
}
