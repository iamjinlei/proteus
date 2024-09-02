package gen

import "github.com/gomarkdown/markdown/ast"

type page struct {
	root     ast.Node
	refs     []string
	headings []*heading
}

type heading struct {
	level    int
	name     string
	children []*heading
}
