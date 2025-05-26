package markdown

import (
	"fmt"

	"github.com/iamjinlei/proteus/gen/color"
)

var (
	defaultDocCss = fmt.Sprintf(`
		table, th, td {
			border:1px solid %v;
  			border-collapse: collapse;
			margin:0.2em;
			padding:0.2em 1em;
		}
		th {
			background-color:%v;
		}
		.codecls {
			padding-left:0.3em;
			padding-right:0.3em;
			background-color:%v;
	    }

		.codeblockcls {
			padding:0.1em 1.5em;
			background-color:%v;
		}
		`,
		color.LightGray,
		color.LightGray,
		color.LightGray,
		color.LightGray,
	)
)
