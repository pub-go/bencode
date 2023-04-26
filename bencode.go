package bencode

import (
	"bytes"
	"encoding/base64"
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
	var str = string(s)
	if !IsUTF8(s) {
		str = base64.StdEncoding.EncodeToString([]byte(s))
		// str = fmt.Sprintf("b'%s'", str) // is this needed? 为了打印时区分
		// 删除上行注释，会这样打印，打印时可以区分出 []byte 和 base64 的 string
		// []byte{111, 222}  / Encode("b'b94='") = "2:o\xde"
		// string("b94=")    / Encode("b94=") = "4:b94="
		// 但多加一层也没啥区分作用
		// string("b'b94='") / Encode("b'b94='") = "7:b'b94='"

		// 不使用 b'' 包裹时，会这样打印，就这样吧
		// []byte{111, 222}  / Encode("b94=") = "2:o\xde"
		// string("b94=")    / Encode("b94=") = "4:b94="
		// string("b'b94='") / Encode("b'b94='") = "7:b'b94='"
	}
	return strconv.Quote(str)
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
