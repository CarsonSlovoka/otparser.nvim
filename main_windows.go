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

func outputProcess_old(f *ttfont.File) {
	for key := range f.Tables {
		if !slices.Contains([]tag.Tag{
			tag.Maxp,
		}, key) {
			// f.Tables[key] = nil
			continue
		}
		ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))
		outData, err := ts.Format(f.Tables[key])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(outData) // 一次寫入在windows上，對記憶體的要求較高，如果取的表格太複雜在windows會出問題跑不出來
	}
}

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
