package gen

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/iamjinlei/proteus/gen/color"
	"github.com/iamjinlei/proteus/gen/markdown"
)

type Config struct {
	Domain                string
	InternalRefHtmlSuffix string
	LazyImageLoading      bool
	Styles                markdown.Styles
}

func DefaultConfig(
	domain string,
	internalRefHtmlSuffix string,
) Config {
	return Config{
		Domain:                domain,
		InternalRefHtmlSuffix: internalRefHtmlSuffix,
		LazyImageLoading:      true,
		Styles: markdown.Styles{
			CodeBlock: fmt.Sprintf(
				"padding:0.1em 1.5em;background-color:%v;",
				color.LightGray,
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
		mdp: markdown.NewParser(),
		mdr: markdown.NewRenderer(
			color.DefaultPalette,
			cfg.Styles,
			cfg.InternalRefHtmlSuffix,
			cfg.LazyImageLoading,
		),
		r: r,
	}, nil
}

type Page struct {
	Html         []byte
	InternalRefs []string
}

func (h *Html) Gen(relPath string, src []byte) (*Page, error) {
	pCfg, md, err := extractPageConfig(src)
	if err != nil {
		return nil, err
	}

	mdDoc, err := h.mdr.Render(h.mdp.Parse(md))
	if err != nil {
		return nil, err
	}

	var toc *HtmlComponent
	if pCfg.leftPane() == "toc" {
		toc = renderToC(mdDoc.Headings, 3)
	} else {
		toc = &HtmlComponent{}
	}

	refs := mdDoc.InternalRefs
	if pCfg.bannerRef() != "" {
		refs = append(refs, pCfg.bannerRef())
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := h.r.render(
		w,
		newTemplateData(
			h.cfg.Domain,
			relPath,
			pCfg.header(),
			pCfg.nav(),
			&HtmlComponent{
				Html: mdDoc.Html,
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
		Html:         b.Bytes(),
		InternalRefs: refs,
	}, nil
}
