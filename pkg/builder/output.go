package builder

import (
	"fmt"
	"os"
	"path/filepath"
)

type Output struct {
	root *os.Root
}

func NewOutput(root string) (*Output, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	r, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}

	return &Output{root: r}, nil
}

func (out *Output) CreateOutputFile(path string) (*os.File, error) {
	if err := out.root.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create dir: %w", err)
	}

	f, err := out.root.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return f, nil
}
