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
	Navi     *HtmlComponent
	Main     *HtmlComponent
	MainLeft *HtmlComponent
	Footer   *HtmlComponent
}

type TemplateData struct {
	Palette    color.Palette
	Dimensions Dimensions
	Content    Content
}

func newTemplateData(
	header *HtmlComponent,
	navi *HtmlComponent,
	main *HtmlComponent,
	mainLeft *HtmlComponent,
	footer *HtmlComponent,
) *TemplateData {
	return &TemplateData{
		Palette:    color.DefaultPalette,
		Dimensions: Dimensions{},
		Content: Content{
			Header:   header,
			Navi:     navi,
			Main:     main,
			MainLeft: mainLeft,
			Footer:   footer,
		},
	}
}
