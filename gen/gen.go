package gen

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Config struct {
	InteralHtmlRefSuffix string
	LazyImageLoading     bool
}

func DefaultConfig(
	internalHtmlRefSuffix string,
) Config {
	return Config{
		InteralHtmlRefSuffix: internalHtmlRefSuffix,
		LazyImageLoading:     true,
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
	relPath string,
	src []byte,
) (*Doc, error) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)
	doc := p.Parse(src)
	refs, err := walkRefs(
		doc,
		relPath,
		h.cfg.InteralHtmlRefSuffix,
		h.cfg.LazyImageLoading,
	)
	if err != nil {
		return nil, err
	}

	flags := html.CommonFlags | html.HrefTargetBlank
	if h.cfg.LazyImageLoading {
		flags |= html.LazyLoadImages
	}
	body := markdown.Render(
		doc,
		html.NewRenderer(
			html.RendererOptions{
				Flags: flags,
			},
		),
	)

	return &Doc{
		Html: bytes.Replace(defaultLayout, placeHolder, body, 1),
		Refs: refs,
	}, nil
}

func walkRefs(
	doc ast.Node,
	relPath string,
	interalHtmlRefSuffix string,
	lazyImgLoading bool,
) ([]string, error) {
	var refs []string
	var walkErr error

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}

		/*
			name := reflect.TypeOf(node).String()
			if strings.Contains(name, "ListItem") ||
				strings.Contains(name, "Text") ||
				strings.Contains(name, "Paragraph") ||
				strings.Contains(name, "List") {
			} else {
				fmt.Printf("node type = %v\n", reflect.TypeOf(node).String())
			}
		*/

		switch v := node.(type) {
		case *ast.Link:
			ref := string(v.Destination)
			if isExternalLink(ref) {
				break
			}

			if !strings.HasPrefix(ref, relPath) {
				ref = filepath.Join(relPath, ref)
			}

			v.Destination = []byte(ref + interalHtmlRefSuffix)
			refs = append(refs, ref)

		case *ast.Image:
			ref := string(v.Destination)
			if isExternalLink(ref) {
				break
			}

			if !strings.HasPrefix(ref, relPath) {
				ref = filepath.Join(relPath, ref)
			}

			refs = append(refs, ref)

		case *ast.HTMLSpan:
			htmlBody, htmlRefs, err := walkHtmlRefs(relPath, v.Literal, lazyImgLoading)
			if err != nil {
				walkErr = err
				return ast.Terminate
			}

			v.Literal = htmlBody
			refs = append(refs, htmlRefs...)
		}

		return ast.GoToNext
	})

	if walkErr != nil {
		return nil, walkErr
	}

	return refs, nil
}

func walkHtmlRefs(
	relPath string,
	data []byte,
	lazyImgLoading bool,
) ([]byte, []string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	var refs []string
	hasUpdate := false
	doc.Find("img").Each(func(_ int, sel *goquery.Selection) {
		if ref, exists := sel.Attr("src"); exists {
			if isExternalLink(ref) {
				return
			}

			if lazyImgLoading {
				sel.SetAttr("loading", "lazy")
				hasUpdate = true
			}

			if !strings.HasPrefix(ref, relPath) {
				ref = filepath.Join(relPath, ref)
			}

			refs = append(refs, ref)
		}
	})

	html := data
	if hasUpdate {
		// NOTE(kmax): be careful, <span> and </span> are treated as
		// different ast.Node. So we only parse <span>, but doc.Html
		// would complete the closing </span>. If we don't handle it
		// properly, it may end up with 2 </span>. Luckily, for now
		// <img> does not have the closing </img> pair.
		htmlStr, err := doc.Html()
		if err != nil {
			return nil, nil, err
		}

		htmlStr = strings.Replace(htmlStr, "<html><head></head><body>", "", 1)
		htmlStr = strings.Replace(htmlStr, "</body></html>", "", 1)
		html = []byte(htmlStr)
	}

	return html, refs, nil
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
