// go test -v

package tst_test

import (
	"fmt"

	"github.com/CarsonSlovoka/otparser.nvim/internal/tst"
)

type Font struct {
	Name   string `json:"name"`
	Tables Tables `json:"tables"`
}

type Tables struct {
	Head Head      `json:"head"`
	Name NameTable `json:"name"`
}

type Head struct {
	Version    float32 `json:"version"`
	UnitsPerEm int     `json:"unitsPerEm"`
}

type NameTable struct {
	Records []Record `json:"records"`
}

type Record struct {
	PlatformID int    `json:"platformID"`
	EncodingID int    `json:"encodingID"`
	NameID     int    `json:"nameID"`
	String     string `json:"string"`
}

func Example_format() {
	// 範例數據
	font := Font{
		Name: "MyFont",
		Tables: Tables{
			Head: Head{Version: 1.0, UnitsPerEm: 1000},
			Name: NameTable{
				Records: []Record{
					{PlatformID: 3, EncodingID: 1, NameID: 1, String: "MyFont"},
					{PlatformID: 3, EncodingID: 1, NameID: 2, String: "Regular"},
				},
			},
		},
	}

	// 配置輸出格式
	ts := tst.New(&tst.Config{
		Delimiter:      " | ",
		HeaderPrefix:   "\n# ",
		KeyValuePrefix: "- ",
	})

	output, err := ts.Format(font)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(output)
	// Output:
	// - name: MyFont
	//
	// # tables.head
	// - version: 1.0
	// - unitsPerEm: 1000
	//
	// # tables.name.records: platformID | encodingID | nameID | string
	// 3 | 1 | 1 | MyFont
	// 3 | 1 | 2 | Regular
}
