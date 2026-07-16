package tmplfn

import (
	"fmt"
	"iter"
)

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

func seq(args ...int) (iter.Seq[int], error) {
	var start, end, step int

	switch len(args) {
	case 1:
		start = 0
		end = args[0]
		step = 1
	case 2:
		start = args[0]
		end = args[1]
		step = 1
	case 3:
		start = args[0]
		end = args[1]
		step = args[2]
	default:
		return nil, fmt.Errorf("expected 1, 2, or 3 arguments")
	}

	return func(yield func(int) bool) {
		for i := start; i < end; i += step {
			if !yield(i) {
				break
			}
		}
	}, nil
}
