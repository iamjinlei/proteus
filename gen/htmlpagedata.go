package gen

import (
	"html/template"
)

var (
	defaultPalette = Palette{
		LightGray: cLightGray,
	}
)

type Palette struct {
	LightGray color
}

type Dimensions struct {
	CenterColWidth string
}

type Content struct {
	Header template.HTML
	Navi   template.HTML
	Main   template.HTML
	Footer template.HTML
}

type HtmlPageData struct {
	Palette    Palette
	Dimensions Dimensions
	Content    Content
}

func newHtmlPageData(
	header template.HTML,
	navi template.HTML,
	main template.HTML,
	footer template.HTML,
) *HtmlPageData {
	return &HtmlPageData{
		Palette: defaultPalette,
		Dimensions: Dimensions{
			CenterColWidth: centerColWidth,
		},
		Content: Content{
			Header: header,
			Navi:   navi,
			Main:   main,
			Footer: footer,
		},
	}
}
