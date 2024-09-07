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
)

const (
	renderSkip = true
	renderNode = false
)

type Styles struct {
	CodeBlock string
}

type Renderer struct {
	palette               color.Palette
	colorMap              map[string]string
	styles                Styles
	internalRefHtmlSuffix string
	lazyImageLoading      bool
	state                 *renderState
}

func NewRenderer(
	palette color.Palette,
	styles Styles,
	internalRefHtmlSuffix string,
	lazyImageLoading bool,
) *Renderer {
	cm := map[string]string{}
	types := reflect.TypeOf(palette)
	vals := reflect.ValueOf(palette)
	for i := 0; i < types.NumField(); i++ {
		cm[strings.ToLower(types.Field(i).Name)] = vals.Field(i).String()
	}

	cm["name"] = string(palette.HighlighterRed)
	cm["b"] = string(palette.HighlighterGreen)
	cm["c"] = string(palette.HighlighterBlue)
	cm["d"] = string(palette.HighlighterYellow)
	cm["e"] = string(palette.HighlighterOrange)

	return &Renderer{
		palette:               palette,
		colorMap:              cm,
		styles:                styles,
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
	err          error
}

func (r *Renderer) Render(root ast.Node) (*Doc, error) {
	flags := html.CommonFlags | html.HrefTargetBlank
	if r.lazyImageLoading {
		flags |= html.LazyLoadImages
	}

	r.state = &renderState{
		htmlTagStack: newHtmlTagStack(),
		ht:           newHeadingTracker(),
		renderer: html.NewRenderer(
			html.RendererOptions{
				Flags:          flags,
				RenderNodeHook: r.render,
			},
		),
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

	if false {
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

func (r *Renderer) renderCodeBlock(
	w io.Writer,
	n *ast.CodeBlock,
	entering bool,
) ast.WalkStatus {
	fmt.Fprintf(w, `<div style="%v">`, r.styles.CodeBlock)
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
		name, val := getTagOnlyAttr(tag)
		if name == "" {
			break
		}

		if color, found := r.colorMap[name]; !found {
			break
		} else {
			r.state.htmlTagStack.push(
				htmlClosingTagMark,
				func(b *htmlTag) ast.WalkStatus {
					content := b.buf.String()
					switch val {
					case "baike", "baidu":
						content = link(content, baiduBaike(content))
					case "wikicn":
						content = link(content, wikipediaCn(content))
					}
					highlight(w, content, color)
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
