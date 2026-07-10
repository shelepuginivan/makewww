package document

import (
	"fmt"
	"io"
	"os"
)

type Raw struct {
	path string
}

func RawFromPath(path string) *Raw {
	return &Raw{path: path}
}

func (r *Raw) Render(w io.Writer) error {
	file, err := os.Open(r.path)
	if err != nil {
		return fmt.Errorf("failed to open raw document: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("failed to write raw document: %w", err)
	}

	return nil
}
