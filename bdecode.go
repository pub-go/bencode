package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	StrLenEnd = ':' // 字节串 ${len}:${str} // len 是字节长度
	IntStart  = 'i' // 整数 i${value}e. 不能有前导 0.  -0 非法.
	ListStart = 'l' // 列表 l[${item}]*e
	DictStart = 'd' // 字典 d[${strKey}${value}]*e
	End       = 'e' // 整数, 列表和字典的结尾标记
)

// Decode 从字节数组中解析出表示的值
func Decode(input []byte) (Value, error) {
	br := bufio.NewReader(bytes.NewReader(input))
	value, err := readValue(br)
	if err != nil {
		return nil, err
	}

	_, err = br.ReadByte()
	if !errors.Is(err, io.EOF) {
		br.UnreadByte()
		return nil, fmt.Errorf("input too long")
	}
	return value, nil
}

func readValue(br *bufio.Reader) (Value, error) {
	var b, err = br.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read value fail: %w", err)
	}
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		br.UnreadByte()                // 把已经读取的第一位数字还回去
		return wrapNil(readString(br)) // 读取 ${strLength}:${strData}
	case IntStart: // 已经把 i 读取
		return wrapNil(readIntUtil(br, End)) // 读取数字, 读到 'e' 为止
	case ListStart: // 已经把 l 读取
		return wrapNil(readList(br)) // 读取每个 item, 读到 'e' 为止
	case DictStart: // 已经把 d 读取
		return wrapNil(readDict(br)) // 读取每个 键值对, 读到 'e' 为止
	default:
		return nil, fmt.Errorf("unexpected byte: '%c'(%v)", b, b)
	}
}

func wrapNil(v Value, err error) (Value, error) {
	if err != nil {
		return nil, err
	}
	return v, nil
}

func readString(br *bufio.Reader) (String, error) {
	len, err := readIntUtil(br, StrLenEnd)
	if err != nil {
		return "", fmt.Errorf("read string length fail: %w", err)
	}
	if len < 0 {
		return "", fmt.Errorf("unexpected string length: %v", len)
	}
	data := make([]byte, len) // make slice length=len
	_, err = io.ReadFull(br, data)
	if err != nil {
		return "", fmt.Errorf("read string data fail(expected length=%d): %w", len, err)
	}
	return String(data), nil
}

func readIntUtil(br *bufio.Reader, util byte) (Integer, error) {
	bs, err := br.ReadBytes(util) // 一直读取直到遇到指定字符，返回数据包含指定的那个字符
	if err != nil {
		return 0, fmt.Errorf("read integer fail, expected ends with '%c'(%v): %w", util, util, err)
	}
	str := string(bs[:len(bs)-1]) // 包含分隔符
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("integer expected, got %q: %w", str, err)
	}
	// 校验有效性
	if i == 0 {
		if str[0] == '-' { // 0 不能有前导负号
			return 0, fmt.Errorf("-0 is invalid")
		}
		if len(str) != 1 { // 00 非法 不能有多余的前导0
			return 0, fmt.Errorf("unexpected lead '0': %v", str)
		}
	} else {
		lead := str[0]
		if i < 0 {
			lead = str[1]
		}
		if lead == '0' {
			return 0, fmt.Errorf("unexpected lead '0': %v", str)
		}
	}
	return Integer(i), nil
}

func readList(br *bufio.Reader) (List, error) {
	var i int
	var list = List{}
	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read list fail(%d item read): %w", i, err)
		}
		if b == End {
			break
		}
		br.UnreadByte() // 不是 'e' 还没结束, 还回去
		item, err := readValue(br)
		if err != nil {
			return nil, fmt.Errorf("read list fail(%d item read): %w", i, err)
		}
		list = append(list, item)
		i++
	}
	return list, nil
}

func readDict(br *bufio.Reader) (Dict, error) {
	var i int
	var dict = Dict{}
	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read dict fail(%d key-value read): %w", i, err)
		}
		if b == End {
			break
		}
		br.UnreadByte() // 不是 'e' 结束, 还回去
		key, err := readString(br)
		if err != nil {
			return nil, fmt.Errorf("read dict key fail(%d key-value read): %w", i, err)
		}
		value, err := readValue(br)
		if err != nil {
			return nil, fmt.Errorf("read dict value fail(%d key-value read): %w", i, err)
		}
		if _, ok := dict[key]; ok {
			return nil, fmt.Errorf("read dict fail, duplicate key %v", key)
		}
		dict[key] = value
		i++
	}
	return dict, nil
}
