package tmplfn

import (
	"fmt"
	"strings"
)

func replace(oldnew []string, s string) (string, error) {
	if len(oldnew)%2 != 0 {
		return "", fmt.Errorf("expected even number of arguments")
	}
	return strings.NewReplacer(oldnew...).Replace(s), nil
}
