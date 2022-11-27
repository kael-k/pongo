package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

var testTypeBoolCases = []testSchemaCase{
	{
		desc:   "type-bool-ok-1",
		schema: pongo.Bool(),
		data:   func() pongo.Data { return true },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-2",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return "true" },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-3",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return "false" },
		want:   func() pongo.Data { return false },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-4",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return "False" },
		want:   func() pongo.Data { return false },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-5",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return 0 },
		want:   func() pongo.Data { return false },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-6",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return 1.3 },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-7",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return byte(0) },
		want:   func() pongo.Data { return false },
		errors: 0,
	}, {
		desc:   "type-bool-ok-8",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return byte(1) },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-9",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return []byte{0} },
		want:   func() pongo.Data { return false },
		errors: 0,
	}, {
		desc:   "type-bool-ok-10",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return []byte{1, 0} },
		want:   func() pongo.Data { return true },
		errors: 0,
	}, {
		desc:   "type-bool-ok-11",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return []byte{1} },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-12",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return []byte{1} },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "type-bool-ok-13",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return []byte{} },
		want:   func() pongo.Data { return false },
		errors: 0,
	},
	{
		desc:   "type-int-ko-1",
		schema: pongo.Bool(),
		data:   func() pongo.Data { return 4 },
		want:   func() pongo.Data { return 4 },
		errors: 1,
	},
	{
		desc:   "type-int-ko-1",
		schema: pongo.Bool().SetCast(true),
		data:   func() pongo.Data { return "foo" },
		want:   func() pongo.Data { return false },
		errors: 1,
	},
}

var testBoolTypeSerializeCases = []testSchemaCase{
	{
		desc:   "bool-serialize-ok-1",
		schema: pongo.Bool(),
		data:   func() pongo.Data { return true },
		want:   func() pongo.Data { return true },
		errors: 0,
	},
	{
		desc:   "bool-serialize-ok-2",
		schema: pongo.Bool(),
		data:   func() pongo.Data { return false },
		want:   func() pongo.Data { return false },
		errors: 0,
	},
	{
		desc:   "bool-serialize-ko-2",
		schema: pongo.Bool(),
		data:   func() pongo.Data { return 1 },
		want:   func() pongo.Data { return 1 },
		errors: 1,
	},
}

func TestTypeBool_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeBoolCases)(t)
}

func TestTypeBool_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testBoolTypeSerializeCases)(t)
}
