package pongo

import (
	"encoding/base64"
	"fmt"
)

type BytesType struct {
	Cast   *ActionFlagProperty  `json:"cast,omitempty"`
	MaxLen *NumberProperty[int] `json:"maxLen,omitempty"`
	MinLen *NumberProperty[int] `json:"mixLen,omitempty"`
}

func Bytes() *BytesType {
	return &BytesType{}
}

func (b BytesType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	var bytes []byte
	if b.Cast.GetAction(action) {
		switch r := dataPointer.Get().(type) {
		case string:
			bytes, err = base64.StdEncoding.DecodeString(r)
			if err != nil {
				return nil, NewSchemaErrorWithError(dataPointer.Path(), err)
			}
		case []byte:
			bytes = r
		case byte:
			bytes = []byte{r}
		default:
			return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s cannot cast to String", dataPointer.Path()))
		}

	} else {
		var ok bool
		bytes, ok = dataPointer.Get().([]byte)
		if !ok {
			return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s is not a \"bytes\"", dataPointer.Path()))
		}
	}

	if l, ok := b.MinLen.Get(); ok && l > len(bytes) {
		return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %d (min: %d)", dataPointer.Path(), len(bytes), b.MinLen))
	}
	if l, ok := b.MaxLen.Get(); ok && l < len(bytes) {
		return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %d (Max: %d)", dataPointer.Path(), len(bytes), b.MaxLen))
	}

	switch action {
	case SchemaActionSerialize:
		return base64.StdEncoding.EncodeToString(bytes), nil
	case SchemaActionParse:
		return bytes, nil
	}

	return nil, ErrInvalidAction(b, action)
}

func (b BytesType) SetMinLen(i int) *BytesType {
	b.MinLen = b.MinLen.Set(i)
	return &b
}

func (b BytesType) SetMaxLen(i int) *BytesType {
	b.MaxLen = b.MaxLen.Set(i)
	return &b
}

func (b BytesType) SetCast(cast bool) *BytesType {
	b.Cast = b.Cast.Set(cast)
	return &b
}

func (b BytesType) SetCastActions(actions ...SchemaAction) *BytesType {
	b.Cast = b.Cast.SetActions(actions...)
	return &b
}

func (b BytesType) UnsetCastActions(actions ...SchemaAction) *BytesType {
	b.Cast.UnsetActions(actions...)
	return &b
}

func (b *BytesType) SchemaTypeID() string {
	return "bytes"
}
