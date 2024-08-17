package bencode

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Value
		wantErr bool
	}{
		{name: "nil", args: args{[]byte(nil)}, want: nil, wantErr: true},
		{name: "empty", args: args{[]byte(``)}, want: nil, wantErr: true},
		{name: "invalid", args: args{[]byte(`a`)}, want: nil, wantErr: true},
		{name: "str-invalid", args: args{[]byte(`1`)}, want: nil, wantErr: true},
		{name: "str-invalid-len", args: args{[]byte(`1:`)}, want: nil, wantErr: true},
		{name: "str-empty", args: args{[]byte(`0:`)}, want: String(""), wantErr: false},
		{name: "str-ok", args: args{[]byte(`4:spam`)}, want: String("spam"), wantErr: false},
		{name: "int-invalid", args: args{[]byte(`i`)}, want: nil, wantErr: true},
		{name: "int-invalid-empty", args: args{[]byte(`ie`)}, want: nil, wantErr: true},
		{name: "int-invalid-0", args: args{[]byte(`i-0e`)}, want: nil, wantErr: true},
		{name: "int-invalid-00", args: args{[]byte(`i00e`)}, want: nil, wantErr: true},
		{name: "int-invalid-lead", args: args{[]byte(`i01e`)}, want: nil, wantErr: true},
		{name: "int-invalid-noend", args: args{[]byte(`i0`)}, want: nil, wantErr: true},
		{name: "int-ok", args: args{[]byte(`i0e`)}, want: Integer(0), wantErr: false},
		{name: "int-ok-neg", args: args{[]byte(`i-1e`)}, want: Integer(-1), wantErr: false},
		{name: "int-ok-pos", args: args{[]byte(`i1e`)}, want: Integer(1), wantErr: false},
		{name: "list-invalid", args: args{[]byte(`l`)}, want: nil, wantErr: true},
		{name: "list-empty", args: args{[]byte(`le`)}, want: List{}, wantErr: false},
		{name: "list-ok", args: args{[]byte(`li0ee`)}, want: List{Integer(0)}, wantErr: false},
		{name: "list-item-invalid", args: args{[]byte(`li0`)}, want: nil, wantErr: true},
		{name: "dict-invalid", args: args{[]byte(`d`)}, want: nil, wantErr: true},
		{name: "dic-empty", args: args{[]byte(`de`)}, want: Dict{}, wantErr: false},
		{name: "dict-key-invalid", args: args{[]byte(`di`)}, want: nil, wantErr: true},
		{name: "dict-key-invalid2", args: args{[]byte(`d:`)}, want: nil, wantErr: true},
		{name: "dict-value-invalid", args: args{[]byte(`d1:a`)}, want: nil, wantErr: true},
		{name: "dict-value-invalid2", args: args{[]byte(`d1:a1`)}, want: nil, wantErr: true},
		{name: "dict-dup-key", args: args{[]byte(`d1:a1:a1:a1:be`)}, want: nil, wantErr: true},
		{name: "dict-ok", args: args{[]byte(`d1:a1:ae`)}, want: Dict{String("a"): String("a")}, wantErr: false},
		{name: "too-long", args: args{[]byte(`1:abc`)}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.input)
			t.Logf("name = %v, input = %s, got = %v , err=%+v", tt.name, tt.args.input, got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
