package pongo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type IntType struct {
	Cast ActionFlagProperty  `json:"cast,omitempty"`
	Min  NumberProperty[int] `json:"min,omitempty"`
	Max  NumberProperty[int] `json:"max,omitempty"`
}

type IntTypeInterface interface {
	Int() int
}

func Int() *IntType {
	return &IntType{}
}

func (i IntType) cast(data *DataPointer) (n int, err error) {
	switch n := data.Get().(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case float32:
		return int(n), nil
	case float64:
		return int(n), nil
	case string:
		v, err := strconv.Atoi(n)
		if err != nil {
			return v, NewSchemaErrorWithError(data.Path(), err)
		}
		return v, nil
	case IntTypeInterface:
		return n.Int(), nil
	}

	return 0, NewSchemaErrorWithError(data.Path(), fmt.Errorf("schema does not validate: cannot cast %s to \"Int\"", data.Path()))
}

func (i IntType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	var n int

	if action != SchemaActionParse && action != SchemaActionSerialize {
		return nil, ErrInvalidAction(i, action)
	}

	if !i.Cast.GetAction(action) {
		var ok bool
		n, ok = dataPointer.Get().(int)
		if !ok {
			return 0, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s is not a \"Int\"", dataPointer.Path()))
		}
	} else {
		n, err = i.cast(dataPointer)
		if err != nil {
			return 0, err
		}
	}

	if m, ok := i.Min.Get(); ok && m > n {
		return 0, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s value is %d (Max: %d)", dataPointer.Path(), n, i.Max))
	}
	if m, ok := i.Max.Get(); ok && m < n {
		return 0, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %d (min: %d)", dataPointer.Path(), n, i.Min))
	}

	return n, nil
}

func (i IntType) SetCast(cast bool) *IntType {
	i.Cast.Set(cast)
	return &i
}

func (i IntType) SetCastActions(actions ...SchemaAction) *IntType {
	i.Cast.SetActions(actions...)
	return &i
}

func (i IntType) UnsetCastActions(actions ...SchemaAction) *IntType {
	i.Cast.UnsetActions(actions...)
	return &i
}

func (i IntType) SetMin(n int) *IntType {
	i.Min.Set(n)
	return &i
}

func (i IntType) SetMax(n int) *IntType {
	i.Max.Set(n)
	return &i
}

func (i *IntType) SchemaTypeID() string {
	return "int"
}

func (i IntType) MarshalJSON() ([]byte, error) {
	var d = map[string]interface{}{}

	if !i.Cast.Empty() {
		d["cast"] = i.Cast
	}

	if m, ok := i.Min.Get(); ok {
		d["min"] = m
	}

	if m, ok := i.Max.Get(); ok {
		d["max"] = m
	}

	return json.Marshal(d)
}
