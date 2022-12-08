package pongo

import "encoding/json"

// AnyOfType SchemaType expose a Process method which run the given SchemaAction until a SchemaNode
// from the SchemaType list given at construction time until a SchemaNode return a result with no error
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

func (e *AnyOfType) SchemaTypeID() string {
	return "anyOf"
}

func (e AnyOfType) MarshalJSONSchema(action SchemaAction) ([]byte, error) {
	var childrenJSON []json.RawMessage

	for _, child := range e.Children() {
		j, err := MarshalJSONSchema(child, action)
		if err != nil {
			return nil, err
		}
		if j == nil {
			continue
		}
		childrenJSON = append(childrenJSON, j)
	}

	return json.Marshal(map[string]interface{}{
		"anyOf": childrenJSON,
	})
}
