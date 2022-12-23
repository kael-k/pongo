package tests

import (
	"testing"

	"github.com/kael-k/pongo/v2/pongo"
)

var testTypeStringCases = []testSchemaCase{
	{
		desc:   "type-string-ok-1",
		schema: pongo.String(),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-2",
		schema: pongo.String().SetMinLen(2).SetMaxLen(5),
		data:   func() pongo.Data { return "aaaaa" },
		want:   func() pongo.Data { return "aaaaa" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-3",
		schema: pongo.String().SetMinLen(2).SetMaxLen(5),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-4",
		schema: pongo.String().SetCast(true),
		data:   func() pongo.Data { return 12345 },
		want:   func() pongo.Data { return "12345" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-5",
		schema: pongo.String(),
		data:   func() pongo.Data { return "" },
		want:   func() pongo.Data { return "" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-6",
		schema: pongo.String().SetMinLen(0),
		data:   func() pongo.Data { return "" },
		want:   func() pongo.Data { return "" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-7",
		schema: pongo.String().SetMaxLen(0),
		data:   func() pongo.Data { return "" },
		want:   func() pongo.Data { return "" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-8",
		schema: pongo.String().SetCast(true),
		data:   func() pongo.Data { return int64(42) },
		want:   func() pongo.Data { return "42" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-9",
		schema: pongo.String().SetCast(true),
		data:   func() pongo.Data { return float32(4.2) },
		want:   func() pongo.Data { return "4.2" },
		errors: 0,
	},
	{
		desc:   "type-string-ok-10",
		schema: pongo.String().SetCastActions(pongo.SchemaActionParse),
		data:   func() pongo.Data { return 4.2 },
		want:   func() pongo.Data { return "4.2" },
		errors: 0,
	},
	{
		desc:   "type-string-ko-1",
		schema: pongo.String().SetMinLen(1),
		data:   func() pongo.Data { return "" },
		want:   func() pongo.Data { return "" },
		errors: 1,
	},
	{
		desc:   "type-string-ko-2",
		schema: pongo.String().SetMinLen(2).SetMaxLen(5),
		data:   func() pongo.Data { return "a" },
		want:   func() pongo.Data { return "a" },
		errors: 1,
	},
	{
		desc:   "type-string-ko-3",
		schema: pongo.String().SetMinLen(2).SetMaxLen(5),
		data:   func() pongo.Data { return "aaaaaaaaaaaaaaa" },
		want:   func() pongo.Data { return "aaaaaaaaaaaaaaa" },
		errors: 1,
	},
	{
		desc:   "type-string-ko-4",
		schema: pongo.String().SetCastActions(pongo.SchemaActionSerialize),
		data:   func() pongo.Data { return 4.2 },
		want:   func() pongo.Data { return "4.2" },
		errors: 1,
	},
}

var testStringTypeSerializeCases = []testSchemaCase{
	{
		desc:   "string-serialize-ok-1",
		schema: pongo.String().SetCast(true),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc:   "string-serialize-ok-2",
		schema: pongo.String(),
		data:   func() pongo.Data { return "1234abc" },
		want:   func() pongo.Data { return "1234abc" },
		errors: 0,
	},
	{
		desc:   "string-serialize-ko-1",
		schema: pongo.String(),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return 123456 },
		errors: 1,
	},
}

func TestTypeString_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeStringCases)(t)
}

func TestTypeString_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testStringTypeSerializeCases)(t)
}
