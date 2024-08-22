package gen

import (
	"bytes"
	"fmt"
)

type Styles struct {
	CodeBlock string
}

type Config struct {
	InteralHtmlRefSuffix string
	LazyImageLoading     bool
	Styles               Styles
}

func DefaultConfig(
	internalHtmlRefSuffix string,
) Config {
	return Config{
		InteralHtmlRefSuffix: internalHtmlRefSuffix,
		LazyImageLoading:     true,
		Styles: Styles{
			CodeBlock: fmt.Sprintf(
				"padding:0.1em 1.5em;background-color:%v;",
				cLightGray,
			),
		},
	}
}

type Html struct {
	cfg Config
}

func NewHtml(cfg Config) *Html {
	return &Html{
		cfg: cfg,
	}
}

type Doc struct {
	Html []byte
	Refs []string
}

func (h *Html) Gen(
	src []byte,
	relDir string,
) (*Doc, error) {
	doc, refs, err := parse(src, relDir, h.cfg)
	if err != nil {
		return nil, err
	}

	body, err := render(doc, h.cfg)
	if err != nil {
		return nil, err
	}

	return &Doc{
		Html: bytes.Replace(defaultLayout, placeHolder, body, 1),
		Refs: refs,
	}, nil
}
