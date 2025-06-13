// go run .
// go build -o otparser .
// sudo ln -siv $(realpath otparser) /usr/bin/
// which otparser

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/CarsonSlovoka/otparser.nvim/internal/app"
)

// OffsetTable 定義 TTF 文件的 Offset Table 結構
type OffsetTable struct {
	SfntVersion   uint32 // 字體版本 (0x00010000 for TrueType)
	NumTables     uint16 // 表數量
	SearchRange   uint16 // (最大 2 次方) * 16
	EntrySelector uint16 // Log2(最大 2 次方)
	RangeShift    uint16 // NumTables * 16 - SearchRange
}

// TableEntry 定義每個表的目錄項
type TableEntry struct {
	Tag      [4]byte // 表標籤 (如 "head", "name")
	Checksum uint32  // 表的校驗和
	Offset   uint32  // 表在文件中的偏移
	Length   uint32  // 表長度
}

func main() {
	log.Println("app version: ", app.Version) // log是輸出在stderr, 所以不會影響到stdout
	if len(os.Args) < 2 {
		fmt.Println("Usage: ttfparser <ttf_file>")
		os.Exit(1)
	}

	// 讀取 TTF 文件
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// 解析 Offset Table
	var offsetTable OffsetTable
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &offsetTable); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing offset table: %v\n", err)
		os.Exit(1)
	}

	// 輸出 Offset Table 資訊
	fmt.Printf("TTF Header Info for %s:\n", filename)
	fmt.Printf("SFNT Version: 0x%08X\n", offsetTable.SfntVersion)
	fmt.Printf("Number of Tables: %d\n", offsetTable.NumTables)
	fmt.Printf("Search Range: %d\n", offsetTable.SearchRange)
	fmt.Printf("Entry Selector: %d\n", offsetTable.EntrySelector)
	fmt.Printf("Range Shift: %d\n", offsetTable.RangeShift)
	fmt.Println("\nTable Directory:")

	// 解析 Table Directory
	for i := 0; i < int(offsetTable.NumTables); i++ {
		var entry TableEntry
		if err := binary.Read(reader, binary.BigEndian, &entry); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing table entry %d: %v\n", i, err)
			os.Exit(1)
		}
		fmt.Printf("Table %d:\n", i+1)
		fmt.Printf("    Tag: %s\n", string(entry.Tag[:]))
		fmt.Printf("    Checksum: 0x%08X\n", entry.Checksum)
		fmt.Printf("    Offset: %d\n", entry.Offset)
		fmt.Printf("    Length: %d\n", entry.Length)
	}
}
