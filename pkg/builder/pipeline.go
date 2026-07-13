package builder

import (
	"io"

	"github.com/shelepuginivan/makewww/pkg/resource"
)

type Pipeline struct {
	global *GlobalContext
}

func NewPipeline(global *GlobalContext) *Pipeline {
	return &Pipeline{
		global: global,
	}
}

func (p *Pipeline) Process(res resource.Resource, w io.Writer) error {
	ok, err := p.tryCopyingAsIs(res, w)
	if ok {
		return err
	}

	if res.IsTemplate() {
		// render
	}

	// markdown => convert to html

	// write to w

	return nil
}

func (p *Pipeline) tryCopyingAsIs(res resource.Resource, w io.Writer) (bool, error) {
	if res.IsTemplate() {
		return false, nil
	}

	switch res.(type) {
	case *resource.MarkdownDocument:
		return false, nil
	}

	writerTo, ok := res.(io.WriterTo)
	if !ok {
		return false, nil
	}

	_, err := writerTo.WriteTo(w)
	return true, err
}
