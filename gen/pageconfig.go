package gen

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	markdownCommentOpen  = []byte("<!---")
	markdownCommentClose = []byte("--->")

	ErrBrokenCommentTag = errors.New("broken comment tag pair")
)

func extractPageConfig(src []byte) (*pageConfig, []byte, error) {
	src = bytes.TrimSpace(src)
	if !bytes.HasPrefix(src, markdownCommentOpen) {
		return newPageConfig(nil), src, nil
	}

	right := bytes.Index(src, markdownCommentClose)
	if right == -1 {
		return nil, nil, ErrBrokenCommentTag
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(
		src[len(markdownCommentOpen):right],
		&cfg,
	); err != nil {
		return nil, nil, err
	}

	content := bytes.TrimSpace(src[right+len(markdownCommentClose):])
	return newPageConfig(cfg), content, nil
}

type pageConfig struct {
	m map[string]interface{}
}

func newPageConfig(m map[string]interface{}) *pageConfig {
	if m == nil {
		m = map[string]interface{}{}
	}
	return &pageConfig{
		m: m,
	}
}

func (c *pageConfig) bannerRef() string {
	if c.m["banner"] == nil {
		return ""
	}

	ref, ok := c.m["banner"].(string)
	if !ok {
		return ""
	}
	return ref
}

func (c *pageConfig) leftPane() string {
	if c.m["left_pane"] == nil {
		return ""
	}

	t, ok := c.m["left_pane"].(string)
	if !ok {
		return ""
	}
	return t
}

func (c *pageConfig) rightPane() string {
	if c.m["right_pane"] == nil {
		return ""
	}

	t, ok := c.m["right_pane"].(string)
	if !ok {
		return ""
	}
	return t
}

func (c *pageConfig) header() *HtmlComponent {
	if c.m["banner"] == nil {
		return &HtmlComponent{
			Html: template.HTML(""),
		}
	}

	return &HtmlComponent{
		Html: template.HTML(fmt.Sprintf(
			`<img src="%v" style="width:100%%;height:%s;object-fit:cover;">`,
			c.m["banner"],
			imgBannerHeight,
		)),
	}
}

func (c *pageConfig) nav() *HtmlComponent {
	if c.m["nav"] == nil {
		return &HtmlComponent{
			Html: template.HTML(""),
		}
	}

	arr, ok := c.m["nav"].([]interface{})
	if !ok {
		return &HtmlComponent{
			Html: template.HTML(""),
		}
	}

	var links []string
	for _, v := range arr {
		str, ok := v.(string)
		if !ok {
			return &HtmlComponent{
				Html: template.HTML(""),
			}
		}

		kv := strings.Split(str, "=")
		if len(kv) != 2 {
			return &HtmlComponent{
				Html: template.HTML(""),
			}
		}

		links = append(links, fmt.Sprintf(`<a href="%s">%s</a>`, kv[1], kv[0]))
	}

	return &HtmlComponent{
		Html: template.HTML(fmt.Sprintf(
			`<span>%s</span><span style="margin-left:1em;">%s</span>`,
			"\U0001F517",
			strings.Join(links, " | "),
		)),
	}
}

func (c *pageConfig) footer() *HtmlComponent {
	return &HtmlComponent{
		Html: template.HTML(`
	<div style="max-width:fit-content;margin-inline:auto;">
		<span style="font-size: 0.8em;">Generated from markdown by
		<a href="https://github.com/iamjinlei/proteus">proteus</a>
		</span>
		</div>
	</div>`),
	}
}
