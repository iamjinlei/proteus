package markdown

import (
	"bytes"
	"errors"
	"fmt"

	"golang.org/x/net/html"
)

var (
	ErrUnexpectedTagLoc = errors.New("unexpected tag location")
)

func parseTag(data []byte) (*html.Node, error) {
	// html.Parse always builds a html->body tree even though the input
	// is a single tag. We will need to traverse down the tree to find
	// the Node that represents the tag. Avoid using goquery pkg for
	// simplicity.
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	fmt.Printf("raw %v\n", string(data))
	n := doc.FirstChild
	if n == nil || n.Data != "html" {
		return nil, ErrUnexpectedTagLoc
	}

	n = doc.FirstChild.LastChild
	if n == nil || n.Data != "body" {
		return nil, ErrUnexpectedTagLoc
	}

	n = doc.FirstChild.LastChild.FirstChild
	if n == nil {
		return nil, ErrUnexpectedTagLoc
	}

	fmt.Printf("node %#v\n", *n)
	return n, nil
}

func getNodeAttrPtr(n *html.Node, name string) *html.Attribute {
	if n == nil {
		return nil
	}

	for i, a := range n.Attr {
		if a.Key == name {
			return &n.Attr[i]
		}
	}

	return nil
}

func getNodeAttr(n *html.Node, name string) string {
	a := getNodeAttrPtr(n, name)
	if a == nil {
		return ""
	}

	return a.Val
}

func setNodeAttr(n *html.Node, name, val string) {
	a := getNodeAttrPtr(n, name)
	if a == nil {
		n.Attr = append(n.Attr, html.Attribute{Key: name, Val: val})
	} else {
		a.Val = val
	}
}

func getNodeOnlyAttr(n *html.Node) string {
	if len(n.Attr) != 1 {
		return ""
	}
	return n.Attr[0].Key
}
