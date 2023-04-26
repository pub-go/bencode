package bencode

import "unicode/utf8"

func IsUTF8(s String) bool {
	return utf8.ValidString(string(s))
}
