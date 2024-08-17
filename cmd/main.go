package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	if *genFlag {
		dstDir := filepath.Clean(*dstFlag)
		rel, err := filepath.Rel(dstDir, srcDir)
		fmt.Printf("rel = %v\n", rel)
		if err != nil {
			fmt.Printf("Error calculating relative path between source and destination directory: %v\n", err)
			return
		}

		g := gen.NewHtml(htmlSuffix)
		if err := filepath.Walk(
			srcDir,
			func(path string, fi fs.FileInfo, err error) error {
				fmt.Printf("seen path %v\n", path)
				if fi.IsDir() {
					path = strings.Replace(path, srcDir, dstDir, 1)
					return os.MkdirAll(path, dirPermMode)
				}

				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if strings.HasSuffix(path, mdSuffix) {
					fmt.Printf("seen md file %v\n", path)
					html := g.Gen(data)
					path = strings.Replace(path[:len(path)-3]+htmlSuffix, srcDir, dstDir, 1)
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
				path = "/index.html"
			}

			path = filepath.Join(srcDir, path)

			if strings.HasSuffix(path, htmlSuffix) {
				data, err := os.ReadFile(strings.Replace(path, htmlSuffix, mdSuffix, 1))
				if err != nil {
					w.Write([]byte("Not Found"))
					fmt.Printf("Error reading file %v: %v\n", path, err)
					return
				}

				html := g.Gen(data)
				fmt.Printf("Transformed, %v bytes\n", len(html))
				w.Write(html)
			} else {
				if _, err := os.Stat(path); err == nil {
					data, err := os.ReadFile(path)
					if err != nil {
						w.Write([]byte("Not Found"))
						fmt.Printf("Error reading file %v: %v\n", path, err)
						return
					}

					w.Write(data)
				}
			}
		})

		fmt.Printf("Start serving at http://localhost:8080\n")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Printf("Error serving http: %v\n", err)
			return
		}
	}
}
