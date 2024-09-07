package gen

import (
	"bufio"
	"bytes"

	"github.com/iamjinlei/proteus/gen/color"
	"github.com/iamjinlei/proteus/gen/markdown"
)

type Config struct {
	Domain                string
	InternalRefHtmlSuffix string
	LazyImageLoading      bool
	Palette               color.Palette
}

func DefaultConfig(
	domain string,
	internalRefHtmlSuffix string,
) Config {
	return Config{
		Domain:                domain,
		InternalRefHtmlSuffix: internalRefHtmlSuffix,
		LazyImageLoading:      true,
		Palette:               color.DefaultPalette,
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
			cfg.Palette,
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
			renderComponent(pCfg.leftPane(), mdDoc, h.cfg.Palette),
			renderComponent(pCfg.rightPane(), mdDoc, h.cfg.Palette),
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

func renderComponent(
	kind string,
	doc *markdown.Doc,
	palette color.Palette,
) *HtmlComponent {
	switch kind {
	case "toc":
		return renderToC(doc.Headings, 3)
	case "kws":
		return renderKeywords(doc.Keywords, palette)
	}

	return &HtmlComponent{}
}
