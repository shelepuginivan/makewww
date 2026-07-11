package source

import (
	"fmt"
	"io"
	"os"
)

type Raw struct {
	path       string
	sourceFile string
}

func rawFromPath(path, sourceFile string) *Raw {
	return &Raw{
		path:       path,
		sourceFile: sourceFile,
	}
}

func (r *Raw) CopyTo(w io.Writer) error {
	file, err := os.Open(r.sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open raw document: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("failed to write raw document: %w", err)
	}

	return nil
}

func (r *Raw) Path() *Path {
	return &Path{r.path}
}
