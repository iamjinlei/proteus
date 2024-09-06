package markdown

import (
	"fmt"

	"github.com/iamjinlei/proteus/gen/color"
)

func bookBibliography(
	palette color.Palette,
	title string,
	coverImgRef string,
	link string,
	author string,
) ([]byte, []byte) {
	return []byte(fmt.Sprintf(`
<div style="width:100%%;margin-bottom:100px;border-bottom: 2px solid %s;display:grid;grid-template-columns: 1fr 2fr;">
	<span><a href="%s"><img src="%s" style="width:100%%"></a></span>
	<span style="padding-left:40px;">
		<div style="font-size:3em;font-weight: bold;">%s</div>
		<div style="font-size:1.4em;margin-top:5px;">作者: %s</div>
	</span>
`,
		palette.DarkGray,
		link,
		coverImgRef,
		title,
		author,
	)), htmlClosingTagDiv
}

func highlight(color string) ([]byte, []byte) {
	return []byte(fmt.Sprintf(
		`<span style="background-color:%s;">`,
		color,
	)), htmlClosingTagSpan
}
