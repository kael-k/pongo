package tests

import (
	"testing"
	"time"

	"github.com/kael-k/pongo/v2/pongo"
)

var testTypeDatetimeCases = []testSchemaCase{
	{
		desc:   "type-datatime-ok-1",
		schema: pongo.Datetime(),
		data:   func() pongo.Data { return time.Unix(1663770398, 0) },
		want:   func() pongo.Data { return time.Unix(1663770398, 0) },
		errors: 0,
	},
	{
		desc:   "type-datatime-ok-2",
		schema: pongo.Datetime().SetBefore(time.Unix(1663800000, 0)).SetAfter(time.Unix(1663770000, 0)),
		data:   func() pongo.Data { return time.Unix(1663770398, 0) },
		want:   func() pongo.Data { return time.Unix(1663770398, 0) },
		errors: 0,
	},
	{
		desc:   "type-datatime-ok-3",
		schema: pongo.Datetime().SetCast(true),
		data:   func() pongo.Data { return "2022-08-09T00:00:00.000000001Z" },
		want:   func() pongo.Data { return time.Unix(1660003200, 1).UTC() },
		errors: 0,
	},
	{
		desc:   "type-datatime-ok-4",
		schema: pongo.Datetime().SetCast(true),
		data:   func() pongo.Data { return 1660003200 },
		want:   func() pongo.Data { return time.Unix(1660003200, 0) },
		errors: 0,
	},
	{
		desc:   "type-datatime-ok-5",
		schema: pongo.Datetime().SetBefore(time.Unix(1663800000, 0)).SetCast(true).SetFormat(time.RFC822).SetAfter(time.Unix(1663770000, 0)),
		data:   func() pongo.Data { return time.Unix(1663800000, 0).UTC().Format(time.RFC822) },
		want:   func() pongo.Data { return time.Unix(1663800000, 0).UTC() },
		errors: 0,
	},
	{
		desc:   "type-datatime-ok-6",
		schema: pongo.Datetime().SetCastActions(pongo.SchemaActionParse),
		data:   func() pongo.Data { return time.Unix(1663770398, 0) },
		want:   func() pongo.Data { return time.Unix(1663770398, 0) },
		errors: 0,
	},
	{
		desc:   "type-datatime-ok-7",
		schema: pongo.Datetime().SetBefore(time.Unix(1663800000, 0)).SetCast(true).SetFormat(time.RFC822).SetAfter(time.Unix(1663770000, 0)),
		data:   func() pongo.Data { return int64(1663800000) },
		want:   func() pongo.Data { return time.Unix(1663800000, 0) },
		errors: 0,
	},
	{
		desc:   "type-datatime-ko-1",
		schema: pongo.Datetime().SetBefore(time.Unix(1663770000, 0)).SetAfter(time.Unix(1663800000, 0)),
		data:   func() pongo.Data { return time.Unix(1663700000, 0) },
		want:   func() pongo.Data { return time.Unix(1663700000, 0) },
		errors: 1,
	},
	{
		desc:   "type-datatime-ko-2",
		schema: pongo.Datetime().SetBefore(time.Unix(1663800000, 0)).SetAfter(time.Unix(1663770000, 0)),
		data:   func() pongo.Data { return time.Unix(1663900000, 0) },
		want:   func() pongo.Data { return time.Unix(1663900000, 0) },
		errors: 1,
	},
	{
		desc:   "type-datatime-ko-3",
		schema: pongo.Datetime().SetBefore(time.Unix(1663800000, 0)).SetCast(true).SetAfter(time.Unix(1663770000, 0)),
		data:   func() pongo.Data { return time.Unix(1663900000, 0) },
		want:   func() pongo.Data { return time.Unix(1663900000, 0) },
		errors: 1,
	},
}

var testDatetimeTypeSerializeCases = []testSchemaCase{
	{
		desc:   "datetime-serialize-ok-1",
		schema: pongo.Datetime(),
		data:   func() pongo.Data { return time.Unix(1660003200, 1).UTC() },
		want:   func() pongo.Data { return "2022-08-09T00:00:00.000000001Z" },
		errors: 0,
	},
	{
		desc:   "datetime-serialize-ok-2",
		schema: pongo.Datetime(),
		data:   func() pongo.Data { return time.Unix(1660003200, 0).UTC() },
		want:   func() pongo.Data { return "2022-08-09T00:00:00Z" },
		errors: 0,
	},
	{
		desc:   "datetime-serialize-ko-1",
		schema: pongo.Datetime(),
		data:   func() pongo.Data { return "not-a-datetime" },
		want:   func() pongo.Data { return "not-a-datetime" },
		errors: 1,
	},
}

func TestTypeDatetime_Parse(t *testing.T) {
	testSchemaCaseParse(testTypeDatetimeCases)(t)
}

func TestTypeDatetime_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testDatetimeTypeSerializeCases)(t)
}
