package markdown

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
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

// Parser parses markdown document, extract and analyzes its content.
// Parser should be read-only, not modifying any AST content.
type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

// Link or reference used in the markdown can be relative the current file
// location, it is ok as browser appends the relative path and the server
// always receives the full path relative to the server root.
func (p *Parser) Parse(src []byte) (*Doc, error) {
	mdp := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)
	root := mdp.Parse(src)

	d, err := p.buildMarkdownDoc(root)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (p *Parser) buildMarkdownDoc(
	root ast.Node,
) (*Doc, error) {
	var walkErr error
	// Accumulate references found in the doc.
	var refs []string
	// Track headings to build table of contents.
	ht := newHeadingTracker()

	ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
		if false {
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
			if !isExternalLink(ref) {
				refs = append(refs, ref)
			}

		case *ast.Image:
			ref := string(v.Destination)
			if !isExternalLink(ref) {
				refs = append(refs, ref)
			}

		case *ast.HTMLSpan:
			// HTMLSpan is not limited to <span> tag, it actually represents
			// a set of HTML tags, such as span, img, etc.
			ref, err := p.processHTMLTag(v.Literal)
			if err != nil {
				walkErr = err
				return ast.Terminate
			}

			if ref != "" {
				refs = append(refs, ref)
			}
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
	data []byte,
) (string, error) {
	// Closing tag is also treated as entering.
	if bytes.HasPrefix(data, htmlClosingTagPrefix) {
		return "", nil
	}

	node, err := parseTag(data)
	if err != nil {
		return "", err
	}

	switch node.Data {
	case "img":
		ref := getTagAttr(node, "src")
		if ref == "" || isExternalLink(ref) {
			return "", nil
		}
		return ref, nil

	case "ins":
		switch getTagAttr(node, "type") {
		case "book_bib":
			coverImgRef := getTagAttr(node, "cover")
			if !isExternalLink(coverImgRef) {
				return coverImgRef, nil
			}

		default:
		}
	}

	return "", nil
}
