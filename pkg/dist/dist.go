package dist

import (
	"fmt"
	"os"
	"path/filepath"
)

type Dist struct {
	root string
}

func FromRoot(root string) *Dist {
	return &Dist{root: root}
}

func (d *Dist) CreateOutputFile(path string) (*os.File, error) {
	filePath := filepath.Join(d.root, path)

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return f, nil
}
