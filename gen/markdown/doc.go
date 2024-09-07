package markdown

import "html/template"

// Link or reference used in the markdown can be relative the current file
// location, it is ok as browser appends the relative path and the server
// always receives the full path relative to the server root.
type Doc struct {
	Html         template.HTML
	InternalRefs []string
	Headings     []*Heading
}

type Heading struct {
	Level    int
	ID       string
	Name     string
	Children []*Heading
}
