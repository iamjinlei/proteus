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
	flag.Parse()

	srcDir := filepath.Clean(*srcFlag)

	cfg, err := loadConfig(srcDir, filepath.Join(srcDir, "config.yaml"))
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if *genFlag {
		dstDir := filepath.Clean(*dstFlag)

		g := gen.NewHtml(htmlSuffix)

		refQueue := []string{
			cfg.Entry,
			cfg.Favicon,
		}
		dirSeen := map[string]bool{}
		for len(refQueue) > 0 {
			relPath := refQueue[0]
			refQueue = refQueue[1:]
			isMarkdown := strings.HasSuffix(relPath, mdSuffix)

			src := filepath.Join(srcDir, relPath)
			dst := filepath.Join(dstDir, relPath)
			if isMarkdown {
				if relPath == cfg.Entry {
					dst = filepath.Join(dstDir, "index.html")
				} else {
					dst += htmlSuffix
				}
			}

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
				relDir := filepath.Dir(relPath)
				doc, err := g.Gen(relDir, data)
				if err != nil {
					fmt.Printf("Error generating HTML doc: %v\n", err)
					return
				}

				data = doc.Html
				refQueue = append(refQueue, doc.Refs...)
			}

			if err := os.WriteFile(dst, data, filePermMode); err != nil {
				fmt.Printf("Error writing destination file %v: %v\n", src, err)
				return
			}
		}
	} else {
		g := gen.NewHtml(htmlSuffix)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			fmt.Printf("Path = %v\n", r.URL.Path)
			if path == "" || path == "/" {
				path = filepath.Join("/", cfg.Entry)
			}

			directCopy := true
			if strings.HasSuffix(path, htmlSuffix) {
				directCopy = false
				path = path[:len(path)-len(htmlSuffix)]
				if !strings.HasSuffix(path, mdSuffix) {
					path += mdSuffix
				}
			} else if strings.HasSuffix(path, mdSuffix) {
				directCopy = false
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
				doc, err := g.Gen("", data)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("Error generating html: %v", err)))
				} else {
					fmt.Printf("Transformed, %v bytes\n", len(doc.Html))
					w.Write(doc.Html)
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

	return cfg, nil
}
