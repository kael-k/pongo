package tests

import (
	"testing"

	"github.com/kael-k/pongo/v2/pongo"
)

var testTypeBytesCases = []testSchemaCase{
	{
		desc:   "type-bytes-ok-1",
		schema: pongo.Bytes(),
		data:   func() pongo.Data { return []byte("abc") },
		want:   func() pongo.Data { return []byte("abc") },
		errors: 0,
	},
	{
		desc:   "type-bytes-ok-2",
		schema: pongo.Bytes().SetCast(true),
		data:   func() pongo.Data { return "YWJj" }, // base64 for "abc"
		want:   func() pongo.Data { return []byte("abc") },
		errors: 0,
	},
	{
		desc:   "type-bytes-ok-3",
		schema: pongo.Bytes().SetCast(true),
		data:   func() pongo.Data { return []byte("abc") },
		want:   func() pongo.Data { return []byte("abc") },
		errors: 0,
	},
	{
		desc:   "type-bytes-ok-4",
		schema: pongo.Bytes().SetCast(true),
		data:   func() pongo.Data { return byte(31) },
		want:   func() pongo.Data { return []byte{31} },
		errors: 0,
	},
	{
		desc:   "type-bytes-ko-1",
		schema: pongo.Bytes(),
		data:   func() pongo.Data { return "YWJj" }, // base64 for "abc"
		want:   func() pongo.Data { return "YWJj" },
		errors: 1,
	},
	{
		desc:   "type-bytes-ko-2",
		schema: pongo.Bytes().SetMinLen(1).SetMaxLen(3),
		data:   func() pongo.Data { return []byte{} },
		want:   func() pongo.Data { return []byte{} },
		errors: 1,
	},
	{
		desc:   "type-bytes-ko-3",
		schema: pongo.Bytes().SetMinLen(1).SetMaxLen(3),
		data:   func() pongo.Data { return []byte{} },
		want:   func() pongo.Data { return []byte{} },
		errors: 1,
	},
}

var testBytesTypeSerializeCases = []testSchemaCase{
	{
		desc:   "byte-serialize-ok-1",
		schema: pongo.Bytes(),
		data:   func() pongo.Data { return []byte("abc") },
		want:   func() pongo.Data { return "YWJj" },
		errors: 0,
	},
	{
		desc:   "byte-serialize-ok-2",
		schema: pongo.Bytes().SetCast(true).SetMinLen(2),
		data:   func() pongo.Data { return "YWJj" },
		want:   func() pongo.Data { return "YWJj" },
		errors: 0,
	},
	{
		desc:   "byte-serialize-ko-1",
		schema: pongo.Bytes().SetMinLen(2),
		data:   func() pongo.Data { return "YWJj" },
		want:   func() pongo.Data { return "YWJj" },
		errors: 1,
	},
	{
		desc:   "byte-serialize-ko-2",
		schema: pongo.Bytes().SetCast(true).SetMinLen(4),
		data:   func() pongo.Data { return "YWJj" },
		want:   func() pongo.Data { return "YWJj" },
		errors: 1,
	},
}

func TestBytesType_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeBytesCases)(t)
}

func TestBytesType_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testBytesTypeSerializeCases)(t)
}
