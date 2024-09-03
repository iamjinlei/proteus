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

type Styles struct {
	MainLeft template.CSS
}

type Content struct {
	Header   template.HTML
	Navi     template.HTML
	Main     template.HTML
	MainLeft template.HTML
	Footer   template.HTML
}

type HtmlPageData struct {
	Palette    Palette
	Dimensions Dimensions
	Styles     Styles
	Content    Content
}

func newHtmlPageData(
	header template.HTML,
	navi template.HTML,
	main template.HTML,
	mainLeft template.HTML,
	mainLeftStyle template.CSS,
	footer template.HTML,
) *HtmlPageData {
	return &HtmlPageData{
		Palette: defaultPalette,
		Dimensions: Dimensions{
			CenterColWidth: centerColWidth,
		},
		Styles: Styles{
			MainLeft: mainLeftStyle,
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
