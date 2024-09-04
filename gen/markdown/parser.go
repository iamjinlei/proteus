package markdown

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"
)

type Doc struct {
	Root     ast.Node
	Refs     []string
	Headings []*Heading
}

type Heading struct {
	Level    int
	ID       string
	Name     string
	Children []*Heading
}

type Parser struct {
	interalHtmlRefSuffix string
	lazyImageLoading     bool
}

func NewParser(
	interalHtmlRefSuffix string,
	lazyImageLoading bool,
) *Parser {
	return &Parser{
		interalHtmlRefSuffix: interalHtmlRefSuffix,
		lazyImageLoading:     lazyImageLoading,
	}
}

func (p *Parser) Parse(src []byte, relDir string) (*Doc, error) {
	mdp := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)
	root := mdp.Parse(src)

	c, err := p.buildMarkdownContent(root, relDir)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *Parser) buildMarkdownContent(
	root ast.Node,
	relPath string,
) (*Doc, error) {
	var walkErr error
	// Accumulate references found in the doc.
	var refs []string
	// Track headings to build table of contents.
	ht := newHeadingTracker()

	ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
		if true {
			name := reflect.TypeOf(node).String()
			if strings.Contains(name, "HTMLSpan") {
				n := node.(*ast.HTMLSpan)
				fmt.Printf("HTMLSpan: %v, entering %v\n", string(n.Literal), entering)
			} else if strings.Contains(name, "Text") ||
				strings.Contains(name, "ListItem") ||
				strings.Contains(name, "Paragraph") ||
				strings.Contains(name, "List") ||
				strings.Contains(name, "HTMLSpan") ||
				strings.Contains(name, "Heading") {
			} else {
				fmt.Printf("node type = %v, entering %v\n",
					name,
					entering,
				)
			}
		}

		if !entering {
			return ast.GoToNext
		}

		switch v := node.(type) {
		case *ast.Heading:
			if len(v.Children) != 1 {
				break
			}

			ht.add(v.Level, v.HeadingID, string(v.Children[0].(*ast.Text).Literal))

		case *ast.Link:
			ref := string(v.Destination)
			if isExternalLink(ref) {
				break
			}

			if !strings.HasPrefix(ref, relPath) {
				ref = filepath.Join(relPath, ref)
			}

			v.Destination = []byte(ref + p.interalHtmlRefSuffix)
			refs = append(refs, ref)

		case *ast.Image:
			ref := string(v.Destination)
			if isExternalLink(ref) {
				break
			}

			if !strings.HasPrefix(ref, relPath) {
				ref = filepath.Join(relPath, ref)
			}

			refs = append(refs, ref)

		case *ast.HTMLSpan:
			if bytes.HasPrefix(v.Literal, []byte("</")) {
				break
			}

			// HTMLSpan is not limited to <span> tag, it actually represents
			// a set of HTML tags, such as span, img, etc.
			htmlBody, htmlRefs, err := p.processHTMLTag(relPath, v.Literal)
			if err != nil {
				walkErr = err
				return ast.Terminate
			}

			v.Literal = htmlBody
			refs = append(refs, htmlRefs...)
		}

		return ast.GoToNext
	})

	if walkErr != nil {
		return nil, walkErr
	}

	return &Doc{
		Root:     root,
		Refs:     refs,
		Headings: ht.getHeadings(),
	}, nil
}

func (p *Parser) processHTMLTag(
	relPath string,
	data []byte,
) ([]byte, []string, error) {
	tag, err := parseTag(data)
	if err != nil {
		return nil, nil, err
	}

	var refs []string
	hasUpdate := false
	switch tag.Data {
	case "img":
		ref := getTagAttr(tag, "src")
		if ref == "" {
			break
		}

		if isExternalLink(ref) {
			break
		}

		if p.lazyImageLoading {
			setTagAttr(tag, "loading", "lazy")
			hasUpdate = true
		}

		if !strings.HasPrefix(ref, relPath) {
			ref = filepath.Join(relPath, ref)
		}

		refs = append(refs, ref)
	}

	if hasUpdate {
		// NOTE(kmax): be careful, <span> and </span> are treated as
		// different ast.Node. So we only parse <span>, but doc.Html
		// would complete the closing </span>. If we don't handle it
		// properly, it may end up with 2 </span>. Luckily, for now
		// <img> does not have the closing </img> pair.
		var buf bytes.Buffer
		if err := html.Render(&buf, tag); err != nil {
			return nil, nil, err
		}
		data = buf.Bytes()
		fmt.Printf("updated %v\n", string(data))
	}

	return data, refs, nil
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
