//go:build windows

package main

import (
	"fmt"
	"os"

	"github.com/CarsonSlovoka/go-font/v1/lib/ttlib/ttfont"
	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"
)

func outputProcess(f *ttfont.File) {
	ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))

	for tagName := range f.Tables {
		table := f.Tables[tagName]
		if table == nil {
			continue
		}
		fmt.Printf("\n@%s\n", tagName)
		outData, err := ts.Format(table)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(outData)
	}
}
