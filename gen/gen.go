package gen

import (
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
	cfg, content, err := extractPageConfig(src)
	if err != nil {
		return nil, err
	}

	doc, refs, err := parsePage(content, relDir, h.cfg)
	if err != nil {
		return nil, err
	}

	body, err := renderPage(doc, h.cfg)
	if err != nil {
		return nil, err
	}

	if cfg.bannerRef() != "" {
		refs = append(refs, cfg.bannerRef())
	}

	return &Doc{
		Html: fillPageTemplate(cfg.header(), body, cfg.footer()),
		Refs: refs,
	}, nil
}
