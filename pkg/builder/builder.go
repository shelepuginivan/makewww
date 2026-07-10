package builder

import (
	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/dist"
	"github.com/shelepuginivan/makewww/pkg/source"
)

type Builder struct {
	src  *source.Source
	dist *dist.Dist
}

func New(cfg *config.Config) *Builder {
	src := source.FromProjectRoot(cfg.Dir)
	dist := dist.FromRoot(cfg.Output)

	return &Builder{src: src, dist: dist}
}

func (b *Builder) Build() error {
	docs, err := b.src.GetDocuments()
	if err != nil {
		return err
	}

	for _, doc := range docs {
		p, err := doc.CanonicalPath(b.src.ContentDir())
		if err != nil {
			return err
		}

		f, err := b.dist.CreateOutputFile(p)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := doc.Render(f); err != nil {
			return err
		}
	}

	return nil
}
