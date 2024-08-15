package bencode_test

import (
	"testing"

	"code.gopub.tech/assert"
	"code.gopub.tech/bencode"
)

func TestAs(t *testing.T) {
	var val bencode.Value = bencode.String("str")
	assert.Equal(t, bencode.AsStr(val), "str")
	assert.Equal(t, bencode.AsInt(val), int64(0))
	assert.DeepEqual(t, bencode.AsList(val), bencode.List{})
	assert.DeepEqual(t, bencode.AsDict(val), bencode.Dict{})

	val = bencode.Integer(10)
	assert.Equal(t, bencode.AsStr(val), "")
	assert.Equal(t, bencode.AsInt(val), int64(10))
	assert.DeepEqual(t, bencode.AsList(val), bencode.List{})
	assert.DeepEqual(t, bencode.AsDict(val), bencode.Dict{})

	val = bencode.Dict{"key": bencode.Integer(10)}
	assert.Equal(t, bencode.AsStr(val), "")
	assert.Equal(t, bencode.AsInt(val), int64(0))
	assert.DeepEqual(t, bencode.AsList(val), bencode.List{})
	assert.DeepEqual(t, bencode.AsDict(val), bencode.Dict{"key": bencode.Integer(10)})

	val = bencode.List{bencode.Integer(10), bencode.String("str")}
	assert.Equal(t, bencode.AsStr(val), "")
	assert.Equal(t, bencode.AsInt(val), int64(0))
	assert.DeepEqual(t, bencode.AsList(val), bencode.List{bencode.Integer(10), bencode.String("str")})
	assert.DeepEqual(t, bencode.AsDict(val), bencode.Dict{})
}
