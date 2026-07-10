package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/shelepuginivan/makewww/pkg/dist"
	"github.com/shelepuginivan/makewww/pkg/source"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	src := source.FromProjectRoot(cwd)
	dist := dist.FromRoot(filepath.Join(cwd, "dist"))

	docs, err := src.GetDocuments()
	if err != nil {
		log.Fatal(err)
	}

	for _, doc := range docs {
		p, err := doc.CanonicalPath(src.ContentDir())
		if err != nil {
			log.Fatal(err)
		}

		f, err := dist.CreateOutputFile(p)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := doc.Render(f); err != nil {
			log.Fatal(err)
		}
	}
}
