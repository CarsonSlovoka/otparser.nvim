// TStruct
// T 可表代Text, Table等概念

package tst

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
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
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if err := writeStruct(&buf, val, "", t.Config); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// writeValue 處理任意值（struct、slice、map 等）
func writeValue(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return writeStruct(buf, v, prefix, config)
	case reflect.Slice, reflect.Array:
		return writeSlice(buf, v, prefix, config)
	case reflect.Map:
		return writeMap(buf, v, prefix, config)
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool,
		reflect.Uintptr, reflect.UnsafePointer:
		// 處理基本類型，輸出為鍵值對
		if prefix != "" {
			// fmt.Fprintf(buf, "%s%s: %v\n", config.KeyValuePrefix, prefix, formatValue(v))
			fmt.Fprintf(buf, "%s%s: %v\n", config.HeaderPrefix, prefix, formatValue(v))
		}
		return nil
	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return writeValue(buf, v.Elem(), prefix, config)
	default:
		return fmt.Errorf("unsupported kind %v at %s", v.Kind(), prefix)
	}
}

// formatValue 格式化基本類型的值
func formatValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		return v.String()
	case reflect.Uintptr:
		return fmt.Sprintf("0x%x", v.Uint()) // 格式化為十六進制
	case reflect.UnsafePointer:
		// 將 unsafe.Pointer 轉為 uintptr 並格式化為十六進制
		return fmt.Sprintf("0x%x", uintptr(v.UnsafePointer()))
	case reflect.Slice, reflect.Array:
		// 簡單切片直接格式化為字符串（避免遞迴）
		var values []string
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.Ptr && item.IsNil() {
				values = append(values, "nil")
				continue
			}
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			// 只處理基本類型，避免複雜結構
			if item.Kind() == reflect.Struct || item.Kind() == reflect.Slice || item.Kind() == reflect.Array || item.Kind() == reflect.Map || item.Kind() == reflect.Interface {
				values = append(values, fmt.Sprintf("%v", item.Interface()))
			} else {
				values = append(values, formatValue(item))
			}
		}
		return strings.Join(values, ", ")
	default:
		// fmt.Println("🔥", v.Kind())
		return fmt.Sprintf("%v", v.Interface())
	}
}

// writeStruct 處理 struct
func writeStruct(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	t := v.Type()
	hasWrittenHeader := false

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 獲取 JSON 標籤名稱
		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}
		if name == "-" {
			continue
		}

		newPrefix := name
		if prefix != "" {
			newPrefix = prefix + "." + name
		}

		// 處理非複雜類型的欄位（鍵值對）
		if value.Kind() != reflect.Struct &&
			value.Kind() != reflect.Slice &&
			value.Kind() != reflect.Array &&
			value.Kind() != reflect.Map &&
			value.Kind() != reflect.Interface &&
			value.Kind() != reflect.Pointer {
			if !hasWrittenHeader && prefix != "" {
				fmt.Fprintf(buf, "%s%s\n", config.HeaderPrefix, prefix)
				hasWrittenHeader = true
			}
			fmt.Fprintf(buf, "%s%s: %v\n", config.KeyValuePrefix, name, formatValue(value))
			continue
		}

		// 處理複雜類型（struct、slice、map、interface）
		if err := writeValue(buf, value, newPrefix, config); err != nil {
			return err
		}
	}

	if hasWrittenHeader {
		fmt.Fprintln(buf) // 鍵值對後添加空行
	}
	return nil
}

// writeSlice 處理陣列或切片
func writeSlice(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	if v.Len() == 0 {
		return nil
	}

	// 檢查元素類型
	elem := v.Index(0)
	// for elem.Kind() == reflect.Ptr {
	// 	if elem.IsNil() {
	// 		return nil // 跳過 nil 指針
	// 	}
	// 	elem = elem.Elem() // 解引用
	// }
	if elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			return nil // 跳過 nil 指針
		}
		elem = elem.Elem() // 解引用
	}

	if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
		// 處理嵌套切片
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.Ptr && item.IsNil() {
				continue // 跳過 nil 指針
			}
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			if err := writeSlice(buf, v.Index(i), newPrefix, config); err != nil {
				return err
			}
		}
		return nil
	}

	if elem.Kind() != reflect.Struct {
		// switch elem.Kind() {
		// case reflect.Bool:
		// 	fmt.Fprintf(buf, "%v\n", elem.Bool())
		// 	return nil
		// default:
		// 	return fmt.Errorf("slice elements must be structs, got %v at %s", elem.Kind(), prefix)
		// }
		return writeValue(buf, elem, prefix, config)
		// return fmt.Errorf("slice elements must be structs, got %v at %s", elem.Kind(), prefix)
	}

	// 獲取欄位名稱
	t := elem.Type()
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
		if item.Kind() == reflect.Ptr {
			if item.IsNil() {
				fmt.Fprintf(buf, "%s\n", strings.Repeat("-", len(headers))) // 輸出空行或標記
				continue
			}
			item = item.Elem()
		}
		var values []string
		for j := 0; j < item.Type().NumField(); j++ {
			if item.Type().Field(j).Tag.Get("json") == "-" {
				continue
			}
			values = append(values, formatValue(item.Field(j)))
		}
		fmt.Fprintf(buf, "%s\n", strings.Join(values, config.Delimiter))
	}

	fmt.Fprintln(buf) // 陣列後添加空行. 如果可以不需要 HeaderPrefix: "\n# "
	return nil
}

// writeMap 處理 map
// writeMap 處理 map
func writeMap(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	if v.Len() == 0 {
		return nil
	}

	// 收集並排序鍵
	type keyPair struct {
		value     reflect.Value // 原始鍵
		formatted string        // 格式化後的鍵（用於排序和顯示）
	}
	keys := v.MapKeys()
	sortedKeys := make([]keyPair, 0, len(keys))
	for _, key := range keys {
		sortedKeys = append(sortedKeys, keyPair{
			value:     key,
			formatted: formatValue(key),
		})
	}
	// 按格式化後的鍵排序
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i].formatted < sortedKeys[j].formatted
	})

	// 寫入上下文標記
	fmt.Fprintf(buf, "%s%s\n", config.HeaderPrefix, prefix)

	// 處理每個鍵值對
	for _, keyPair := range sortedKeys {
		key := keyPair.formatted         // 用於顯示
		val := v.MapIndex(keyPair.value) // 使用原始鍵查找值
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			val = val.Elem()
		}
		// 遞迴處理值，支援接口、結構體、切片等
		newPrefix := fmt.Sprintf("%s.%s", prefix, key)
		if err := writeValue(buf, val, newPrefix, config); err != nil {
			return err
		}
	}

	fmt.Fprintln(buf) // map 後添加空行
	return nil
}
