package pongo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type BytesType struct {
	Cast   ActionFlagProperty  `json:"cast,omitempty"`
	MaxLen NumberProperty[int] `json:"maxLen,omitempty"`
	MinLen NumberProperty[int] `json:"mixLen,omitempty"`
}

func Bytes() *BytesType {
	return &BytesType{}
}

func (b BytesType) Parse(data *DataPointer) (parsedData Data, err error) {
	return b.parse(data, SchemaActionParse)
}

func (b BytesType) Serialize(data *DataPointer) (serializedData Data, err error) {
	bytes, err := b.parse(data, SchemaActionSerialize)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (b BytesType) parse(data *DataPointer, action SchemaAction) (bytes []byte, err error) {
	if b.Cast.GetAction(action) {
		switch r := data.Get().(type) {
		case string:
			bytes, err = base64.StdEncoding.DecodeString(r)
			if err != nil {
				return nil, NewSchemaErrorWithError(data.Path(), err)
			}
		case []byte:
			bytes = r
		case byte:
			bytes = []byte{r}
		default:
			return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s cannot cast to String", data.Path()))
		}

	} else {
		var ok bool
		bytes, ok = data.Get().([]byte)
		if !ok {
			return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s is not a \"bytes\"", data.Path()))
		}
	}

	if l, ok := b.MinLen.Get(); ok && l > len(bytes) {
		return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s length is %d (min: %d)", data.Path(), len(bytes), b.MinLen))
	}
	if l, ok := b.MaxLen.Get(); ok && l < len(bytes) {
		return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: %s length is %d (Max: %d)", data.Path(), len(bytes), b.MaxLen))
	}

	return bytes, nil
}

func (b BytesType) SetMinLen(i int) *BytesType {
	b.MinLen.Set(i)
	return &b
}

func (b BytesType) SetMaxLen(i int) *BytesType {
	b.MaxLen.Set(i)
	return &b
}

func (b BytesType) SetCast(cast bool) *BytesType {
	b.Cast.Set(cast)
	return &b
}

func (b BytesType) SetCastActions(actions ...SchemaAction) *BytesType {
	b.Cast.SetActions(actions...)
	return &b
}

func (b BytesType) UnsetCastActions(actions ...SchemaAction) *BytesType {
	b.Cast.UnsetActions(actions...)
	return &b
}

func (b *BytesType) SchemaTypeID() (string, error) {
	return "bytes", nil
}

func (b *BytesType) Schema() *Schema {
	return NewBaseSchema(b)
}

func (b BytesType) MarshalJSON() ([]byte, error) {
	var d = map[string]interface{}{}
	if !b.Cast.Empty() {
		d["cast"] = b.Cast
	}

	if m, ok := b.MinLen.Get(); ok {
		d["minLen"] = m
	}

	if m, ok := b.MaxLen.Get(); ok {
		d["maxLen"] = m
	}

	return json.Marshal(d)
}
