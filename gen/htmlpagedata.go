package gen

import (
	"html/template"
)

var (
	defaultPalette = Palette{
		LightGray: cLightGray,
		Red:       cRed,
	}
)

type Palette struct {
	LightGray color
	Red       color
}

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
	Palette    Palette
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
		Palette: defaultPalette,
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
