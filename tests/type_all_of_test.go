package tests

import (
	"testing"

	"github.com/kael-k/pongo/v2/pongo"
)

var testTypeAllOfCases = []testSchemaCase{
	{
		desc: "type-all-of-1-ok",
		schema: pongo.AllOf(
			pongo.String(),
			pongo.String().SetMaxLen(10),
			pongo.String().SetMinLen(1),
		),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 0,
	},
	{
		desc: "type-all-of-2-ok",
		schema: pongo.AllOf(
			pongo.Int(),
			pongo.String().SetCast(true),
		),
		data:   func() pongo.Data { return 123 },
		want:   func() pongo.Data { return "123" },
		errors: 0,
	},
	{
		desc: "type-all-of-3-ok",
		schema: pongo.AllOf(
			pongo.Int(),
			pongo.String().SetCast(true),
			pongo.String().SetMinLen(2),
		).SetChain(true),
		data:   func() pongo.Data { return 123 },
		want:   func() pongo.Data { return "123" },
		errors: 0,
	},
	{
		desc: "type-all-of-1-ko",
		schema: pongo.AllOf(
			pongo.String(),
			pongo.String().SetMinLen(10),
			pongo.String().SetMinLen(9),
			pongo.String().SetMinLen(1),
		),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 2,
	},
	{
		desc: "type-all-of-2-ko",
		schema: pongo.AllOf(
			pongo.Int(),
			pongo.String().SetCast(true),
			pongo.String().SetMinLen(2),
		),
		data:   func() pongo.Data { return 123 },
		want:   func() pongo.Data { return "123" },
		errors: 1,
	},
}

var testAllOfTypeSerializeCases = []testSchemaCase{
	{
		desc: "all-of-serialize-ok-1",
		schema: pongo.AllOf(
			pongo.String().SetMinLen(1),
			pongo.String().SetMaxLen(9),
		),
		data:   func() pongo.Data { return "123456" },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc: "all-of-serialize-ok-2",
		schema: pongo.AllOf(
			pongo.String().SetCast(true),
			pongo.String().SetCast(true).SetMinLen(3),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc: "all-of-serialize-ok-3",
		schema: pongo.AllOf(
			pongo.Int(),
			pongo.String().SetCast(true),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc: "all-of-serialize-ko-1",
		schema: pongo.AllOf(
			pongo.String(),
			pongo.Bool(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return 123456 },
		errors: 2,
	},
	{
		desc: "all-of-serialize-ko-2",
		schema: pongo.AllOf(
			pongo.String().SetCast(true),
			pongo.String(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 1,
	},
}

func TestTypeAllOf_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeAllOfCases)(t)
}

func TestTypeAllOf_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testAllOfTypeSerializeCases)(t)
}
