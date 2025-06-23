// go run .
// go build -o otparser .
// sudo ln -siv $(realpath otparser) /usr/bin/
// which otparser

package main

import (
	"fmt"

	"github.com/CarsonSlovoka/go-font/v1"
	"github.com/CarsonSlovoka/go-font/v1/lib/ttlib/ttfont"
	"github.com/CarsonSlovoka/go-font/v1/type/tag"
	"github.com/CarsonSlovoka/otparser.nvim/internal/app"
	"log"
	"os"
)

func init() {
	log.Println("app version: ", app.Version) // log是輸出在stderr, 所以不會影響到stdout
	if len(os.Args) < 2 {
		fmt.Println("Usage: otparser <ttf_file>")
		os.Exit(1)
	}
}

func main() {
	filepath := os.Args[1]

	iFont, err := font.Dump(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	f := iFont.(*ttfont.File)

	_, _ = f.GetTable(tag.Head)
	_, _ = f.GetTable(tag.Maxp)
	_, _ = f.GetTable(tag.OS2)
	outputProcess(f)
}
