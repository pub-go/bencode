# bencode
a Go implementation of Bencode.

[bep_0003.rst](https://github.com/bittorrent/bittorrent.org/blob/master/beps/bep_0003.rst) | [bep_0003.md(中文)](https://github.com/pub-go/bep/bep_0003.md)

## Usage
```go
import "code.gopub.tech/bencode"

var str = bencode.String("spam")
var data = bencode.Encode(str)   // []byte(`4:spam`)
val, err := bencode.Decode(data) // val==str, err==nil

var i = bencode.Integer(42)
var list = bencode.List{
    bencode.String("item1"),
    bencode.Integer(2),
    bencode.List{},
    bencode.Dict{},
}
var dict = bencode.Dict{
    bencode.String("key-must-be-string"): bencode.String("value can be any bencode.Value"),
}
```

## Implement

### Encode
```go
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

```

### Decode
```go
func Decode(input []byte) (Value, error) {
    ...
}
```
