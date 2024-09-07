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
	top: 10em;
	margin-right: 2em;
}
.toc a {
	text-decoration: none;
	color: #000000;
}
.toc0_ul {
	list-style-type: none;
	font-size: 1.2em;
}
.toc1_ul {
	list-style-type: none;
	font-size: 0.9em;
}
.toc2_ul {
	list-style-type: none;
	font-size: 0.9em;
}
.toc_tgl {
  	background-color: transparent;
	font-family: "Menlo", "Lucida Console", "Monaco", "Consolas", monospace;
	padding-left: 0.8em;
	font-size: 0.8em;
	border: none;
	cursor: pointer;
}
`)

	defaultToJs = template.JS(`
function toc_tgl(id) {
	var btn = document.getElementById(id);
	var c = document.getElementById(id.replace("tgl", "div"));
	console.log(btn.innerHTML);
	console.log(id.replace("tgl", "div"));
	if (btn.innerHTML === "[+]") {
	   btn.innerHTML = "[-]";
	   c.style.display = "inline";
	} else {
	   btn.innerHTML = "[+]";
	   c.style.display = "none";
	}
}
`)
)

func renderToC(
	hs []*markdown.Heading,
	maxDepth int,
) *HtmlComponent {
	if len(hs) == 0 {
		return &HtmlComponent{}
	}

	return &HtmlComponent{
		Html: template.HTML(fmt.Sprintf(
			`<div class="toc">%s</div>`,
			renderHeadingList(hs, "", 0, maxDepth),
		)),
		Css: defaultToCCss,
		Js:  defaultToJs,
	}
}

func renderHeadingList(
	hs []*markdown.Heading,
	idPrefix string,
	depth int,
	maxDepth int,
) string {
	html := fmt.Sprintf(`<ul class="toc%d_ul">`, depth)
	for idx, h := range hs {
		id := fmt.Sprintf("%d", idx)
		if idPrefix != "" {
			id = idPrefix + "." + id
		}
		html += fmt.Sprintf(
			`<li class="toc%d_li"><a href="#%s">%s</a>`,
			depth,
			h.ID,
			h.Name,
		)
		if len(h.Children) > 0 && depth+1 < maxDepth {
			html += fmt.Sprintf(
				`<button id="toc%s_tgl" class="toc_tgl" onclick="toc_tgl(this.id)" style="display:inline;">[+]</button></li>
				<div id="toc%s_div" style="display:none;">%s</div>`,
				id,
				id,
				renderHeadingList(h.Children, id, depth+1, maxDepth),
			)
		} else {
			html += `</li>`
		}
	}
	html += "</ul>"

	return html
}
