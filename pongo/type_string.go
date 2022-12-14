package pongo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringType struct {
	Cast   *ActionFlagProperty  `json:"cast,omitempty"`
	MinLen *NumberProperty[int] `json:"minLen,omitempty"`
	MaxLen *NumberProperty[int] `json:"maxLen,omitempty"`
}

func String() *StringType {
	return &StringType{}
}

func (s StringType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	var str string
	if s.Cast.GetAction(action) {
		switch r := dataPointer.Get().(type) {
		case string:
			str = r
		case int64:
			str = strconv.FormatInt(r, 10)
		case int:
			str = strconv.FormatInt(int64(r), 10)
		case float32:
			str = strconv.FormatFloat(float64(r), 'f', -1, 32)
		case float64:
			str = strconv.FormatFloat(r, 'f', -1, 64)
		case fmt.Stringer:
			str = r.String()
		default:
			return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s cannot cast to \"String\"", dataPointer.Path()))
		}

	} else {
		var ok bool
		str, ok = dataPointer.Get().(string)
		if !ok {
			return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s is not a \"String\"", dataPointer.Path()))
		}
	}

	if l, ok := s.MinLen.Get(); ok && l > len(str) {
		return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %d (min: %d)", dataPointer.Path(), len(str), l))
	}
	if l, ok := s.MaxLen.Get(); ok && l < len(str) {
		return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %d (Max: %d)", dataPointer.Path(), len(str), l))
	}

	return str, nil
}

func (s StringType) SetCast(cast bool) *StringType {
	s.Cast = s.Cast.Set(cast)
	return &s
}

func (s StringType) SetCastActions(actions ...SchemaAction) *StringType {
	s.Cast = s.Cast.SetActions(actions...)
	return &s
}

func (s StringType) UnsetCastActions(actions ...SchemaAction) *StringType {
	s.Cast.UnsetActions(actions...)
	return &s
}

func (s StringType) SetMinLen(i int) *StringType {
	s.MinLen = s.MinLen.Set(i)
	return &s
}

func (s StringType) SetMaxLen(i int) *StringType {
	s.MaxLen = s.MaxLen.Set(i)
	return &s
}

func (s *StringType) SchemaTypeID() string {
	return "string"
}

func (s StringType) MarshalJSONSchema(action SchemaAction) ([]byte, error) {
	if action != SchemaActionParse && action != SchemaActionSerialize {
		return nil, NewErrInvalidAction(s, action)
	}

	var jsonObject = map[string]interface{}{"type": "string"}

	if l, ok := s.MinLen.Get(); ok {
		jsonObject["minLength"] = l
	}
	if l, ok := s.MaxLen.Get(); ok {
		jsonObject["maxLength"] = l
	}

	if !s.Cast.GetAction(action) {
		return json.Marshal(jsonObject)
	}

	return json.Marshal(map[string]interface{}{
		"oneOf": []map[string]interface{}{
			jsonObject,
			{
				"type": "number",
			},
		},
	})
}
