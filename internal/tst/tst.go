// TStruct
// T 可表代Text, Table等概念

package tst

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Config 定義輸出格式的客製化選項
type Config struct {
	Delimiter      string // 陣列數據的分隔符，例如 "|"
	HeaderPrefix   string // 標頭行前綴，例如 "# "
	KeyValuePrefix string // 鍵值對縮進前綴，例如 "  "
}

func NewConfig(sep, header, prefix string) *Config {
	return &Config{sep, header, prefix}
}

func New(c *Config) *TStruct {
	return &TStruct{*c}
}

type TStruct struct {
	Config
}

// Format 將任意 struct 轉換為客製化格式
func (t *TStruct) Format(v any) (string, error) {
	var buf bytes.Buffer
	err := writeStruct(&buf, reflect.ValueOf(v), "", t.Config)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// writeStruct 遞迴處理 struct
func writeStruct(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", v.Kind())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 獲取 JSON 標籤名稱，若無則使用欄位名稱
		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}
		if name == "-" {
			continue // 忽略標記為 "-" 的欄位
		}

		// 處理嵌套 struct
		if value.Kind() == reflect.Struct {
			newPrefix := prefix
			if prefix != "" {
				newPrefix += "."
			}
			newPrefix += name
			if err := writeStruct(buf, value, newPrefix, config); err != nil {
				return err
			}
			continue
		}

		// 處理陣列或切片
		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			if err := writeSlice(buf, value, prefix+"."+name, config); err != nil {
				return err
			}
			continue
		}

		// 處理簡單欄位（鍵值對）
		if prefix != "" {
			// 只有在第一個鍵值對前添加上下文標記
			if buf.Len() == 0 || buf.Bytes()[buf.Len()-1] != '\n' {
				fmt.Fprintf(buf, "%s%s\n", config.HeaderPrefix, prefix)
			}
			fmt.Fprintf(buf, "%s%s: %v\n", config.KeyValuePrefix, name, value.Interface())
		}
	}

	return nil
}

// writeSlice 處理陣列或切片，生成單標頭格式
func writeSlice(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	if v.Len() == 0 {
		return nil
	}

	// 假設切片元素是 struct
	if v.Index(0).Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be structs, got %v", v.Index(0).Kind())
	}

	// 獲取欄位名稱
	t := v.Index(0).Type()
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		if name == "" {
			name = t.Field(i).Name
		}
		if name != "-" {
			headers = append(headers, name)
		}
	}

	// 寫入標頭
	fmt.Fprintf(buf, "%s%s: %s\n", config.HeaderPrefix, prefix, strings.Join(headers, config.Delimiter))

	// 寫入數據
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		var values []string
		for j := 0; j < item.Type().NumField(); j++ {
			if item.Type().Field(j).Tag.Get("json") == "-" {
				continue
			}
			val := item.Field(j).Interface()
			switch v := val.(type) {
			case float32, float64:
				values = append(values, strconv.FormatFloat(v.(float64), 'f', 1, 64))
			default:
				values = append(values, fmt.Sprintf("%v", val))
			}
		}
		fmt.Fprintf(buf, "%s\n", strings.Join(values, config.Delimiter))
	}

	fmt.Fprintln(buf) // 陣列後添加空行
	return nil
}
