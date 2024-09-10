package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"reflect"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"

	"github.com/iamjinlei/proteus/gen/color"
	"github.com/iamjinlei/proteus/gen/keyword"
)

const (
	renderSkip = true
	renderNode = false
)

type Styles struct {
	Code      string
	CodeBlock string
}

var (
	defaultStyles = Styles{
		Code: fmt.Sprintf(
			"padding-left:0.3em;padding-right:0.3em;background-color:%v;",
			color.LightGray,
		),
		CodeBlock: fmt.Sprintf(
			"padding:0.1em 1.5em;background-color:%v;",
			color.LightGray,
		),
	}
)

type Renderer struct {
	palette               color.Palette
	colorMap              map[string]color.Color
	internalRefHtmlSuffix string
	lazyImageLoading      bool
	state                 *renderState
}

func NewRenderer(
	palette color.Palette,
	internalRefHtmlSuffix string,
	lazyImageLoading bool,
) *Renderer {
	cm := map[string]color.Color{}
	types := reflect.TypeOf(palette)
	vals := reflect.ValueOf(palette)
	for i := 0; i < types.NumField(); i++ {
		cm[strings.ToLower(types.Field(i).Name)] = vals.Field(i).Interface().(color.Color)
	}

	cm["name"] = palette.HighlighterRed
	cm["b"] = palette.HighlighterGreen
	cm["c"] = palette.HighlighterBlue
	cm["d"] = palette.HighlighterYellow
	cm["e"] = palette.HighlighterOrange

	return &Renderer{
		palette:               palette,
		colorMap:              cm,
		internalRefHtmlSuffix: internalRefHtmlSuffix,
		lazyImageLoading:      lazyImageLoading,
	}
}

type renderState struct {
	renderer     *html.Renderer
	reentry      bool
	htmlTagStack *htmlTagStack
	internalRefs []string
	ht           *headingTracker
	kws          *Keywords
	err          error
}

func (r *Renderer) Render(root ast.Node) (*Doc, error) {
	flags := html.CommonFlags
	if r.lazyImageLoading {
		flags |= html.LazyLoadImages
	}

	r.state = &renderState{
		renderer: html.NewRenderer(
			html.RendererOptions{
				Flags:          flags,
				RenderNodeHook: r.render,
			},
		),
		htmlTagStack: newHtmlTagStack(),
		ht:           newHeadingTracker(),
		kws:          newKeywords(r.colorMap),
	}

	// Traverse AST using ast.WalkFunc()
	data := markdown.Render(root, r.state.renderer)
	rs := r.state
	r.state = nil

	if rs.err != nil {
		return nil, rs.err
	}

	return &Doc{
		Html:         template.HTML(data),
		InternalRefs: rs.internalRefs,
		Headings:     rs.ht.getHeadings(),
		Keywords:     rs.kws,
	}, nil
}

func (r *Renderer) render(
	w io.Writer,
	n ast.Node,
	entering bool,
) (ast.WalkStatus, bool) {
	if r.state.reentry {
		return ast.GoToNext, renderNode
	}

	if true {
		name := reflect.TypeOf(n).String()
		if strings.Contains(name, "ListItem") ||
			strings.Contains(name, "Text") ||
			strings.Contains(name, "Paragraph") ||
			strings.Contains(name, "List") ||
			strings.Contains(name, "HTMLSpan") ||
			strings.Contains(name, "CodeBlock") ||
			strings.Contains(name, "Heading") {
		} else {
			fmt.Printf("node type = %v, entering %v\n",
				reflect.TypeOf(n).String(),
				entering,
			)
		}
	}

	// If tag buf stack is not empty, write to the top of the stack.
	if !r.state.htmlTagStack.empty() {
		w = &r.state.htmlTagStack.top().buf
	}

	// A non-reentry call should never return false to use default Render
	// which uses the default writer. Always use r.renderNodeDefault() to
	// trigger a reentry call to the default Render, with picked writer.
	switch v := n.(type) {
	case *ast.Heading:
		if !entering {
			break
		}

		if len(v.Children) != 1 {
			break
		}

		r.state.ht.add(v.Level, v.HeadingID, string(v.Children[0].(*ast.Text).Literal))

	case *ast.Code:
		return r.renderCode(w, v, entering), renderSkip

	case *ast.CodeBlock:
		return r.renderCodeBlock(w, v, entering), renderSkip

	case *ast.Link:
		if !entering {
			break
		}

		ref := string(v.Destination)
		if !isExternalLink(ref) {
			r.state.internalRefs = append(r.state.internalRefs, ref)
			v.Destination = []byte(ref + r.internalRefHtmlSuffix)
		}

	case *ast.Image:
		if !entering {
			break
		}

		if v.Attribute == nil {
			v.Attribute = &ast.Attribute{}
		}
		if v.Attribute.Attrs == nil {
			v.Attribute.Attrs = map[string][]byte{}
		}
		v.Attribute.Attrs["style"] = []byte("width:100%;margin-top:0.5em;margin-bottom:0.5em;")

		ref := string(v.Destination)
		if !isExternalLink(ref) {
			r.state.internalRefs = append(r.state.internalRefs, ref)
		}

	case *ast.HTMLSpan:
		return r.processHTMLTag(w, v, entering), renderSkip
	}

	/*
	 * RenderNode()
	 *   |
	 *   +-> render() [reentry == false]
	 *   | 	   |
	 *   | 	   +-> Set retrny = true
	 *   | 	   |
	 *   | 	   +-> RenderNode()
	 *   |     |	 |
	 *   |     |	 +-> render() [reentry == true]
	 *   |     |     |   Do nothing, return ast.GoToNext, ## false ##
	 *   |     |	 |
	 *   |     |	 +-> do ## rendering ##
	 *   |     |
	 *   | 	   +-> Set retrny = false
	 *   |     |
	 *   | 	   +-> Return ast.GoToNext, ## true ##
	 *   |
	 *   +-> Skip ## rendering ##
	 */
	return r.renderNodeDefault(w, n, entering), renderSkip
}

