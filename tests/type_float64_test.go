package tests

import (
	"testing"

	"github.com/kael-k/pongo/v2/pongo"
)

var testFloat64TypeCases = []testSchemaCase{
	{
		desc:   "type-float64-ok-1",
		schema: pongo.Float64(),
		data:   func() pongo.Data { return 1.1 },
		want:   func() pongo.Data { return 1.1 },
		errors: 0,
	},
	{
		desc:   "type-float64-ok-2",
		schema: pongo.Float64().SetCast(true),
		data:   func() pongo.Data { return 1 },
		want:   func() pongo.Data { return float64(1) },
		errors: 0,
	},
	{
		desc:   "type-float64-ok-3",
		schema: pongo.Float64().SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 2.22 },
		want:   func() pongo.Data { return 2.22 },
		errors: 0,
	},
	{
		desc:   "type-float64-ok-4",
		schema: pongo.Float64().SetCast(true).SetMax(5),
		data:   func() pongo.Data { return float32(4) },
		want:   func() pongo.Data { return float64(4) },
		errors: 0,
	},
	{
		desc:   "type-float64-ok-5",
		schema: pongo.Float64().SetCast(true).SetMax(5),
		data:   func() pongo.Data { return int64(4) },
		want:   func() pongo.Data { return float64(4) },
		errors: 0,
	},
	{
		desc:   "type-float64-ok-5",
		schema: pongo.Float64().SetCast(true).SetMax(5),
		data:   func() pongo.Data { return "4.2" },
		want:   func() pongo.Data { return 4.2 },
		errors: 0,
	},
	{
		desc:   "type-float64-ko-1",
		schema: pongo.Float64().SetCast(true).SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 0.1 },
		want:   func() pongo.Data { return 0.1 },
		errors: 1,
	},
	{
		desc:   "type-float64-ko-2",
		schema: pongo.Float64().SetMin(1).SetMax(3),
		data:   func() pongo.Data { return 4.1 },
		want:   func() pongo.Data { return 4.1 },
		errors: 1,
	},
	{
		desc:   "type-float64-ko-3",
		schema: pongo.Float64().SetCast(true).SetMax(5),
		data:   func() pongo.Data { return "not-a-number" },
		want:   func() pongo.Data { return 4.2 },
		errors: 1,
	},
}

var testFloat64TypeSerializeCases = []testSchemaCase{
	{
		desc:   "float64-serialize-ok-1",
		schema: pongo.Float64().SetMin(10),
		data:   func() pongo.Data { return 12.1 },
		want:   func() pongo.Data { return 12.1 },
		errors: 0,
	},
	{
		desc:   "float64-serialize-ko-1",
		schema: pongo.Float64().SetMin(10),
		data:   func() pongo.Data { return 1 },
		want:   func() pongo.Data { return 1 },
		errors: 1,
	},
}

func TestFloat64Type_Parse(t *testing.T) {
	testSchemaCaseParse(testFloat64TypeCases)(t)
}

func TestFloat64Type_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testFloat64TypeSerializeCases)(t)
}
