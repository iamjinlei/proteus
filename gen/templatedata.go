package gen

import (
	"html/template"

	"github.com/iamjinlei/proteus/gen/color"
)

type Dimensions struct {
}

type HtmlComponent struct {
	Html template.HTML
	Css  template.CSS
	Js   template.JS
}

type Content struct {
	Header   *HtmlComponent
	Nav      *HtmlComponent
	Main     *HtmlComponent
	MainLeft *HtmlComponent
	Footer   *HtmlComponent
}

type TemplateData struct {
	CanonicalDomain string
	RelPath         string
	Palette         color.Palette
	Dimensions      Dimensions
	Content         Content
}

func newTemplateData(
	domain string,
	relPath string,
	header *HtmlComponent,
	nav *HtmlComponent,
	main *HtmlComponent,
	mainLeft *HtmlComponent,
	footer *HtmlComponent,
) *TemplateData {
	return &TemplateData{
		CanonicalDomain: normalizeDomain(domain),
		RelPath:         normalizeRelPath(relPath),
		Palette:         color.DefaultPalette,
		Dimensions:      Dimensions{},
		Content: Content{
			Header:   header,
			Nav:      nav,
			Main:     main,
			MainLeft: mainLeft,
			Footer:   footer,
		},
	}
}
