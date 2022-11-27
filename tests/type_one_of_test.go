package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

var testTypeOneOfCases = []testSchemaCase{
	{
		desc: "type-one-of-1-ok",
		schema: pongo.OneOf(
			pongo.String(),
			pongo.Int(),
		),
		data:   func() pongo.Data { return 123 },
		want:   func() pongo.Data { return 123 },
		errors: 0,
	}, {
		desc: "type-one-of-2-ok",
		schema: pongo.OneOf(
			pongo.String().SetMinLen(4),
			pongo.String().SetMinLen(1),
			pongo.String().SetMaxLen(2),
		),
		data:   func() pongo.Data { return "aaa" },
		want:   func() pongo.Data { return "aaa" },
		errors: 0,
	}, {
		desc: "type-one-of-1-ko",
		schema: pongo.OneOf(
			pongo.String().SetMinLen(3),
			pongo.String().SetMinLen(1),
			pongo.String().SetMaxLen(2),
		),
		data:   func() pongo.Data { return "a" },
		want:   func() pongo.Data { return "a" },
		errors: 1,
	},
}

var testOneOfTypeSerializeCases = []testSchemaCase{
	{
		desc: "one-of-serialize-ok-1",
		schema: pongo.OneOf(
			pongo.String().SetCast(true),
			pongo.Bool(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return "123456" },
		errors: 0,
	},
	{
		desc: "one-of-serialize-ok-2",
		schema: pongo.OneOf(
			pongo.String(),
			pongo.Int(),
			pongo.Bool(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return 123456 },
		errors: 0,
	},
	{
		desc: "one-of-serialize-ko-1",
		schema: pongo.OneOf(
			pongo.String().SetCast(true),
			pongo.Bool(),
			pongo.Int(),
		),
		data:   func() pongo.Data { return 123456 },
		want:   func() pongo.Data { return 123456 },
		errors: 1,
	},
}

func TestTypeOneOf_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeOneOfCases)(t)
}

func TestTypeOneOf_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testOneOfTypeSerializeCases)(t)
}
