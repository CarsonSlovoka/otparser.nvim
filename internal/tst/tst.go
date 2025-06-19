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
	err := writeStruct(&buf, reflect.ValueOf(v), "", t.Config)
	if err != nil {
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
		reflect.Float32, reflect.Float64, reflect.Bool:
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
	case reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// Handle [4]byte or similar arrays
			bytes := make([]byte, v.Len())
			for i := 0; i < v.Len(); i++ {
				bytes[i] = byte(v.Index(i).Uint())
			}
			return string(bytes) // Convert to string (or use fmt.Sprintf("%x", bytes) for hex)
		}
		return fmt.Sprintf("%v", v.Interface())
	default:
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

		// Handle non-complex fields as key-value pairs
		if value.Kind() != reflect.Struct && value.Kind() != reflect.Slice && value.Kind() != reflect.Array && value.Kind() != reflect.Map && value.Kind() != reflect.Interface {
			if !hasWrittenHeader && prefix != "" {
				fmt.Fprintf(buf, "%s%s\n", config.HeaderPrefix, prefix)
				hasWrittenHeader = true
			}
			fmt.Fprintf(buf, "%s%s: %v\n", config.KeyValuePrefix, name, formatValue(value))
			continue
		}

		// 處理值
		if err := writeValue(buf, value, newPrefix, config); err != nil {
			return err
		}
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
	if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
		// 處理嵌套切片
		for i := 0; i < v.Len(); i++ {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			if err := writeSlice(buf, v.Index(i), newPrefix, config); err != nil {
				return err
			}
		}
		return nil
	}

	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be structs, got %v at %s", elem.Kind(), prefix)
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
		var values []string
		for j := 0; j < item.Type().NumField(); j++ {
			if item.Type().Field(j).Tag.Get("json") == "-" {
				continue
			}
			val := item.Field(j).Interface()
			switch v := val.(type) {
			case float32, float64:
				values = append(values, strconv.FormatFloat(float64(v.(float32)), 'f', 1, 64))
			default:
				values = append(values, fmt.Sprintf("%v", val))
			}
		}
		fmt.Fprintf(buf, "%s\n", strings.Join(values, config.Delimiter))
	}

	fmt.Fprintln(buf) // 陣列後添加空行
	return nil
}

// writeMap 處理 map
func writeMap(buf *bytes.Buffer, v reflect.Value, prefix string, config Config) error {
	if v.Len() == 0 {
		return nil
	}

	// 寫入上下文標記
	fmt.Fprintf(buf, "%s%s\n", config.HeaderPrefix, prefix)

	// 按鍵排序以確保一致輸出
	keys := v.MapKeys()
	sortedKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		if key.Kind() == reflect.String {
			sortedKeys = append(sortedKeys, key.String())
		}
	}
	sort.Strings(sortedKeys)

	// 寫入鍵值對
	for _, key := range sortedKeys {
		val := v.MapIndex(reflect.ValueOf(key))
		fmt.Fprintf(buf, "%s%s: %v\n", config.KeyValuePrefix, key, val.Interface())
	}

	fmt.Fprintln(buf) // map 後添加空行
	return nil
}
