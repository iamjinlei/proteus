package gen

import (
	"html/template"
	"io"
)

type renderer struct {
	tpl *template.Template
}

func newRenderer(layout string) (*renderer, error) {
	tpl, err := template.New("test").Parse(layout)
	if err != nil {
		return nil, err
	}

	return &renderer{
		tpl: tpl,
	}, nil
}

func (r *renderer) render(w io.Writer, d *HtmlPageData) error {
	if err := r.tpl.Execute(w, d); err != nil {
		return err
	}

	return nil
}
