//go:build windows

package main

import (
	"fmt"
	"github.com/CarsonSlovoka/go-font/v1/lib/ttlib/ttfont"
	"github.com/CarsonSlovoka/go-font/v1/type"
	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"
	"os"
)

type Font struct {
	Header  ttfont.TableHeader
	Records []t.TableRecord
}

func outputProcess(f *ttfont.File) {
	ts := tst.New(tst.NewConfig(" | ", "\n# ", "- "))

	var copyFont Font
	copyFont.Header = f.TTFont.TableHeader

	// tableRecords, _ := ts.Format(f.TTFont.TableRecords) -- 用不了會被報錯
	fmt.Printf("\n@Table\n\n")
	copyFont.Records = make([]t.TableRecord, len(f.TTFont.TableRecords))
	for i := range f.TTFont.TableRecords {
		copyFont.Records[i] = *f.TTFont.TableRecords[i]
	}
	headerStr, _ := ts.Format(copyFont)
	fmt.Println(headerStr)

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
