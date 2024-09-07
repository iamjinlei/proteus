package gen

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/iamjinlei/proteus/gen/color"
	"github.com/iamjinlei/proteus/gen/keyword"
	"github.com/iamjinlei/proteus/gen/markdown"
)

type kwType string

const (
	kwName kwType = "name"
)

const (
	defaultKwsCss = `
.kws {
	position: -webkit-sticky; /* Safari */
	position: sticky;
	float: left;
	top: 10em;
	margin-left: 2em;
	font-size: 1em;
}
.kws .namebox {
	display: inline-block;
	background-color: {{ .NameBoxBgColor }};
	border: 1px solid {{ .Palette.LightGray }};
	border-radius:4px;
	padding: 4px 8px;
	margin: 4px;
}
.kws a {
	text-decoration: none;
	color: #000000;
}
`
)

func renderKeywords(
	kws *markdown.Keywords,
	palette color.Palette,
) *HtmlComponent {
	list := kws.Get(keyword.Name)

	seen := map[string]bool{}
	spans := ""
	for _, kw := range list {
		if seen[kw.Value] {
			continue
		}
		seen[kw.Value] = true

		spans += fmt.Sprintf(
			`<span class="namebox"><a href="#%s">%s</a></span>`,
			kw.Target,
			kw.Value,
		)
	}
	css := strings.Replace(
		defaultKwsCss,
		"{{ .Palette.LightGray }}",
		palette.LightGray.Hex(),
		-1,
	)
	css = strings.Replace(
		css,
		"{{ .NameBoxBgColor }}",
		kws.Color(keyword.Name).Hex(),
		-1,
	)

	return &HtmlComponent{
		Html: template.HTML(fmt.Sprintf(
			`<div class="kws">%s</div>`,
			spans,
		)),
		Css: template.CSS(css),
	}
}
