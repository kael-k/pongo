package pongo

import "fmt"

// OneOfType SchemaType expose a Process method which run the given SchemaAction on all Schema(s)
// from the SchemaType list given at construction time, exactly one schema must process with no error
type OneOfType struct {
	SchemaList `json:"elements"`
}

func OneOf(schemaElements ...SchemaType) *OneOfType {
	return &OneOfType{
		SchemaList: L(schemaElements).SchemaList(),
	}
}

func (e OneOfType) Process(action SchemaAction, data *DataPointer) (processedData Data, err error) {
	var schemaError *SchemaError
	var processed = false

	for _, v := range e.SchemaList {
		processedCaseData, err := v.Process(action, data)
		if err != nil {
			if schemaError != nil {
				schemaError = schemaError.MergeWithCast(data.Path(), err)
			} else {
				schemaError = NewSchemaWithCasting(data.Path(), err)
			}
		} else {
			if processed {
				return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("cannot serialize %s, multiple types match the schema, expected exactly one match", data.Path()))
			}
			processedData = processedCaseData
			processed = true
		}
	}

	if !processed {
		return nil, NewSchemaErrorWithError(data.Path(), fmt.Errorf("cannot serialize %s, no type match the schema, expected exactly one match", data.Path()))
	}

	return processedData, nil
}

func (e OneOfType) SchemaTypeID() (string, error) {
	return "oneOf", nil
}

func (e *OneOfType) Schema() *Schema {
	return NewProcessableSchema(e)
}
