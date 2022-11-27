package pongo

import (
	"encoding/json"
	"errors"
)

func MarshalSchemaJSON(schema SchemaType) ([]byte, error) {
	return MarshalSchemaJSONWithMetadata(schema, nil)
}

func MarshalSchemaJSONWithMetadata(schema SchemaType, metadata *Metadata) ([]byte, error) {
	d := map[string]interface{}{
		"version": "1.0",
		"schema":  schema.Schema(),
	}
	if metadata != nil {
		d["metadata"] = metadata
	}
	return json.Marshal(d)
}

type marshalSchemaType struct {
	Body     *json.RawMessage `json:"$body,omitempty"`
	Metadata *Metadata        `json:"$metadata,omitempty"`
	Type     *string          `json:"$type"`
}

func UnmarshalSchemaJSON(jsonSchema []byte) (schema *Schema, metadata *Metadata, err error) {
	return UnmarshalSchemaJSONWithMapper(jsonSchema, DefaultSchemaUnmarshalMap())
}

func UnmarshalSchemaJSONWithMapper(jsonSchema []byte, mapper *SchemaUnmarshalMapper) (schema *Schema, metadata *Metadata, err error) {
	var root *map[string]json.RawMessage

	err = json.Unmarshal(jsonSchema, &root)
	if err != nil {
		return nil, nil, err
	}

	version, ok := (*root)["version"]
	if !ok {
		return nil, nil, errors.New("expected schema version \"1.0\" in JSON, no version found")
	}

	if string(version) != "\"1.0\"" {
		return nil, nil, errors.New("expected schema version \"1.0\" in JSON, found " + string(version))
	}

	jsonSchema, ok = (*root)["schema"]
	if !ok {
		return nil, nil, errors.New("expected schema body in JSON, no schema found")
	}

	schema = NewEmptySchema()
	schema.rawJSON = jsonSchema

	jsonMetadata, ok := (*root)["metadata"]
	if ok {
		metadata = &Metadata{}
		err = json.Unmarshal(jsonMetadata, metadata)
		if err != nil {
			return nil, nil, err
		}
	}

	return schema, metadata, schema.unmarshalRawJSON(mapper)
}
