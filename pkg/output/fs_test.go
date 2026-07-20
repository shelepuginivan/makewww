package output_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shelepuginivan/makewww/pkg/output"
)

type fileMock struct {
	relpath string
	content string
}

func TestOutputFS(t *testing.T) {
	root := t.TempDir()

	existingFiles := []*fileMock{
		{"existing.1", "some"},
		{"nested/existing.1", "another"},
	}

	for _, f := range existingFiles {
		path := filepath.Join(root, f.relpath)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte(f.content), 0644)
	}

	out, err := output.NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(root); err != nil {
		t.Fatal("output dir must exist")
	}

	files, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatal("output.FS must cleanup output dir")
	}

	newFiles := []*fileMock{
		{"index.html", "<!doctype html><html></html>"},
		{"a/b/c/some.md", "# Hello"},
	}

	for _, f := range newFiles {
		w, err := out.WriterForPath(f.relpath)
		if err != nil {
			t.Fatal(err)
		}

		w.Write([]byte(f.content))
		w.Close()
	}

	for _, f := range newFiles {
		content, err := os.ReadFile(filepath.Join(root, f.relpath))
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != f.content {
			t.Fatalf("%s: content mismatch", f.relpath)
		}
	}
}

func TestOutputFS_PathTraversal(t *testing.T) {
	root := t.TempDir()
	out, err := output.NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	escapePaths := []string{
		"/etc/nginx/nginx.conf",
		"///etc/fstab",
		"../../../../../../../../etc/passwd",
		"\x2E\x2E/\x2E\x2E/\x2E\x2E/.ssh/id_ed25519",
		"\x2Fgrub\x2Fgrub.cfg",
	}

	for _, path := range escapePaths {
		if _, err := out.WriterForPath(path); err == nil {
			t.Fatalf("%s escapes output root", path)
		}
	}
}
