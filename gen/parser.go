package gen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func parse(
	src []byte,
	relDir string,
	cfg Config,
) (ast.Node, []string, error) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)
	doc := p.Parse(src)
	refs, err := walkMarkdownAST(
		doc,
		relDir,
		cfg.InteralHtmlRefSuffix,
		cfg.LazyImageLoading,
	)
	if err != nil {
		return nil, nil, err
	}

	return doc, refs, nil
}

func walkMarkdownAST(
	doc ast.Node,
	relPath string,
	interalHtmlRefSuffix string,
	lazyImgLoading bool,
) ([]string, error) {
	var refs []string
	var walkErr error

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if false {
			name := reflect.TypeOf(node).String()
			if strings.Contains(name, "ListItem") ||
				strings.Contains(name, "Text") ||
				strings.Contains(name, "Paragraph") ||
				strings.Contains(name, "List") ||
				strings.Contains(name, "HTMLSpan") ||
				strings.Contains(name, "Heading") {
			} else {
				fmt.Printf("node type = %v, entering %v\n",
					reflect.TypeOf(node).String(),
					entering,
				)
			}
		}

		if !entering {
			return ast.GoToNext
		}

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
			htmlBody, htmlRefs, err := walkHtmlDOMs(relPath, v.Literal, lazyImgLoading)
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

func walkHtmlDOMs(
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

		html = []byte(stripHTMLWrapper(htmlStr))
	}

	return html, refs, nil
}

func stripHTMLWrapper(html string) string {
	html = strings.Replace(html, "<html><head></head><body>", "", 1)
	return strings.Replace(html, "</body></html>", "", 1)
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
