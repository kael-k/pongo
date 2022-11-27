package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

var testTypeIntCases = []testSchemaCase{
	{
		desc:   "type-int-ok-1",
		schema: pongo.Int(),
		data:   func() pongo.Data { return 1 },
		want:   func() pongo.Data { return 1 },
		errors: 0,
	},
	{
		desc:   "type-int-ok-2",
		schema: pongo.Int().SetMin(1).SetMax(1),
		data:   func() pongo.Data { return 1 },
		want:   func() pongo.Data { return 1 },
		errors: 0,
	},
	{
		desc:   "type-int-ok-3",
		schema: pongo.Int().SetCast(true).SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 2 },
		want:   func() pongo.Data { return 2 },
		errors: 0,
	},
	{
		desc:   "type-int-ok-4",
		schema: pongo.Int().SetCast(true).SetMin(1).SetMax(3),
		data:   func() pongo.Data { return "2" },
		want:   func() pongo.Data { return 2 },
		errors: 0,
	},
	{
		desc:   "type-int-ok-5",
		schema: pongo.Int().SetCast(true).SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 2.0 },
		want:   func() pongo.Data { return 2 },
		errors: 0,
	},
	{
		desc:   "type-int-ok-6",
		schema: pongo.Int().SetCast(true).SetMin(1).SetMax(3),
		data:   func() pongo.Data { return float32(2.0) },
		want:   func() pongo.Data { return 2 },
		errors: 0,
	},
	{
		desc:   "type-int-ko-1",
		schema: pongo.Int().SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 0 },
		want:   func() pongo.Data { return 0 },
		errors: 1,
	},
	{
		desc:   "type-int-ko-2",
		schema: pongo.Int().SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 4 },
		want:   func() pongo.Data { return 4 },
		errors: 1,
	},
}

var testIntTypeSerializeCases = []testSchemaCase{
	{
		desc:   "int-serialize-ok-1",
		schema: pongo.Int().SetMin(10),
		data:   func() pongo.Data { return 12 },
		want:   func() pongo.Data { return 12 },
		errors: 0,
	},
	{
		desc:   "int-serialize-ko-1",
		schema: pongo.Int().SetMin(10),
		data:   func() pongo.Data { return 1 },
		want:   func() pongo.Data { return 1 },
		errors: 1,
	},
}

func TestTypeInt_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeIntCases)(t)
}

func TestTypeInt_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testIntTypeSerializeCases)(t)
}
