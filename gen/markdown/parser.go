package markdown

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

type Doc struct {
	Root     ast.Node
	Refs     []string
	Headings []*Heading
}

type Heading struct {
	Level    int
	Name     string
	Children []*Heading
}

func Parse(
	src []byte,
	relDir string,
	interalHtmlRefSuffix string,
	lazyImageLoading bool,
) (*Doc, error) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)
	root := p.Parse(src)

	c, err := buildMarkdownContent(
		root,
		relDir,
		interalHtmlRefSuffix,
		lazyImageLoading,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func buildMarkdownContent(
	root ast.Node,
	relPath string,
	interalHtmlRefSuffix string,
	lazyImgLoading bool,
) (*Doc, error) {
	var walkErr error
	// Accumulate references found in the doc.
	var refs []string
	// Track headings to build table of contents.
	ht := newHeadingTracker()

	ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
		if false {
			name := reflect.TypeOf(node).String()
			if strings.Contains(name, "Heading") {
				children := node.(*ast.Heading).Children
				c0 := children[0].(*ast.Text)
				fmt.Printf("cc = %v, node = %#v\n", len(children), string(c0.Literal))
			} else if strings.Contains(name, "Text") ||
				strings.Contains(name, "ListItem") ||
				strings.Contains(name, "Paragraph") ||
				strings.Contains(name, "List") ||
				strings.Contains(name, "HTMLSpan") ||
				strings.Contains(name, "Heading") {
			} else {
				fmt.Printf("node type = %v, entering %v\n",
					name,
					entering,
				)
			}
		}

		if !entering {
			return ast.GoToNext
		}

		switch v := node.(type) {
		case *ast.Heading:
			if len(v.Children) != 1 {
				break
			}

			ht.add(v.Level, string(v.Children[0].(*ast.Text).Literal))

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

	return &Doc{
		Root:     root,
		Refs:     refs,
		Headings: ht.getHeadings(),
	}, nil
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
