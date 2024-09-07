package gen

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"sort"
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
	domain string
	rels   []string
	UrlSet UrlSet `xml:"urlset"`
}

func NewSitemap(domain string) *Sitemap {
	return &Sitemap{
		domain: normalizeDomain(domain),
		UrlSet: UrlSet{
			Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		},
	}
}

func (m *Sitemap) Add(rel string) {
	m.rels = append(m.rels, normalizeRelPath(rel))
}

func (m *Sitemap) Gen() ([]byte, error) {
	now := time.Now()
	sort.Strings(m.rels)
	for _, rel := range m.rels {
		loc, _ := url.JoinPath(m.domain, rel)
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
