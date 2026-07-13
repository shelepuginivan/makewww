package builder

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/shelepuginivan/makewww/pkg/resource"
	"github.com/yuin/goldmark"
)

type Pipeline struct {
	global     *GlobalContext
	components *template.Template
	layouts    map[string]*template.Template
	md         goldmark.Markdown
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
	if err != nil {
		return err
	}

	switch r := res.(type) {
	case *resource.Markdown:
		err = p.convertMarkdown(r, content, w)
	default:
		_, err = w.Write(content)
	}

	return err
}

func (p *Pipeline) tryCopyingAsIs(res resource.Resource, w io.Writer) (bool, error) {
	if res.IsTemplate() {
		return false, nil
	}

	switch res.(type) {
	case *resource.Markdown:
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

func (p *Pipeline) convertMarkdown(res *resource.Markdown, content []byte, w io.Writer) error {
	buffer := new(bytes.Buffer)

	err := p.md.Convert(content, buffer)
	if err != nil {
		return err
	}

	metadata := res.Metadata()

	layout, exists := p.layouts[metadata.Template]
	if !exists {
		return fmt.Errorf("layout %s does not exist", metadata.Template)
	}

	data := map[string]any{
		"Global":  p.global,
		"Content": string(content),
		"Document": map[string]any{
			"Metadata": metadata,
			"Path":     res.Path(),
		},
	}

	return layout.Execute(w, data)
}