func (r *Renderer) renderNodeDefault(
	w io.Writer,
	n ast.Node,
	entering bool,
) ast.WalkStatus {
	r.state.reentry = true
	s := r.state.renderer.RenderNode(w, n, entering)
	r.state.reentry = false
	return s
}

func (r *Renderer) renderCode(
	w io.Writer,
	n *ast.Code,
	entering bool,
) ast.WalkStatus {
	fmt.Fprintf(w, `<span style="%v">`, defaultStyles.Code)
	r.state.renderer.Code(w, n)
	fmt.Fprintf(w, "</span>")
	return ast.GoToNext
}

func (r *Renderer) renderCodeBlock(
	w io.Writer,
	n *ast.CodeBlock,
	entering bool,
) ast.WalkStatus {
	fmt.Fprintf(w, `<div style="%v">`, defaultStyles.CodeBlock)
	r.state.renderer.CodeBlock(w, n)
	fmt.Fprintf(w, "</div>")
	return ast.GoToNext
}

func (r *Renderer) processHTMLTag(
	w io.Writer,
	n *ast.HTMLSpan,
	entering bool,
) ast.WalkStatus {
	if !bytes.HasPrefix(n.Literal, htmlClosingTagPrefix) {
		return r.processHTMLOpeningTag(w, n, entering)
	}
	return r.processHTMLClosingTag(w, n, entering)
}

func (r *Renderer) processHTMLOpeningTag(
	w io.Writer,
	n *ast.HTMLSpan,
	entering bool,
) ast.WalkStatus {
	tag, err := parseTag(n.Literal)
	if err != nil {
		r.state.err = err
		return ast.Terminate
	}

	switch tag.Data {
	case "img":
		if !r.lazyImageLoading {
			break
		}

		setTagAttr(tag, "loading", "lazy")
		style := parseStyle(getTagAttr(tag, "style"))
		if style["margin-top"] == "" && style["margin-bottom"] == "" {
			style["margin-top"] = "0.5em"
			style["margin-bottom"] = "0.5em"
			setTagAttr(tag, "style", encodeStyle(style))
		}

		if v, err := renderTag(tag); err != nil {
			r.state.err = err
			return ast.Terminate
		} else {
			n.Literal = v
		}

	case "ins":
		switch getTagAttr(tag, "type") {
		case "book_bib":
			r.state.htmlTagStack.push(
				htmlClosingTagIns,
				func(b *htmlTag) ast.WalkStatus {
					// Content inside ins tag is ignored.
					bookBibliography(
						w,
						r.palette,
						getTagAttr(tag, "title"),
						getTagAttr(tag, "cover"),
						getTagAttr(tag, "link"),
						getTagAttr(tag, "author"),
					)
					return ast.GoToNext
				},
			)

			return ast.GoToNext
		}

	case "mark":
		kind, val := getTagOnlyAttr(tag)
		if kind == "" {
			break
		}

		if color, found := r.colorMap[kind]; !found {
			break
		} else {
			r.state.htmlTagStack.push(
				htmlClosingTagMark,
				func(b *htmlTag) ast.WalkStatus {
					content := b.buf.String()
					id := ""
					if keyword.ValidType(kind) {
						id = hash20([]byte(content))
						r.state.kws.add(keyword.Type(kind), content, id)
					}

					switch val {
					case "baike", "baidu":
						content = link(content, baiduBaike(content))
					case "wikicn":
						content = link(content, wikipediaCn(content))
					}
					highlight(w, id, content, color)

					return ast.GoToNext
				},
			)
		}

		return ast.GoToNext
	}

	return r.renderNodeDefault(w, n, entering)
}

func (r *Renderer) processHTMLClosingTag(
	w io.Writer,
	n *ast.HTMLSpan,
	entering bool,
) ast.WalkStatus {
	if r.state.htmlTagStack.empty() {
		return r.renderNodeDefault(w, n, entering)
	}

	top := r.state.htmlTagStack.top()
	if !bytes.Equal(n.Literal, top.closingTag) {
		return r.renderNodeDefault(w, n, entering)
	}

	r.state.htmlTagStack.pop()

	return top.close()
}
