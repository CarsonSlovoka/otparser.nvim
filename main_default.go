//go:build !windows

package main

import (
	"fmt"
	"github.com/CarsonSlovoka/go-font/v1/lib/ttlib/ttfont"
	"github.com/CarsonSlovoka/go-font/v1/type/tag"
	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"
	"os"
)

func outputProcess(f *ttfont.File) {
	f.Tables[tag.Cmap] = nil
	f.Tables[tag.Glyf] = nil
	f.Tables[tag.Loca] = nil
	f.Tables[tag.Hmtx] = nil
	f.Tables[tag.Post] = nil

	// 移除其他數據
	f.GlyphOrder = nil // 這個數據也很大，如果不想要也可以拿掉

	ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))
	outData, err := ts.Format(f.TTFont)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(outData)
}
