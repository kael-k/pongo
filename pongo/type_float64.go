package pongo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Float64Type struct {
	Cast *ActionFlagProperty      `json:"cast,omitempty"`
	Min  *NumberProperty[float64] `json:"min,omitempty"`
	Max  *NumberProperty[float64] `json:"max,omitempty"`
}

func Float64() *Float64Type {
	return &Float64Type{}
}

func (f64 Float64Type) cast(data *DataPointer) (n float64, err *SchemaError) {
	switch n := data.Get().(type) {
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case float32:
		return float64(n), nil
	case float64:
		return n, nil
	case string:
		v, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return v, NewSchemaErrorWithError(data.Path(), err)
		}
		return v, nil
	}

	return 0, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: cannot cast %s to \"Float64\"", data.Path()))
}

func (f64 Float64Type) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	var n float64

	if action != SchemaActionParse && action != SchemaActionSerialize {
		return nil, ErrInvalidAction(f64, action)
	}

	if !f64.Cast.GetAction(action) {
		var ok bool
		n, ok = dataPointer.Get().(float64)
		if !ok {
			return 0, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s is not a \"Float64\"", dataPointer.Path()))
		}
	} else {
		var schemaErr *SchemaError
		n, schemaErr = f64.cast(dataPointer)
		if schemaErr != nil {
			return 0, schemaErr
		}
	}

	if m, ok := f64.Min.Get(); ok && m > n {
		return 0, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s value is %f (Min: %f)", dataPointer.Path(), n, m))
	}
	if m, ok := f64.Max.Get(); ok && m < n {
		return 0, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %f (Max: %f)", dataPointer.Path(), n, m))
	}

	return n, nil
}

func (f64 Float64Type) SetCast(cast bool) *Float64Type {
	f64.Cast = f64.Cast.Set(cast)
	return &f64
}

func (f64 Float64Type) SetCastActions(actions ...SchemaAction) *Float64Type {
	f64.Cast = f64.Cast.SetActions(actions...)
	return &f64
}

func (f64 Float64Type) UnsetCastActions(actions ...SchemaAction) *Float64Type {
	f64.Cast.UnsetActions(actions...)
	return &f64
}

func (f64 Float64Type) SetMin(f float64) *Float64Type {
	f64.Min = f64.Min.Set(f)
	return &f64
}

func (f64 Float64Type) SetMax(f float64) *Float64Type {
	f64.Max = f64.Max.Set(f)
	return &f64
}

func (f64 *Float64Type) SchemaTypeID() string {
	return "float64"
}

func (f64 Float64Type) MarshalJSONSchema(action SchemaAction) ([]byte, error) {
	if action != SchemaActionParse && action != SchemaActionSerialize {
		return nil, ErrInvalidAction(f64, action)
	}

	var jsonObject = map[string]interface{}{"type": "number"}

	if l, ok := f64.Min.Get(); ok {
		jsonObject["minimum"] = l
	}
	if l, ok := f64.Max.Get(); ok {
		jsonObject["maximum"] = l
	}

	if !f64.Cast.GetAction(action) {
		return json.Marshal(jsonObject)
	}

	return json.Marshal(map[string]interface{}{
		"oneOf": []map[string]interface{}{
			jsonObject,
			{
				"type":    "string",
				"pattern": "^-?([0-9]+)(.[0-9]+)?$",
			},
		},
	})
}
