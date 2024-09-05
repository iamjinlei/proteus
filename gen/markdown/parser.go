package markdown

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"

	"github.com/iamjinlei/proteus/gen/color"
)

var (
	htmlClosingTagPrefix = []byte("</")
	htmlClosingTagMark   = []byte("</mark>")
	htmlClosingTagSpan   = []byte("</span>")
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

type closingTagMapping struct {
	from []byte
	to   []byte
}

type Parser struct {
	palette              map[string]string
	interalHtmlRefSuffix string
	lazyImageLoading     bool
	closingTagQueue      []*closingTagMapping
}

func NewParser(
	palette color.Palette,
	interalHtmlRefSuffix string,
	lazyImageLoading bool,
) *Parser {
	p := map[string]string{}
	types := reflect.TypeOf(palette)
	vals := reflect.ValueOf(palette)
	for i := 0; i < types.NumField(); i++ {
		p[strings.ToLower(types.Field(i).Name)] = vals.Field(i).String()
	}

	p["name"] = string(palette.HighlighterRed)
	p["b"] = string(palette.HighlighterGreen)
	p["c"] = string(palette.HighlighterBlue)
	p["d"] = string(palette.HighlighterYellow)
	p["e"] = string(palette.HighlighterOrange)

	return &Parser{
		palette:              p,
		interalHtmlRefSuffix: interalHtmlRefSuffix,
		lazyImageLoading:     lazyImageLoading,
	}
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

	c, err := p.buildMarkdownContent(root)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *Parser) buildMarkdownContent(
	root ast.Node,
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

			v.Destination = []byte(ref + p.interalHtmlRefSuffix)
			refs = append(refs, ref)

		case *ast.Image:
			ref := string(v.Destination)
			if isExternalLink(ref) {
				break
			}

			refs = append(refs, ref)

		case *ast.HTMLSpan:
			// HTMLSpan is not limited to <span> tag, it actually represents
			// a set of HTML tags, such as span, img, etc.
			tag, err := p.processHTMLTag(v.Literal)
			if err != nil {
				walkErr = err
				return ast.Terminate
			}

			v.Literal = tag.html
			if tag.ref != "" {
				refs = append(refs, tag.ref)
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

type htmlTag struct {
	html []byte
	ref  string
}

func (p *Parser) processHTMLTag(
	data []byte,
) (*htmlTag, error) {
	tag := &htmlTag{
		html: data,
	}

	// Closing tag is also treated as entering.
	if bytes.HasPrefix(data, htmlClosingTagPrefix) {
		if len(p.closingTagQueue) > 0 {
			m := p.closingTagQueue[0]
			if bytes.Equal(data, m.from) {
				tag.html = m.to
				p.closingTagQueue = p.closingTagQueue[1:]
			}
		}

		return tag, nil
	}

	node, err := parseTag(data)
	if err != nil {
		return nil, err
	}

	updateNode := func() error {
		// NOTE(kmax): be careful, <span> and </span> are treated as
		// different ast.Node. So we only parse <span>, but doc.Html
		// would complete the closing </span>. If we don't handle it
		// properly, it may end up with 2 </span>. Luckily, for now
		// <img> does not have the closing </img> pair.
		var buf bytes.Buffer
		if err := html.Render(&buf, node); err != nil {
			return err
		}

		tag.html = buf.Bytes()
		fmt.Printf("updated %v\n", string(tag.html))
		return nil
	}

	switch node.Data {
	case "img":
		ref := getNodeAttr(node, "src")
		if ref == "" {
			break
		}

		if isExternalLink(ref) {
			break
		}

		if p.lazyImageLoading {
			setNodeAttr(node, "loading", "lazy")
			if err := updateNode(); err != nil {
				return nil, err
			}
		}

		tag.ref = ref

	case "mark":
		attrName := getNodeOnlyAttr(node)
		if attrName == "" {
			break
		}

		if color, found := p.palette[attrName]; found {
			tag.html = []byte(fmt.Sprintf(
				`<span style="background-color:%s;">`,
				color,
			))
			p.closingTagQueue = append(p.closingTagQueue, &closingTagMapping{
				from: htmlClosingTagMark,
				to:   htmlClosingTagSpan,
			})
			break
		}
	}

	return tag, nil
}

func isExternalLink(ref string) bool {
	return strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://")
}
