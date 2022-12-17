package pongo

import (
	"reflect"
)

// SchemaType is SchemaNode type that expose a generic Process function
type SchemaType interface {
	Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error)
}

type ProcessFn func(dataPointer *DataPointer, action SchemaAction) (data Data, err error)

// CustomSchemaTypeID is a SchemaType with a custom SchemaTypeID
// implementation should return a constant string
type CustomSchemaTypeID interface {
	SchemaType

	SchemaTypeID() string
}

// ParentSchema is a SchemaType type nested inside one or more schemas in a wrapped []*SchemaNode (SchemaList)
// the implementation must return all the *SchemaNode direct children
type ParentSchema interface {
	SchemaType
	// Children return all direct Children of the *SchemaNode as the original SchemaMap
	Children() SchemaList
}

// JSONSchemaMarshaler is a SchemaType which can be marshaled into a jsonschema
type JSONSchemaMarshaler interface {
	SchemaType

	MarshalJSONSchema(action SchemaAction) ([]byte, error)
}

func SchemaTypeID(s SchemaType) string {
	// we must remove the first char of type, which is always a `*`
	// since SchemaType is an interface

	if customSchemaTypeID, ok := s.(CustomSchemaTypeID); ok {
		return customSchemaTypeID.SchemaTypeID()
	}

	return reflect.TypeOf(s).String()[1:]
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

type SchemaAction string

const (
	SchemaActionParse     SchemaAction = "PARSE"
	SchemaActionSerialize SchemaAction = "SERIALIZE"
)

// Parse is wrapper for SchemaNode.Parse that automatically
// transforms Data into a DataPointer
func Parse(schema SchemaType, data Data) (Data, error) {
	return Process(schema, SchemaActionParse, data)
}

// Serialize is wrapper for SchemaNode.Serialize that automatically
// transforms Data into a DataPointer
func Serialize(schema SchemaType, data Data) (Data, error) {
	return Process(schema, SchemaActionSerialize, data)
}

// Process is wrapper for SchemaNode.Process that automatically
// transforms Data into a DataPointer
func Process(schema SchemaType, action SchemaAction, data Data) (Data, error) {
	return schema.Process(action, NewDataPointer(data, schema))
}
