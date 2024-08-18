package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
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

	configPath := filepath.Join(srcDir, "config.yaml")
	cfg, err := loadConfig(srcDir, configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if *genFlag {
		dstDir := filepath.Clean(*dstFlag)
		entryPath := filepath.Join(srcDir, cfg.Entry)

		g := gen.NewHtml(htmlSuffix)
		if err := filepath.Walk(
			srcDir,
			func(path string, fi fs.FileInfo, err error) error {
				fmt.Printf("seen path %v, err %v\n", path, err)
				if err != nil {
					return err
				}

				if path == configPath {
					return nil
				}

				isHidden := isHidden(path)
				if fi.IsDir() {
					if isHidden {
						return filepath.SkipDir
					}

					path = strings.Replace(path, srcDir, dstDir, 1)
					return os.MkdirAll(path, dirPermMode)
				} else if isHidden {
					return nil
				}

				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if strings.HasSuffix(path, mdSuffix) {
					fmt.Printf("seen md file %v\n", path)

					if path == entryPath {
						path = filepath.Join(dstDir, "index.html")
					} else {
						path = path + htmlSuffix
						path = strings.Replace(path, srcDir, dstDir, 1)
					}

					html := g.Gen(data)
					return os.WriteFile(path, html, filePermMode)
				} else {
					fmt.Printf("copy file %v\n", path)
					path = strings.Replace(path, srcDir, dstDir, 1)
					return os.WriteFile(path, data, filePermMode)
				}
			},
		); err != nil {
			fmt.Printf("Error generating HTML files: %v\n", err)
			return
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
				html := g.Gen(data)
				fmt.Printf("Transformed, %v bytes\n", len(html))
				w.Write(html)
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
