package pongo

import (
	"encoding/json"
	"fmt"
)

type ListType struct {
	Type   *SchemaNode          `json:"type"`
	MinLen *NumberProperty[int] `json:"minLen,omitempty"`
	MaxLen *NumberProperty[int] `json:"maxLen,omitempty"`
}

func List(schema SchemaType) *ListType {
	if schema == nil {
		return &ListType{Type: nil}
	}
	return &ListType{
		Type: Schema(schema),
	}
}

func (l ListType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	var schemaError = NewSchemaError()

	if l.Type == nil {
		return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("cannot %s data as ListType at %s, BaseSchemaType provided for \"List\" items is nil", action, dataPointer.Path()))
	}

	d, ok := dataPointer.Get().([]interface{})
	if !ok {
		return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("cannot %s data as ListType at %s, not an \"List\"", action, dataPointer.Path()))
	}

	// validate array length
	if m, ok := l.MinLen.Get(); ok && m > len(d) {
		return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("cannot %s data as ListType at %s, expected min lenght of the list at %d, got %d", action, dataPointer.Path(), l.MinLen, len(d)))
	}
	if m, ok := l.MaxLen.Get(); ok && m < len(d) {
		return nil, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("cannot %s data as ListType at %s, expected min lenght of the list at %d, got %d", action, dataPointer.Path(), l.MinLen, len(d)))
	}

	var processedSlice = []interface{}{}

	for key := range d {
		ptr := dataPointer.Push(fmt.Sprintf("[%d]", key), d[key], l.Type)
		var item interface{}

		switch action {
		case SchemaActionSerialize:
			item, err = l.Type.Serialize(ptr)
			processedSlice = append(processedSlice, item)
		case SchemaActionParse:
			item, err = l.Type.Parse(ptr)
			processedSlice = append(processedSlice, item)
		}

		if err != nil {
			schemaError = schemaError.MergeWithCast(dataPointer.Path(), err)
			continue
		}
	}

	if len(schemaError.Errors) > 0 {
		return nil, schemaError
	}

	return processedSlice, nil
}

func (l ListType) SetMinLen(i int) *ListType {
	l.MinLen = l.MinLen.Set(i)
	return &l
}

func (l ListType) SetMaxLen(i int) *ListType {
	l.MaxLen = l.MaxLen.Set(i)
	return &l
}

func (l *ListType) SchemaTypeID() string {
	return "list"
}

func (l *ListType) Children() SchemaList {
	return SchemaList{l.Type}
}

func (l ListType) MarshalJSONSchema(action SchemaAction) ([]byte, error) {
	jsonObject := map[string]interface{}{
		"type": "array",
	}

	if m, ok := l.MinLen.Get(); ok {
		jsonObject["minItems"] = m
	}
	if m, ok := l.MaxLen.Get(); ok {
		jsonObject["maxItems"] = m
	}

	if l.Type == nil {
		return json.Marshal(jsonObject)
	}

	var childJSON json.RawMessage
	var err error

	childJSON, err = MarshalJSONSchema(l.Type, action)
	if err != nil {
		return nil, err
	}

	jsonObject["items"] = childJSON

	return json.Marshal(jsonObject)
}
