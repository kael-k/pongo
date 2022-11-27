package pongo

// AnyOfType SchemaType expose a Process method which run the given SchemaAction until a Schema
// from the SchemaType list given at construction time until a Schema return a result with no error
type AnyOfType struct {
	SchemaList `json:"elements"`
}

func AnyOf(schemaElements ...SchemaType) *AnyOfType {
	return &AnyOfType{
		SchemaList: L(schemaElements).SchemaList(),
	}
}

func (e AnyOfType) Process(action SchemaAction, data *DataPointer) (processedData Data, err error) {
	var schemaError *SchemaError

	for _, v := range e.SchemaList {
		processedData, err = v.Process(action, data)
		if err == nil {
			return processedData, nil
		}

		if schemaError != nil {
			err = schemaError.MergeWithCast(data.Path(), err)
		} else {
			schemaError = NewSchemaWithCasting(data.Path(), err)
		}
	}

	return nil, err
}

func (e *AnyOfType) SchemaTypeID() (string, error) {
	return "anyOf", nil
}

func (e *AnyOfType) Schema() *Schema {
	return NewProcessableSchema(e)
}
