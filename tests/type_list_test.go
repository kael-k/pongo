package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

var testTypeListCases = []testSchemaCase{
	{
		desc:   "type-list-1-ok",
		schema: pongo.List(pongo.String().SetMinLen(2).SetMaxLen(5)),
		data:   func() pongo.Data { return []any{"aaa", "AAAA"} },
		want:   func() pongo.Data { return []any{"aaa", "AAAA"} },
		errors: 0,
	},
	{
		desc: "type-list-2-ok",
		schema: pongo.List(
			pongo.AnyOf(
				pongo.String(),
				pongo.Int(),
			),
		),
		data:   func() pongo.Data { return []any{"zpdijf0sj", 12345, "test", "342", 11} },
		want:   func() pongo.Data { return []any{"zpdijf0sj", 12345, "test", "342", 11} },
		errors: 0,
	},
	{
		desc:   "type-list-1-ko",
		schema: pongo.List(pongo.String()),
		data:   func() pongo.Data { return []any{"zpdijf0sj", 12345, "test", "342", 11} },
		want:   func() pongo.Data { return []any{"zpdijf0sj", 12345, "test", "342", 11} },
		errors: 2,
	},
	{
		desc:   "type-list-1-ko",
		schema: pongo.List(pongo.String()),
		data:   func() pongo.Data { return nil },
		want:   func() pongo.Data { return nil },
		errors: 1,
	},
}

var testListTypeSerializeCases = []testSchemaCase{
	{
		desc:   "list-serialize-ok-1",
		schema: pongo.List(pongo.String().SetCast(true)),
		data:   func() pongo.Data { return []interface{}{123456, "abc"} },
		want:   func() pongo.Data { return []interface{}{"123456", "abc"} },
		errors: 0,
	},
	{
		desc:   "list-serialize-ok-2",
		schema: pongo.List(pongo.String().SetCast(true)).SetMinLen(1),
		data:   func() pongo.Data { return []interface{}{123456, "abc"} },
		want:   func() pongo.Data { return []interface{}{"123456", "abc"} },
		errors: 0,
	},
	{
		desc:   "list-serialize-ok-3",
		schema: pongo.List(pongo.String().SetCast(true)).SetMaxLen(4),
		data:   func() pongo.Data { return []interface{}{123456, "abc"} },
		want:   func() pongo.Data { return []interface{}{"123456", "abc"} },
		errors: 0,
	},
	{
		desc:   "list-serialize-ko-1",
		schema: pongo.List(pongo.String()),
		data:   func() pongo.Data { return []interface{}{123456, "abc"} },
		want:   func() pongo.Data { return []interface{}{123456, "abc"} },
		errors: 1,
	},
	{
		desc:   "list-serialize-ko-2",
		schema: pongo.List(pongo.String()).SetMinLen(10),
		data:   func() pongo.Data { return []interface{}{123456, "abc"} },
		want:   func() pongo.Data { return []interface{}{123456, "abc"} },
		errors: 1,
	},
	{
		desc:   "list-serialize-ko-3",
		schema: pongo.List(pongo.String()).SetMaxLen(1),
		data:   func() pongo.Data { return []interface{}{123456, "abc"} },
		want:   func() pongo.Data { return []interface{}{123456, "abc"} },
		errors: 1,
	},
	{
		desc:   "list-serialize-ko-4",
		schema: pongo.List(nil),
		data:   func() pongo.Data { return []interface{}{123456, 12, 1} },
		want:   func() pongo.Data { return []interface{}{123456, 12, 1} },
		errors: 1,
	},
}

func TestTypeList_Parse(t *testing.T) {
	testSchemaCaseProcess(testTypeListCases, pongo.SchemaActionParse)(t)
}

func TestTypeList_Serialize(t *testing.T) {
	testSchemaCaseProcess(testListTypeSerializeCases, pongo.SchemaActionSerialize)(t)
}
