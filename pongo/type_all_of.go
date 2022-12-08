package pongo

import "encoding/json"

// AllOfType SchemaType expose a Process method which run the given SchemaAction on all SchemaNode
// from the SchemaType list given at construction time return the result of the last SchemaType if
// no error is encountered during the processing of the previous SchemaNode
type AllOfType struct {
	SchemaList `json:"elements"`
	Chain      ActionFlagProperty `json:"chain"`
}

func AllOf(schemaElements ...SchemaType) *AllOfType {
	return &AllOfType{
		SchemaList: L(schemaElements).SchemaList(),
	}
}

func (e AllOfType) Process(action SchemaAction, data *DataPointer) (processedData Data, err error) {
	var schemaError *SchemaError

	for _, v := range e.SchemaList {
		processedData, err = v.Process(action, data)
		if err != nil {
			if schemaError != nil {
				schemaError = schemaError.MergeWithCast(data.Path(), err)
			} else {
				schemaError = NewSchemaWithCasting(data.Path(), err)
			}

			continue
		}

		if e.Chain.GetAction(action) {
			err = data.Path().SetOverride(processedData)
			if err != nil {
				if schemaError != nil {
					schemaError = schemaError.MergeWithCast(data.Path(), err)
				} else {
					schemaError = NewSchemaWithCasting(data.Path(), err)
				}
			}
		}
	}

	if schemaError != nil {
		return nil, schemaError
	}

	return processedData, nil
}

func (e AllOfType) SetChain(cast bool) *AllOfType {
	e.Chain.Set(cast)
	return &e
}

func (e AllOfType) SetChainActions(actions ...SchemaAction) *AllOfType {
	e.Chain.SetActions(actions...)
	return &e
}

func (e AllOfType) UnsetChainActions(actions ...SchemaAction) *AllOfType {
	e.Chain.UnsetActions(actions...)
	return &e
}

func (e *AllOfType) SchemaTypeID() string {
	return "allOf"
}

func (e AllOfType) MarshalJSONSchema(action SchemaAction) ([]byte, error) {
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
		"allOf": childrenJSON,
	})
}
