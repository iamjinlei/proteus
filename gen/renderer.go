package gen

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

func renderPage(page *page, cfg Config) ([]byte, error) {
	flags := html.CommonFlags | html.HrefTargetBlank
	if cfg.LazyImageLoading {
		flags |= html.LazyLoadImages
	}

	rh := &renderHook{
		styles: cfg.Styles,
	}

	r := html.NewRenderer(
		html.RendererOptions{
			Flags:          flags,
			RenderNodeHook: rh.render,
		},
	)
	rh.r = r

	data := markdown.Render(page.root, r)
	if rh.err != nil {
		return nil, rh.err
	}

	return data, nil
}

type renderHook struct {
	styles Styles
	r      *html.Renderer
	err    error
}

func (h *renderHook) renderCodeBlock(
	w io.Writer,
	node *ast.CodeBlock,
	entering bool,
) (ast.WalkStatus, bool) {
	io.WriteString(w, fmt.Sprintf("<div style=\"%v\">", h.styles.CodeBlock))
	h.r.CodeBlock(w, node)
	io.WriteString(w, "</div>")
	return ast.GoToNext, true /* skip current node */
}

func (h *renderHook) render(
	w io.Writer,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, bool) {
	if false {
		name := reflect.TypeOf(node).String()
		if strings.Contains(name, "ListItem") ||
			strings.Contains(name, "Text") ||
			strings.Contains(name, "Paragraph") ||
			strings.Contains(name, "List") ||
			strings.Contains(name, "HTMLSpan") ||
			strings.Contains(name, "CodeBlock") ||
			strings.Contains(name, "Heading") {
		} else {
			fmt.Printf("node type = %v, entering %v\n",
				reflect.TypeOf(node).String(),
				entering,
			)
		}
	}

	switch v := node.(type) {
	case *ast.CodeBlock:
		return h.renderCodeBlock(w, v, entering)
	case *ast.TableCell:
		//return h.renderTabelCell(w, v, entering)
	}

	return ast.GoToNext, false
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
