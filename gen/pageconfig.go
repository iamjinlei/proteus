package gen

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	pageConfigDivider = []byte("+++\n")
)

func extractPageConfig(src []byte) (*pageConfig, []byte, error) {
	left := bytes.Index(src, pageConfigDivider)
	if left == -1 {
		return newPageConfig(nil), src, nil
	}

	s := src[left+len(pageConfigDivider):]
	right := bytes.Index(s, pageConfigDivider)
	if right == -1 {
		return newPageConfig(nil), src, nil
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(s[:right], &cfg); err != nil {
		return nil, nil, err
	}

	return newPageConfig(cfg), s[right+len(pageConfigDivider):], nil
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

func (c *pageConfig) header() []byte {
	if c.m["banner"] == nil {
		return []byte(fmt.Sprintf("<div style=\"width:100%%;height:%s;\"></div>", emptyBannerHeight))
	}

	return []byte(fmt.Sprintf("<img src=\"%v\" style=\"width:100%%;height:%s;object-fit:cover;\">", c.m["banner"], imgBannerHeight))
}

func (c *pageConfig) navi() []byte {
	if c.m["navi"] == nil {
		return []byte("")
	}

	arr, ok := c.m["navi"].([]interface{})
	if !ok {
		return []byte("")
	}

	var links []string
	for _, v := range arr {
		str, ok := v.(string)
		if !ok {
			return []byte("")
		}
		fmt.Printf("%v\n", str)
		kv := strings.Split(str, "=")
		if len(kv) != 2 {
			return []byte("")
		}

		links = append(links, fmt.Sprintf("<a href=\"%s\">%s</a>", kv[1], kv[0]))
	}

	return []byte(fmt.Sprintf("<span>%s</span><span style=\"margin-left:1em;\">%s</span>",
		"\U0001f517",
		strings.Join(links, " | "),
	))
	//return []byte(fmt.Sprintf("< <a href=\"%s\">%s</a>", c.m["navi"]))
}

func (c *pageConfig) footer() []byte {
	return []byte(`
	<div style="max-width:fit-content;margin-inline:auto;">
		<span style="font-size: 80%;">Generated from markdown by
		<a href="https://github.com/iamjinlei/proteus">proteus</a>
		</span>
		</div>
	</div>`)
}
