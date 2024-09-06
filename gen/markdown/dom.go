package markdown

import (
	"fmt"
	"io"

	"github.com/iamjinlei/proteus/gen/color"
)

func bookBibliography(
	w io.Writer,
	palette color.Palette,
	title string,
	coverImgRef string,
	link string,
	author string,
) {
	fmt.Fprintf(
		w,
		`
<div style="width:100%%;margin-bottom:100px;border-bottom: 2px solid %s;display:grid;grid-template-columns: 1fr 2fr;">
	<span><a href="%s"><img src="%s" style="width:100%%"></a></span>
	<span style="padding-left:40px;">
		<div style="font-size:2em;font-weight: bold;">%s</div>
		<div style="font-size:1.2em;margin-top:5px;">作者: %s</div>
	</span>
</div>`,
		palette.DarkGray,
		link,
		coverImgRef,
		title,
		author,
	)
}

func highlight(w io.Writer, content []byte, color string) {
	fmt.Fprintf(
		w,
		`<span style="background-color:%s;">%s</span>`,
		color,
		string(content),
	)
}
