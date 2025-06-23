//go:build windows

package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/CarsonSlovoka/go-font/v1/lib/ttlib/ttfont"
	"github.com/CarsonSlovoka/go-font/v1/type/tag"
	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"
)

func outputProcess(f *ttfont.File) {
	f.Tables[tag.Cmap] = nil
	f.Tables[tag.Glyf] = nil
	f.Tables[tag.Loca] = nil
	f.Tables[tag.Hmtx] = nil
	f.Tables[tag.Post] = nil

	// 移除其他數據
	f.GlyphOrder = nil // 這個數據也很大，如果不想要也可以拿掉
	for tagName := range f.Tables {
		table := f.Tables[tagName]
		if table == nil {
			continue
		}
		ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))
		fmt.Printf("\n@%s\n", tagName)
		outData, err := ts.Format(table)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(outData)
	}
}
