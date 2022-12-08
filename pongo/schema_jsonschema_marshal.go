package pongo

import (
	"encoding/json"
	"fmt"
)

const jsonSchemaDraft07Schema string = "http://json-schema.org/draft-07/schema#"

func MarshalJSONSchemaWithMetadata(schema *SchemaNode, action SchemaAction) ([]byte, error) {
	var jsonObject map[string]json.RawMessage
	var metadata = schema.Metadata

	var jsonBytes, err = MarshalJSONSchema(schema, action)
	if err != nil {
		return nil, err
	}

	if jsonBytes != nil {
		err = json.Unmarshal(jsonBytes, &jsonObject)
		if err != nil {
			return nil, err
		}
	} else {
		jsonObject = map[string]json.RawMessage{}
	}

	var id, ok = metadata.Get("$id")
	if ok {
		jsonObject["$id"] = []byte(id)
	}

	jsonObject["$schema"] = []byte(fmt.Sprintf("\"%s\"", jsonSchemaDraft07Schema))
	return json.Marshal(jsonObject)
}

func MarshalJSONSchema(schema *SchemaNode, action SchemaAction) ([]byte, error) {
	var schemaType, ok = schema.Type().(JSONSchemaMarshaler)
	if !ok {
		return nil, nil
	}

	return schemaType.MarshalJSONSchema(action)
}
