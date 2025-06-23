//go:build !windows

package main

import (
	"fmt"
	"github.com/CarsonSlovoka/go-font/v1/lib/ttlib/ttfont"
	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"
	"os"
)

func outputProcess(f *ttfont.File) {
	ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))
	outData, err := ts.Format(f.TTFont)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(outData)
}
