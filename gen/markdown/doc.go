package markdown

import "html/template"

// Link or reference used in the markdown can be relative the current file
// location, it is ok as browser appends the relative path and the server
// always receives the full path relative to the server root.
type Doc struct {
	Html         template.HTML
	Css          template.CSS
	InternalRefs []string
	Headings     []*Heading
	Keywords     *Keywords
}

type Heading struct {
	Level    int
	ID       string // Original heading ID
	ParentID string // Original parent heading ID
	Name     string
	Children []*Heading
}

func (h *Heading) GetHtmlDomID() string {
	if h.ParentID == "" {
		return h.ID
	}
	return h.ParentID + "_" + h.ID
}
