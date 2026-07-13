package builder

import (
	"bytes"
	"io"
	"text/template"

	"github.com/shelepuginivan/makewww/pkg/resource"
)

type Pipeline struct {
	global     *GlobalContext
	components *template.Template
	layouts    map[string]*template.Template
}

func NewPipeline(
	global *GlobalContext,
	components *template.Template,
	layouts map[string]*template.Template,
) *Pipeline {
	return &Pipeline{
		global:     global,
		components: components,
		layouts:    layouts,
	}
}

func (p *Pipeline) Process(res resource.Resource, w io.Writer) error {
	ok, err := p.tryCopyingAsIs(res, w)
	if ok {
		return err
	}

	var content []byte
	if res.IsTemplate() {
		content, err = p.renderTemplate(res)
	} else {
		content, err = res.Content()
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

func (p *Pipeline) renderTemplate(res resource.Resource) ([]byte, error) {
	content, err := res.Content()
	if err != nil {
		return nil, err
	}

	base, err := p.components.Clone()
	if err != nil {
		return nil, err
	}

	tmpl, err := base.Parse(string(content))
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"Global": p.global,
		"Path":   res.Path(),
	}
	if withMetadata, ok := res.(resource.WithMetadata); ok {
		data["Metadata"] = withMetadata.Metadata()
	}

	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
