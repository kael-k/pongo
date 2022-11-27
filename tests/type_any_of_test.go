package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

var testTypeAnyOfCases = []testSchemaCase{
	{
		desc: "type-any-of-1-ok",
		schema: pongo.AnyOf(
			pongo.String(),
			pongo.String().SetMaxLen(10),
			pongo.String().SetMinLen(1),
		),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 0,
	},
	{
		desc: "type-any-of-2-ok",
		schema: pongo.AnyOf(
			pongo.String(),
			pongo.String().SetMinLen(10),
			pongo.String().SetMinLen(9),
			pongo.String().SetMinLen(1),
		),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 0,
	}, {
		desc: "type-any-of-1-ko",
		schema: pongo.AnyOf(
			pongo.String().SetMinLen(10),
			pongo.String().SetMinLen(9),
		),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 2,
	},
}

var testAnyOfTypeSerializeCases = []testSchemaCase{
	{
		desc: "any-of-serialize-ok-1",
		schema: pongo.AnyOf(
			pongo.String().SetCast(true),
			pongo.Bool(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc: "any-of-serialize-ok-2",
		schema: pongo.AnyOf(
			pongo.String(),
			pongo.Int(),
			pongo.Bool(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return 123456 },
		errors: 0,
	},
	{
		desc: "any-of-serialize-ok-3",
		schema: pongo.AnyOf(
			pongo.String().SetCast(true),
			pongo.Bool(),
			pongo.Int(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc: "any-of-serialize-ko-1",
		schema: pongo.AnyOf(
			pongo.String(),
			pongo.Bool(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return 123456 },
		errors: 2,
	},
}

func TestTypeAnyOf_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeAnyOfCases)(t)
}

func TestTypeAnyOf_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testAnyOfTypeSerializeCases)(t)
}
