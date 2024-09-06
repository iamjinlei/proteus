package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrUnexpectedTagLoc = errors.New("unexpected tag location")
)

var (
	htmlClosingTagPrefix = []byte("</")
	htmlClosingTagMark   = []byte("</mark>")
	htmlClosingTagIns    = []byte("</ins>")
	htmlClosingTagDiv    = []byte("</div>")
	htmlClosingTagSpan   = []byte("</span>")
	htmlClosingTagASpan  = []byte("</a></span>")
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

	return n, nil
}

func renderTag(n *html.Node) ([]byte, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getTagAttrPtr(n *html.Node, name string) *html.Attribute {
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

func getTagAttr(n *html.Node, name string) string {
	a := getTagAttrPtr(n, name)
	if a == nil {
		return ""
	}

	return a.Val
}

func setTagAttr(n *html.Node, name, val string) {
	a := getTagAttrPtr(n, name)
	if a == nil {
		n.Attr = append(n.Attr, html.Attribute{Key: name, Val: val})
	} else {
		a.Val = val
	}
}

func getTagOnlyAttr(n *html.Node) (string, string) {
	if len(n.Attr) != 1 {
		return "", ""
	}
	return n.Attr[0].Key, n.Attr[0].Val
}

func parseStyle(style string) map[string]string {
	m := map[string]string{}
	parts := strings.Split(style, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			continue
		}

		m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	return m
}

func encodeStyle(style map[string]string) string {
	var keys []string
	for k, _ := range style {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var arr []string
	for _, k := range keys {
		arr = append(arr, fmt.Sprintf("%s:%s", k, style[k]))
	}

	return strings.Join(arr, ";")
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
