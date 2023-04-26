package bencode

import (
	"bytes"
	"fmt"
	"sort"
)

func Encode(v Value) []byte {
	return v.Encode()
}

type Value interface {
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
	return []byte(fmt.Sprintf("%d:%s", len(s), s))
}

func (i Integer) Encode() []byte {
	return []byte(fmt.Sprintf("i%de", i))
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

func (d Dict) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('d')
	var keys sort.StringSlice
	for key := range d {
		keys = append(keys, string(key))
	}
	keys.Sort()
	for _, k := range keys {
		key := String(k)
		value := d[key]
		buf.Write(key.Encode())
		buf.Write(value.Encode())
	}
	buf.WriteRune('e')
	return buf.Bytes()
}
