package gen

import (
	"strings"

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

func (h *Html) Gen(src []byte) []byte {
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := updateRefUrls(p.Parse(src), h.interalHtmlRefSuffix)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func updateRefUrls(doc ast.Node, interalHtmlRefSuffix string) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if link, ok := node.(*ast.Link); ok && entering {
			ref := string(link.Destination)
			if !isExternalLink(ref) {
				link.Destination = []byte(ref + interalHtmlRefSuffix)
			}
		}

		return ast.GoToNext
	})

	return doc
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
