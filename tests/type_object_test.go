package tests

import (
	"testing"
	"time"

	"github.com/kael-k/pongo/pongo"
)

var testObjectTypeCases = []testSchemaCase{
	{
		desc: "schema-1-ok",
		schema: pongo.Object(
			pongo.O{
				"test": pongo.String(),
				"nestedObject": pongo.Object(
					pongo.O{
						"test2": pongo.String(),
					},
				),
			},
		),
		data: func() pongo.Data {
			return map[string]any{
				"test": "value",
				"nestedObject": map[string]any{
					"test2": "value",
				},
			}
		},
		want: func() pongo.Data {
			return map[string]any{
				"test": "value",
				"nestedObject": map[string]any{
					"test2": "value",
				},
			}
		},
		errors: 0,
	},
	{
		desc: "schema-2-ok",
		schema: pongo.Object(
			pongo.O{
				"test": pongo.String(),
				"nestedObject": pongo.Object(
					pongo.O{
						"test2": pongo.String(),
					},
				),
			},
		).Require("test"),
		data: func() pongo.Data {
			return map[string]any{
				"test": "value",
				"nestedObject": map[string]any{
					"test2": "value",
				},
			}
		},
		want: func() pongo.Data {
			return map[string]any{
				"test": "value",
				"nestedObject": map[string]any{
					"test2": "value",
				},
			}
		},
		errors: 0,
	},
	{
		desc: "schema-3-ok",
		schema: pongo.Object(
			pongo.O{
				"test": pongo.String(),
				"nestedObject": pongo.Object(
					pongo.O{
						"test2": pongo.String(),
					},
				),
			},
		).Require("test"),
		data: func() pongo.Data {
			return map[string]any{
				"test": "value",
			}
		},
		want: func() pongo.Data {
			return map[string]any{
				"test": "value",
			}
		},
		errors: 0,
	},
	{
		desc: "schema-1-ko",
		schema: pongo.Object(pongo.O{
			"test": pongo.String(),
		}),
		data: func() pongo.Data {
			return map[string]any{
				"test": 12,
			}
		},
		want: func() pongo.Data {
			return map[string]any{
				"test": 12,
			}
		},
		errors: 1,
	},
	{
		desc: "schema-2-ok",
		schema: pongo.Object(
			pongo.O{
				"test": pongo.String(),
				"nestedObject": pongo.Object(
					pongo.O{
						"test2": pongo.String(),
					},
				),
			},
		).Require("nestedObject"),
		data: func() pongo.Data {
			return map[string]any{
				"test": "value",
			}
		},
		want: func() pongo.Data {
			return map[string]any{
				"test": "value",
			}
		},
		errors: 1,
	},
	{
		desc: "schema-3-ko",
		schema: pongo.Object(pongo.O{
			"string1": pongo.String().SetCast(true),
			"string2": pongo.String(),
		}),
		data: func() pongo.Data {
			return map[string]any{
				"string1": 12,
				"string2": 12,
			}
		},
		want: func() pongo.Data {
			return map[string]any{
				"string1": 12,
				"string2": 12,
			}
		},
		errors: 1,
	},
}

var testObjectTypeSerializeCases = []testSchemaCase{
	{
		desc: "object-serialize-ok-1",
		schema: pongo.Object(pongo.O{
			"aDatetime": pongo.Datetime(),
			"aString":   pongo.String().SetCast(true),
		}),
		data: func() pongo.Data {
			return map[string]interface{}{
				"aDatetime": time.Unix(1660003200, 0).UTC(),
				"aString":   12345,
			}
		},
		want: func() pongo.Data {
			return map[string]interface{}{
				"aDatetime": "2022-08-09T00:00:00Z",
				"aString":   "12345",
			}
		},
		errors: 0,
	},
	{
		desc: "object-serialize-ko-1",
		schema: pongo.Object(pongo.O{
			"aDatetime": pongo.Datetime(),
			"aString":   pongo.String(),
		}),
		data: func() pongo.Data {
			return map[string]interface{}{
				"aDatetime": time.Unix(1660003200, 0).UTC(),
				"aString":   12345,
			}
		},
		want: func() pongo.Data {
			return map[string]interface{}{
				"aDatetime": time.Unix(1660003200, 0).UTC(),
				"aString":   12345,
			}
		},
		errors: 1,
	},
}

func TestObjectType_Parse(t *testing.T) {
	testSchemaCaseParse(testObjectTypeCases)(t)
}

func TestObjectType_Serialize(t *testing.T) {
	testSchemaCaseSerialize(testObjectTypeSerializeCases)(t)
}
