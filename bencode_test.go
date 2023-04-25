package bencode_test

import (
	"testing"

	"code.gopub.tech/bencode"
)

func TestBencode(t *testing.T) {
	t.Logf("lne(%s)=%d", "你好", len("你好")) // lne(你好)=6
	var tests = []struct {
		name string
		args bencode.Bencode
		want string
	}{
		{name: "string/empty", args: bencode.String(""), want: "0:"},
		{name: "string/ascii", args: bencode.String("spam"), want: "4:spam"},
		{name: "string/emoji", args: bencode.String("⭕️"), want: "6:⭕️"}, // [54 58 226 173 149 239 184 143]
		{name: "string/cn", args: bencode.String("你好"), want: "6:你好"},    // [54 58 228 189 160 229 165 189]
		{name: "int", args: bencode.Integer(0), want: "i0e"},
		{name: "int-0", args: bencode.Integer(-0), want: "i0e"},
		{name: "int1", args: bencode.Integer(1), want: "i1e"},
		{name: "int-1", args: bencode.Integer(-1), want: "i-1e"},
		{name: "list-empty", args: bencode.List{}, want: "le"},
		{name: "list", args: bencode.List{bencode.String("spam"), bencode.Integer(1)}, want: "l4:spami1ee"},
		{name: "dict", args: bencode.Dict{
			bencode.String("spam"): bencode.String("eggs"),
			bencode.String("cow"):  bencode.String("moo"),
		}, want: "d3:cow3:moo4:spam4:eggse"}, // key sorted
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Encode()
			if string(got) != tt.want {
				t.Errorf("Encode() = %s, want = %s", got, tt.want)
			}
		})
	}
}
