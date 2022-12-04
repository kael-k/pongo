package pongo

import (
	"encoding/json"
	"fmt"
)

type ObjectType struct {
	SchemaMap `json:"properties"`
	Required  []string `json:"required"`
}

func Object(properties O) *ObjectType {
	return &ObjectType{
		SchemaMap: properties.SchemaMap(),
		Required:  []string{},
	}
}

func (o ObjectType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	d, ok := dataPointer.Get().(map[string]interface{})
	if !ok {
		return nil, NewSchemaError().Append(dataPointer.Path(), fmt.Errorf("cannot validate data as ObjectType at %s, not an \"Object\"", dataPointer.Path()))
	}

	var schemaError = NewSchemaError()

	// check required keys
	if len(o.Required) > 0 {
		if diff := ListMapDiff[string](o.Required, d); len(diff) > 0 {
			return nil, NewSchemaError().Append(
				dataPointer.Path(),
				fmt.Errorf(
					"cannot validate data as ObjectType at %s, missing required properties %v", dataPointer.Path(),
					diff,
				),
			)
		}
	}

	var processedObject = map[string]interface{}{}
	for key := range d {
		// load BaseSchemaType, run all pre-checks related to ObjectType
		validator, ok := o.SchemaMap[key]
		if !ok {
			schemaError = schemaError.Append(dataPointer.Path(), fmt.Errorf("cannot validate data as ObjectType at %s, cannot get key %s", dataPointer.Path(), key))
			continue
		}

		// navigate the DataPointer
		ptr := dataPointer.Push(key, d[key], &o)

		switch action {
		case SchemaActionSerialize:
			processedObject[key], err = validator.Serialize(ptr)
		case SchemaActionParse:
			processedObject[key], err = validator.Parse(ptr)
		}

		if err != nil {
			schemaError = schemaError.MergeWithCast(dataPointer.Path(), err)
			continue
		}
	}

	if len(schemaError.Errors) > 0 {
		return nil, schemaError
	}

	return processedObject, nil
}

func (o ObjectType) Require(requires ...string) *ObjectType {
	o.Required = requires
	return &o
}

func (o *ObjectType) SchemaTypeID() string {
	return "object"
}

func (o ObjectType) MarshalJSON() ([]byte, error) {
	var d = map[string]interface{}{}
	d["properties"] = o.SchemaMap

	if len(o.Required) > 0 {
		d["required"] = o.Required
	}

	return json.Marshal(d)
}
