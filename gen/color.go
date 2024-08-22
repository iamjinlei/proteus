package gen

import "strings"

type color string

const (
	cLightGray color = "#F0F0F0"
)

var (
	colorMap = map[string]color{
		"#LightGray": cLightGray,
	}
)

func fillColors(str string) string {
	for name, color := range colorMap {
		str = strings.Replace(str, name, string(color), -1)
	}
	return str
}
