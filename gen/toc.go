package gen

import (
	"fmt"
	"html/template"

	"github.com/iamjinlei/proteus/gen/markdown"
)

const (
	defaultToCCss = template.CSS(`
	.toc {
	position: -webkit-sticky; /* Safari */
	position: sticky;
	float: right;
	top: 5em;
	margin-right: 2em;
	}
	`)
)

func renderToC(
	hs []*markdown.Heading,
	maxDepth int,
) (template.HTML, template.CSS) {
	if len(hs) == 0 {
		return template.HTML(""), template.CSS("")
	}

	toc := renderHeadingList(hs, 0, maxDepth)
	return template.HTML(`<div class="toc">` + toc + `</div>`), defaultToCCss
}

func renderHeadingList(
	hs []*markdown.Heading,
	depth int,
	maxDepth int,
) string {
	html := fmt.Sprintf(`<ul class="toc_ul_%d">`, depth)
	for _, h := range hs {
		html += fmt.Sprintf(
			`<li class="toc_li_%d"><a href="#%s">%s</a></li>`,
			depth,
			h.ID,
			h.Name,
		)
		if len(h.Children) > 0 && depth+1 < maxDepth {
			html += renderHeadingList(h.Children, depth+1, maxDepth)
		}
	}
	html += "</ul>"

	return html
}
