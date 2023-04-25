package bencode

import (
	"bytes"
	"fmt"
	"sort"
)

type Bencode interface {
	Encode() []byte
}

type String string
type Integer int64
type List []Bencode
type Dict map[String]Bencode

var (
	_ Bencode = String("")
	_ Bencode = Integer(0)
	_ Bencode = List{}
	_ Bencode = Dict{}
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
