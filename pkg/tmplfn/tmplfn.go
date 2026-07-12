// Package tmplfn provides template functions.
package tmplfn

import (
	"fmt"
	"html/template"
)

var FuncMap = template.FuncMap{
	"props": props,
}

func props(keyval ...any) (map[string]any, error) {
	if len(keyval)%2 != 0 {
		return nil, fmt.Errorf("expected even number of arguments")
	}

	res := make(map[string]any)

	for i := 0; i < len(keyval); i += 2 {
		key := keyval[i]
		val := keyval[i+1]

		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("keys must be strings")
		}

		res[keyStr] = val
	}

	return res, nil
}
