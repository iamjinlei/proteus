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
	r   *renderer
}

func NewHtml(cfg Config) (*Html, error) {
	r, err := newRenderer(defaultLayout)
	if err != nil {
		return nil, err
	}

	return &Html{
		cfg: cfg,
		r:   r,
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

	mdDoc, err := markdown.Parse(
		md,
		relDir,
		h.cfg.InteralHtmlRefSuffix,
		h.cfg.LazyImageLoading,
	)
	if err != nil {
		return nil, err
	}

	mdHtml, err := markdown.RenderHtml(
		mdDoc.Root,
		h.cfg.LazyImageLoading,
		h.cfg.Styles,
	)
	if err != nil {
		return nil, err
	}

	tocHtml, tocCss := renderToC(mdDoc.Headings, 3)

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
			mdHtml,
			tocHtml,
			tocCss,
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
