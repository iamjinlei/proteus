package markdown

import (
	"github.com/iamjinlei/proteus/gen/color"
	"github.com/iamjinlei/proteus/gen/keyword"
)

type Keyword struct {
	Value  string
	Target string
}

type Keywords struct {
	colorMap map[string]color.Color
	index    map[keyword.Type][]*Keyword
}

func newKeywords(
	colorMap map[string]color.Color,
) *Keywords {
	return &Keywords{
		colorMap: colorMap,
		index:    map[keyword.Type][]*Keyword{},
	}
}

func (k *Keywords) add(
	t keyword.Type,
	val string,
	target string,
) {
	k.index[t] = append(k.index[t], &Keyword{
		Value:  val,
		Target: target,
	})
}

func (k *Keywords) Get(t keyword.Type) []*Keyword {
	return k.index[t]
}

func (k *Keywords) Color(t keyword.Type) color.Color {
	return k.colorMap[string(t)]
}
