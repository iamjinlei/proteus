package gen

import (
	"html/template"

	"github.com/iamjinlei/proteus/gen/color"
)

type Dimensions struct {
	CenterColWidth string
}

type HtmlComponent struct {
	Html template.HTML
	Css  template.CSS
	Js   template.JS
}

type Content struct {
	Header   *HtmlComponent
	Navi     *HtmlComponent
	Main     *HtmlComponent
	MainLeft *HtmlComponent
	Footer   *HtmlComponent
}

type HtmlPageData struct {
	Palette    color.Palette
	Dimensions Dimensions
	Content    Content
}

func newHtmlPageData(
	header *HtmlComponent,
	navi *HtmlComponent,
	main *HtmlComponent,
	mainLeft *HtmlComponent,
	footer *HtmlComponent,
) *HtmlPageData {
	return &HtmlPageData{
		Palette: color.DefaultPalette,
		Dimensions: Dimensions{
			CenterColWidth: centerColWidth,
		},
		Content: Content{
			Header:   header,
			Navi:     navi,
			Main:     main,
			MainLeft: mainLeft,
			Footer:   footer,
		},
	}
}
