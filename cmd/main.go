package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/iamjinlei/proteus/gen"
)

const (
	htmlSuffix = ".html"
	mdSuffix   = ".md"

	dirPermMode  = 0755
	filePermMode = 0644
)

func main() {
	genFlag := flag.Bool(
		"g",
		false,
		"If true, generate html files from the markdown files instead of serving",
	)
	srcFlag := flag.String(
		"s",
		"",
		"Source directory that contains markdown files to serve or convert",
	)
	dstFlag := flag.String(
		"d",
		"",
		"Destination directory to store generated html files",
	)
	forceFlag := flag.Bool(
		"f",
		false,
		"Force regenerate all html files",
	)
	flag.Parse()

	srcDir := filepath.Clean(*srcFlag)

	cfg, err := loadConfig(srcDir, filepath.Join(srcDir, "config.yaml"))
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Printf("yaml domain = %v\n", cfg.Domain)
	g, err := gen.NewHtml(gen.DefaultConfig(
		cfg.Domain,
		htmlSuffix,
	))
	if err != nil {
		fmt.Printf("Error creating html renderer: %v\n", err)
		return
	}

	if *genFlag {
		dstDir := filepath.Clean(*dstFlag)
		sm := gen.NewSitemap(cfg.Domain)

		refQueue := []string{
			cfg.Entry,
		}
		for path, _ := range cfg.Assets {
			refQueue = append(refQueue, path)
		}

		mdCnt := 0
		dirSeen := map[string]bool{}
		for len(refQueue) > 0 {
			// relPath is the path relative to the source repo dir.
			relPath := refQueue[0]
			refQueue = refQueue[1:]
			isMarkdown := strings.HasSuffix(relPath, mdSuffix)

			src := filepath.Join(srcDir, relPath)
			dst := filepath.Join(dstDir, relPath)
			switch relPath {
			case cfg.Entry:
				dst = filepath.Join(dstDir, "index.html")
			default:
				if v := cfg.Assets[relPath]; v != "" {
					dst = filepath.Join(dstDir, v)
					isMarkdown = false
				}

				if isMarkdown {
					dst += htmlSuffix
				}
			}

			relPath = dst[len(dstDir):]
			if isMarkdown {
				sm.Add(relPath)
			} else if !*forceFlag && !updateRequired(src, dst) {
				continue
			}

			fmt.Printf("Processing %v -> %v\n", src, dst)

			dstDir := filepath.Dir(dst)
			if !dirSeen[dstDir] {
				if err := os.MkdirAll(dstDir, dirPermMode); err != nil {
					fmt.Printf("Error creating directory %v: %v\n", dstDir, err)
					return
				}
				dirSeen[dstDir] = true
			}

			data, err := os.ReadFile(src)
			if err != nil {
				fmt.Printf("Error reading source file %v: %v\n", src, err)
				return
			}

			if isMarkdown {
				page, err := g.Gen(relPath, data)
				if err != nil {
					fmt.Printf("Error generating HTML page: %v\n", err)
					return
				}

				relDir := filepath.Dir(relPath)
				for _, ref := range page.InternalRefs {
					// If the link path referenced in the markdown is relative
					// to the markdown file location, update it to be relative
					// to the src repo location before pushing into queue.
					if relDir != "" && !strings.HasPrefix(ref, relDir) {
						ref = filepath.Join(relDir, ref)
					}
					refQueue = append(refQueue, ref)
				}

				data = page.Html
				mdCnt++
			}

			if err := os.WriteFile(dst, data, filePermMode); err != nil {
				fmt.Printf("Error writing destination file %v: %v\n", src, err)
				return
			}
		}

		if cfg.EnableSitemap && cfg.Domain != "" {
			data, err := sm.Gen()
			if err != nil {
				fmt.Printf("Error generating sitemap file: %v\n", err)
				return
			}
			if err := os.WriteFile(
				filepath.Join(dstDir, "sitemap.xml"),
				data,
				filePermMode,
			); err != nil {
				fmt.Printf("Error writing sitemap file: %v\n", err)
				return
			}
		}

		fmt.Printf("Total markdown files processed: %v\n", mdCnt)
	} else {
		rassets := map[string]string{}
		for from, to := range cfg.Assets {
			rassets[to] = from
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			fmt.Printf("Path = %v\n", r.URL.Path)
			switch path {
			case "", "/", "/index.html":
				path = cfg.Entry + htmlSuffix
			default:
				if v := rassets[path]; v != "" {
					path = v
				}
			}

			directCopy := true
			if strings.HasSuffix(path, htmlSuffix) {
				directCopy = false
				path = path[:len(path)-len(htmlSuffix)]
				if !strings.HasSuffix(path, mdSuffix) {
					path += mdSuffix
				}
			}

			path = filepath.Clean(filepath.Join(srcDir, path))
			if !strings.HasPrefix(path, srcDir) || isHidden(path) {
				// Serve 404
				w.Write([]byte("Not Found"))
				return
			}

			data, err := os.ReadFile(path)
			if err != nil {
				// TODO(lei): serve 404
				w.Write([]byte("Not Found"))
				fmt.Printf("Error reading file %v: %v\n", path, err)
				return
			}

			if directCopy {
				w.Write(data)
			} else {
				page, err := g.Gen(path, data)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("Error generating html page: %v", err)))
				} else {
					fmt.Printf("Transformed, %v bytes\n", len(page.Html))
					w.Write(page.Html)
				}
			}
		})

		fmt.Printf("Start serving at http://localhost:8000\n")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Printf("Error serving http: %v\n", err)
			return
		}
	}
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func isHidden(path string) bool {
	base := filepath.Base(path)
	return base[0] == '.'
}

func updateRequired(src, dst string) bool {
	dstFi, err := os.Stat(dst)
	if err != nil {
		return true
	}
	srcFi, err := os.Stat(src)
	if err != nil {
		return true
	}

	return srcFi.ModTime().After(dstFi.ModTime())
}

func loadConfig(srcDir, path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if !fileExists(filepath.Join(srcDir, cfg.Entry)) {
		return Config{}, errors.New("entry point undefined")
	}

	cfg.Entry = filepath.Join("/", cfg.Entry)
	assets := map[string]string{}
	for from, to := range cfg.Assets {
		assets[filepath.Join("/", from)] = filepath.Join("/", to)
	}
	cfg.Assets = assets

	return cfg, nil
}
