package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

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
	case 'i': // 已经把 i 读取
		return wrapNil(readIntUtil(br, 'e')) // 读取数字, 读到 'e' 为止
	case 'l': // 已经把 l 读取
		return wrapNil(readList(br)) // 读取每个 item, 读到 'e' 为止
	case 'd': // 已经把 d 读取
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
	len, err := readIntUtil(br, ':')
	if err != nil {
		return "", fmt.Errorf("read string length fail: %w", err)
	}
	if len < 0 {
		return "", fmt.Errorf("unexpected string length: %v", len)
	}
	data := make([]byte, len) // make slice length=len
	_, err = io.ReadFull(br, data)
	if err != nil {
		return "", fmt.Errorf("read string data fail(length=%d): %w", len, err)
	}
	return String(data), nil
}

func readIntUtil(br *bufio.Reader, util byte) (Integer, error) {
	bs, err := br.ReadBytes(util) // 一直读取知道遇到指定字符，返回数据包含指定的那个字符
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
			return 0, fmt.Errorf("unexpected lead '0'")
		}
	} else {
		lead := str[0]
		if i < 0 {
			lead = str[1]
		}
		if lead == '0' {
			return 0, fmt.Errorf("unexpected lead '0'")
		}
	}
	return Integer(i), nil
}

func readList(br *bufio.Reader) (List, error) {
	var list = List{}
	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read list fail: %w", err)
		}
		if b == 'e' {
			break
		}
		br.UnreadByte() // 不是 'e' 还没结束, 还回去
		item, err := readValue(br)
		if err != nil {
			return nil, fmt.Errorf("read list item fail: %w", err)
		}
		list = append(list, item)
	}
	return list, nil
}

func readDict(br *bufio.Reader) (Dict, error) {
	var dict = Dict{}
	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read dict fail: %w", err)
		}
		if b == 'e' {
			break
		}
		br.UnreadByte() // 不是 'e' 结束, 还回去
		key, err := readString(br)
		if err != nil {
			return nil, fmt.Errorf("read dict key fail: %w", err)
		}
		value, err := readValue(br)
		if err != nil {
			return nil, fmt.Errorf("read dict value fail: %w", err)
		}
		dict[key] = value
	}
	return dict, nil
}
