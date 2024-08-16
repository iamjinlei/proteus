package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/iamjinlei/proteus/gen"
)

func main() {
	genFlag := flag.Bool(
		"g",
		false,
		"If true, generate html files from the markdown files instead of serving",
	)
	dirFlag := flag.String(
		"d",
		"",
		"Directory that contains markdown files to serve or convert",
	)
	flag.Parse()

	if *genFlag {
	} else {
		g := gen.NewHtml()

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			fmt.Printf("Path = %v\n", r.URL.Path)
			if path == "" || path == "/" {
				path = "/index"
			}

			data, err := os.ReadFile(filepath.Join(*dirFlag, path+".md"))
			if err != nil {
				w.Write([]byte("Not Found"))
				fmt.Printf("Error reading file: %v\n", err)
				return
			}

			html := g.Gen(data)
			fmt.Printf("Transformed, %v bytes\n", len(html))
			w.Write(html)
		})

		fmt.Printf("Start serving at http://localhost:8080\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Error serving http: %v\n", err)
			return
		}
	}
}
