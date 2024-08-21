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

type Html struct {
	interalHtmlRefSuffix string
}

func NewHtml(interalHtmlRefSuffix string) *Html {
	return &Html{
		interalHtmlRefSuffix: interalHtmlRefSuffix,
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
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(src)
	refs, err := walkRefs(doc, relPath, h.interalHtmlRefSuffix)
	if err != nil {
		return nil, err
	}

	return &Doc{
		Html: markdown.Render(
			doc,
			html.NewRenderer(
				html.RendererOptions{
					Flags: html.CommonFlags | html.HrefTargetBlank,
				},
			),
		),
		Refs: refs,
	}, nil
}

func walkRefs(
	doc ast.Node,
	relPath string,
	interalHtmlRefSuffix string,
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
			htmlRefs, err := walkHtmlRefs(relPath, v.Literal)
			if err != nil {
				walkErr = err
				return ast.Terminate
			}

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
) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var refs []string

	doc.Find("img").Each(func(_ int, sel *goquery.Selection) {
		if ref, exists := sel.Attr("src"); exists {
			if isExternalLink(ref) {
				return
			}

			if !strings.HasPrefix(ref, relPath) {
				ref = filepath.Join(relPath, ref)
			}

			refs = append(refs, ref)
		}
	})

	return refs, nil
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
