package bencode

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Encode(v Value) []byte {
	return v.Encode()
}

type Value interface {
	fmt.Stringer // for print purpose
	Encode() []byte
}

type String string
type Integer int64
type List []Value
type Dict map[String]Value

var (
	_ Value = String("")
	_ Value = Integer(0)
	_ Value = List{}
	_ Value = Dict{}
)

func (s String) Encode() []byte {
	data := []byte(s) // 不直接对 String 使用 %s 输出，否则会调用到 String() 方法
	// 这里需要直接输出内容，不需要 String() 里的转义逻辑
	return []byte(fmt.Sprintf("%d:%s", len(data), data))
}

func (s String) String() string {
	return strconv.Quote(string(s))
}

func (i Integer) Encode() []byte {
	return []byte(fmt.Sprintf("i%de", i))
}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (l List) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('l')
	for _, item := range l {
		buf.Write(item.Encode())
	}
	buf.WriteRune('e')
	return buf.Bytes()
}

func (l List) String() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, item := range l {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(item.String())
	}
	sb.WriteRune(']')
	return sb.String()
}

func (d Dict) Range(fn func(index int, key String, value Value)) {
	var keys sort.StringSlice
	for key := range d {
		keys = append(keys, string(key))
	}
	keys.Sort()
	for i, k := range keys {
		key := String(k)
		value := d[key]
		fn(i, key, value)
	}
}

func (d Dict) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('d')
	d.Range(func(_ int, key String, value Value) {
		buf.Write(key.Encode())
		buf.Write(value.Encode())
	})
	buf.WriteRune('e')
	return buf.Bytes()
}

func (d Dict) String() string {
	var sb strings.Builder
	sb.WriteRune('{')
	d.Range(func(index int, key String, value Value) {
		if index > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(key.String())
		sb.WriteRune(':')
		sb.WriteString(value.String())
	})
	sb.WriteRune('}')
	return sb.String()
}
