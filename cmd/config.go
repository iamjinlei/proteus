package main

type Config struct {
	Domain        string            `yaml:"domain"`
	EnableSitemap bool              `yaml:"enable_sitemap"`
	Entry         string            `yaml:"entry"`
	Assets        map[string]string `yaml:"assets"`
}
