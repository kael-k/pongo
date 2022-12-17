package pongo

import (
	"encoding/json"
	"errors"
)

func MarshalPongoSchema(schema SchemaType) ([]byte, error) {
	return MarshalPongoSchemaWithMetadata(schema, nil)
}

func MarshalPongoSchemaWithMetadata(schema SchemaType, metadata *Metadata) ([]byte, error) {
	d := map[string]interface{}{
		"$version": "1.0",
		"$body":    Schema(schema),
	}
	if metadata != nil {
		d["$metadata"] = metadata
	}
	return json.Marshal(d)
}

type marshalSchemaType struct {
	Body     *json.RawMessage `json:"$body,omitempty"`
	Metadata *Metadata        `json:"$metadata,omitempty"`
	Type     *string          `json:"$type"`
}

func UnmarshalPongoSchema(jsonSchema []byte) (schema *SchemaNode, metadata *Metadata, err error) {
	return UnmarshalSchemaWithMapper(jsonSchema, GlobalSchemaUnmarshalMap())
}

func UnmarshalSchemaWithMapper(jsonSchema []byte, mapper *SchemaUnmarshalMapper) (schema *SchemaNode, metadata *Metadata, err error) {
	var root *map[string]json.RawMessage

	err = json.Unmarshal(jsonSchema, &root)
	if err != nil {
		return nil, nil, err
	}

	version, ok := (*root)["$version"]
	if !ok {
		return nil, nil, errors.New("expected schema version \"1.0\" in JSON, no version found")
	}

	if string(version) != "\"1.0\"" {
		return nil, nil, errors.New("expected schema version \"1.0\" in JSON, found " + string(version))
	}

	jsonSchema, ok = (*root)["$body"]
	if !ok {
		return nil, nil, errors.New("expected schema body in JSON, no schema found")
	}

	schema = NewEmptySchema()
	schema.rawJSON = jsonSchema

	jsonMetadata, ok := (*root)["$metadata"]
	if ok {
		metadata = &Metadata{}
		err = json.Unmarshal(jsonMetadata, metadata)
		if err != nil {
			return nil, nil, err
		}
	}

	return schema, metadata, schema.unmarshalRawJSON(mapper)
}

type SchemaFactory func() SchemaType

type SchemaUnmarshalMapper struct {
	schemaElementsMap map[string]SchemaFactory
}

func SchemaUnmarshalMap() *SchemaUnmarshalMapper {
	return &SchemaUnmarshalMapper{
		map[string]SchemaFactory{},
	}
}

func (p SchemaUnmarshalMapper) SchemaElements() map[string]SchemaFactory {
	return p.schemaElementsMap
}

func (p SchemaUnmarshalMapper) Get(schemaElementID string) SchemaType {
	schemaType, ok := p.schemaElementsMap[schemaElementID]
	if !ok {
		return nil
	}
	return schemaType()
}

func (p SchemaUnmarshalMapper) Set(schema SchemaFactory) *SchemaUnmarshalMapper {
	id := SchemaTypeID(schema())
	p.schemaElementsMap[id] = schema

	return &p
}

func (p SchemaUnmarshalMapper) Clone() *SchemaUnmarshalMapper {
	cloneMapper := SchemaUnmarshalMap()
	for key, el := range p.schemaElementsMap {
		cloneMapper.schemaElementsMap[key] = el
	}

	return cloneMapper
}

var globalUnmarshalMapper = &SchemaUnmarshalMapper{
	map[string]SchemaFactory{
		"anyOf":    func() SchemaType { return AnyOf(nil) },
		"oneOf":    func() SchemaType { return OneOf(nil) },
		"allOf":    func() SchemaType { return AllOf(nil) },
		"list":     func() SchemaType { return List(nil) },
		"object":   func() SchemaType { return Object(nil) },
		"string":   func() SchemaType { return String() },
		"int":      func() SchemaType { return Int() },
		"float64":  func() SchemaType { return Float64() },
		"bytes":    func() SchemaType { return Bytes() },
		"bool":     func() SchemaType { return Bool() },
		"datetime": func() SchemaType { return Datetime() },
	},
}

func GlobalSchemaUnmarshalMap() *SchemaUnmarshalMapper {
	return globalUnmarshalMapper.Clone()
}

func SetGlobalSchemaUnmarshalMap(newGlobalScheme *SchemaUnmarshalMapper) {
	globalUnmarshalMapper = newGlobalScheme
}
