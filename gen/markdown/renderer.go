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
	xhtml "golang.org/x/net/html"

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

func (r *Renderer) Render(
	root ast.Node,
) (template.HTML, error) {
	flags := html.CommonFlags | html.HrefTargetBlank
	if r.lazyImageLoading {
		flags |= html.LazyLoadImages
	}

	rh := &renderHook{
		palette:               r.palette,
		colorMap:              r.colorMap,
		styles:                r.styles,
		internalRefHtmlSuffix: r.internalRefHtmlSuffix,
		lazyImageLoading:      r.lazyImageLoading,
	}

	htmlRenderer := html.NewRenderer(
		html.RendererOptions{
			Flags:          flags,
			RenderNodeHook: rh.render,
		},
	)
	rh.r = htmlRenderer

	data := markdown.Render(root, htmlRenderer)
	if rh.err != nil {
		return template.HTML(""), rh.err
	}

	return template.HTML(data), nil
}

type htmlTagBuf struct {
	w          io.Writer
	openingTag *xhtml.Node
	buf        bytes.Buffer
	closingTag []byte
	fn         func(*htmlTagBuf) ast.WalkStatus
}

func (b *htmlTagBuf) close() ast.WalkStatus {
	return b.fn(b)
}

type renderHook struct {
	palette               color.Palette
	colorMap              map[string]string
	styles                Styles
	internalRefHtmlSuffix string
	lazyImageLoading      bool
	r                     *html.Renderer
	reentry               bool
	tagBufStack           []*htmlTagBuf
	err                   error
}

func (h *renderHook) render(
	w io.Writer,
	n ast.Node,
	entering bool,
) (ast.WalkStatus, bool) {
	if h.reentry {
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
	if len(h.tagBufStack) > 0 {
		w = &h.tagBufStack[len(h.tagBufStack)-1].buf
	}

	// A non-reentry call should never return false to use default Render
	// which uses the default writer. Always use h.renderNodeDefault() to
	// trigger a reentry call to the default Render, with picked writer.
	switch v := n.(type) {
	case *ast.CodeBlock:
		return h.renderCodeBlock(w, v, entering), renderSkip

	case *ast.Link:
		ref := string(v.Destination)
		if !isExternalLink(ref) {
			v.Destination = []byte(ref + h.internalRefHtmlSuffix)
		}

	case *ast.HTMLSpan:
		return h.processHTMLTag(w, v, entering), renderSkip
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
	return h.renderNodeDefault(w, n, entering), renderSkip
}

func (h *renderHook) renderNodeDefault(
	w io.Writer,
	n ast.Node,
	entering bool,
) ast.WalkStatus {
	h.reentry = true
	s := h.r.RenderNode(w, n, entering)
	h.reentry = false
	return s
}

func (h *renderHook) renderCodeBlock(
	w io.Writer,
	n *ast.CodeBlock,
	entering bool,
) ast.WalkStatus {
	fmt.Fprintf(w, `<div style="%v">`, h.styles.CodeBlock)
	h.r.CodeBlock(w, n)
	fmt.Fprintf(w, "</div>")
	return ast.GoToNext
}

func (h *renderHook) processHTMLTag(
	w io.Writer,
	n *ast.HTMLSpan,
	entering bool,
) ast.WalkStatus {
	if bytes.HasPrefix(n.Literal, htmlClosingTagPrefix) {
		return h.processHTMLClosingTag(w, n, entering)
	}
	return h.processHTMLOpeningTag(w, n, entering)
}

func (h *renderHook) processHTMLOpeningTag(
	w io.Writer,
	n *ast.HTMLSpan,
	entering bool,
) ast.WalkStatus {
	tag, err := parseTag(n.Literal)
	if err != nil {
		h.err = err
		return ast.Terminate
	}

	switch tag.Data {
	case "img":
		if !h.lazyImageLoading {
			break
		}

		setTagAttr(tag, "loading", "lazy")
		if v, err := renderTag(tag); err != nil {
			h.err = err
			return ast.Terminate
		} else {
			n.Literal = v
		}

	case "ins":
		switch getTagAttr(tag, "type") {
		case "book_bib":
			h.tagBufStack = append(h.tagBufStack, &htmlTagBuf{
				w:          w,
				openingTag: tag,
				closingTag: htmlClosingTagIns,
				fn: func(b *htmlTagBuf) ast.WalkStatus {
					// Content inside ins tag is ignored.
					bookBibliography(
						b.w,
						h.palette,
						getTagAttr(b.openingTag, "title"),
						getTagAttr(b.openingTag, "cover"),
						getTagAttr(b.openingTag, "link"),
						getTagAttr(b.openingTag, "author"),
					)
					return ast.GoToNext
				},
			})

			return ast.GoToNext
		}

	case "mark":
		name, _ := getTagOnlyAttr(tag)
		if name == "" {
			break
		}

		if color, found := h.colorMap[name]; !found {
			break
		} else {
			h.tagBufStack = append(h.tagBufStack, &htmlTagBuf{
				w:          w,
				openingTag: tag,
				closingTag: htmlClosingTagMark,
				fn: func(b *htmlTagBuf) ast.WalkStatus {
					highlight(b.w, b.buf.Bytes(), color)
					return ast.GoToNext
				},
			})
		}

		return ast.GoToNext
	}

	return h.renderNodeDefault(w, n, entering)
}

func (h *renderHook) processHTMLClosingTag(
	w io.Writer,
	n *ast.HTMLSpan,
	entering bool,
) ast.WalkStatus {
	if len(h.tagBufStack) == 0 {
		return h.renderNodeDefault(w, n, entering)
	}

	tb := h.tagBufStack[len(h.tagBufStack)-1]
	if !bytes.Equal(n.Literal, tb.closingTag) {
		return h.renderNodeDefault(w, n, entering)
	}

	h.tagBufStack = h.tagBufStack[:len(h.tagBufStack)-1]

	return tb.close()
}
