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
	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"

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
	_, _ = f.GetTable(tag.Head)
	_, _ = f.GetTable(tag.Maxp)
	_, _ = f.GetTable(tag.OS2)

	// 由於GetTable相依的表也會載入，所以移除不需要的表
	// f.Tables[tag.Cmap] = nil
	f.Tables[tag.Glyf] = nil
	f.Tables[tag.Loca] = nil
	f.Tables[tag.Hmtx] = nil
	// f.Tables[tag.Post] = nil

	// 移除其他數據
	f.GlyphOrder = nil // 這個數據也很大，如果不想要也可以拿掉

	// outData, _ := f.ExportJSON()

	ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))
	outData, err := ts.Format(f.TTFont)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(outData)
}
