package gen

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/iamjinlei/proteus/gen/markdown"
)

type Config struct {
	InteralHtmlRefSuffix string
	LazyImageLoading     bool
	Styles               markdown.Styles
}

func DefaultConfig(
	internalHtmlRefSuffix string,
) Config {
	return Config{
		InteralHtmlRefSuffix: internalHtmlRefSuffix,
		LazyImageLoading:     true,
		Styles: markdown.Styles{
			CodeBlock: fmt.Sprintf(
				"padding:0.1em 1.5em;background-color:%v;",
				cLightGray,
			),
		},
	}
}

type Html struct {
	cfg Config
	mdp *markdown.Parser
	mdr *markdown.Renderer
	r   *renderer
}

func NewHtml(cfg Config) (*Html, error) {
	r, err := newRenderer(defaultLayout)
	if err != nil {
		return nil, err
	}

	return &Html{
		cfg: cfg,
		mdp: markdown.NewParser(
			cfg.InteralHtmlRefSuffix,
			cfg.LazyImageLoading,
		),
		mdr: markdown.NewRenderer(
			cfg.LazyImageLoading,
		),
		r: r,
	}, nil
}

type Page struct {
	Html []byte
	Refs []string
}

func (h *Html) Gen(
	src []byte,
	relDir string,
) (*Page, error) {
	pCfg, md, err := extractPageConfig(src)
	if err != nil {
		return nil, err
	}

	mdDoc, err := h.mdp.Parse(md, relDir)
	if err != nil {
		return nil, err
	}

	mdHtml, err := h.mdr.Render(
		mdDoc.Root,
		h.cfg.Styles,
	)
	if err != nil {
		return nil, err
	}

	var toc *HtmlComponent
	if pCfg.leftPane() == "toc" {
		toc = renderToC(mdDoc.Headings, 3)
	} else {
		toc = &HtmlComponent{}
	}

	refs := mdDoc.Refs
	if pCfg.bannerRef() != "" {
		refs = append(refs, pCfg.bannerRef())
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := h.r.render(
		w,
		newHtmlPageData(
			pCfg.header(),
			pCfg.navi(),
			&HtmlComponent{
				Html: mdHtml,
			},
			toc,
			pCfg.footer(),
		),
	); err != nil {
		return nil, err
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	return &Page{
		Html: b.Bytes(),
		Refs: refs,
	}, nil
}
