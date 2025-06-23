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
	fmt.Printf("@TableHeader\n\n")
	tableHeader, _ := ts.Format(f.TTFont.TableHeader)
	fmt.Println(tableHeader)

	// tableRecords, _ := ts.Format(f.TTFont.TableRecords) -- 用不了會被報錯
	fmt.Printf("\n@TableRecords\n\n")
	for _, tr := range f.TTFont.TableRecords {
		tableRecords, _ := ts.Format(tr)
		fmt.Println(tableRecords)
	}

	for tagName := range f.Tables {
		table := f.Tables[tagName]
		if table == nil {
			continue
		}
		fmt.Printf("\n@%s\n\n", tagName)
		outData, err := ts.Format(table)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error format with TStruct: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(outData)
	}
}
