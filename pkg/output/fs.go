package output

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FS struct {
	root *os.Root
}

func NewFS(root string) (*FS, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	oldEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range oldEntries {
		entryPath := filepath.Join(root, entry.Name())

		if err := os.RemoveAll(entryPath); err != nil {
			return nil, err
		}
	}

	r, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}

	return &FS{root: r}, nil
}

func (fs *FS) WriterForPath(path string) (io.WriteCloser, error) {
	if err := fs.root.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create dir: %w", err)
	}

	f, err := fs.root.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return f, nil
}
