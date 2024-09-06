package gen

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"sort"
	"strings"
	"time"
)

type SitemapURL struct {
	Loc     string    `xml:"loc"`
	LastMod time.Time `xml:"lastmod"`
}

type UrlSet struct {
	Xmlns string        `xml:"xmlns,attr"`
	Urls  []*SitemapURL `xml:"url"`
}

type Sitemap struct {
	base   string
	rels   []string
	UrlSet UrlSet `xml:"urlset"`
}

func NewSitemap(base string) *Sitemap {
	if !strings.HasPrefix(base, "http://") &&
		!strings.HasPrefix(base, "https://") {
		base = "https://" + base
	}

	return &Sitemap{
		base: base,
		UrlSet: UrlSet{
			Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		},
	}
}

func (m *Sitemap) Add(rel string) {
	for len(rel) > 0 && rel[0] == '/' {
		rel = rel[1:]
	}
	m.rels = append(m.rels, rel)
}

func (m *Sitemap) Gen() ([]byte, error) {
	now := time.Now()
	sort.Strings(m.rels)
	for _, rel := range m.rels {
		loc, _ := url.JoinPath(m.base, rel)
		m.UrlSet.Urls = append(m.UrlSet.Urls, &SitemapURL{
			Loc:     loc,
			LastMod: now,
		})
	}

	v, err := xml.Marshal(m.UrlSet)
	if err != nil {
		return nil, err
	}

	v = bytes.Replace(v, []byte("UrlSet"), []byte("urlset"), -1)
	return append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), v...), nil
}
