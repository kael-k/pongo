package pongo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BoolType struct {
	Cast ActionFlagProperty `json:"cast,omitempty"`
}

func Bool() *BoolType {
	return &BoolType{}
}

func (b BoolType) Process(action SchemaAction, data *DataPointer) (Data, error) {
	var parsedBool bool
	var err error

	if b.Cast.GetAction(action) {
		switch r := data.Get().(type) {
		case string:
			switch strings.ToLower(r) {
			case "true":
				parsedBool = true
			case "false":
				parsedBool = false
			default:
				return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s cannot cast to \"Bool\"", data.Path()))
			}
		case byte:
			parsedBool = r != byte(0)
		case []byte:
			parsedBool = len(r) > 1 || len(r) == 1 && r[0] != byte(0)
		case uint, uint64, int, int64, float64, float32:
			parsedBool = r != 0
		default:
			return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s cannot cast to \"Bool\"", data.Path()))
		}
	} else {
		var ok bool
		parsedBool, ok = data.Get().(bool)
		if !ok {
			err = NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s is not a \"Bool\"", data.Path()))
		}
	}

	if err != nil {
		return nil, err
	}

	if action == SchemaActionSerialize || action == SchemaActionParse {
		return parsedBool, nil
	}

	return nil, NewSchemaErrorWithError(data.Path(), ErrInvalidAction(&b, action))
}

func (b *BoolType) SetCast(cast bool) *BoolType {
	b.Cast.Set(cast)
	return b
}

func (b BoolType) SetCastActions(actions ...SchemaAction) *BoolType {
	b.Cast.SetActions(actions...)
	return &b
}

func (b BoolType) UnsetCastActions(actions ...SchemaAction) *BoolType {
	b.Cast.UnsetActions(actions...)
	return &b
}

func (b *BoolType) SchemaTypeID() (string, error) {
	return "bool", nil
}

func (b *BoolType) Schema() *Schema {
	return NewProcessableSchema(b)
}

func (b BoolType) MarshalJSON() ([]byte, error) {
	var d = map[string]interface{}{}
	if !b.Cast.Empty() {
		d["cast"] = b.Cast
	}

	return json.Marshal(d)
}
