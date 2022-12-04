package pongo

import (
	"fmt"
	"reflect"
)

// SchemaType is an interface that should represent logically one of
// BaseSchemaType or ProcessableSchemaType, but golang at the moment
// does not allow type union on behavioural interfaces
type SchemaType interface {
	Schema() *Schema
}

// BaseSchemaType is SchemaType that can Parse and Serialize a *DataPointer
type BaseSchemaType interface {
	SchemaType

	// Parse the DataPointer, if the BaseSchemaType does not validate, it will return an error
	Parse(data *DataPointer) (parsedData Data, err error)

	// Serialize implementation must return a new Data that is directly mappable as JSON: string, int,
	// float64, bool and their map[string] and slices, and slices and map[string] of slices and map[string]
	Serialize(data *DataPointer) (serializedData Data, err error)
}

// ProcessableSchemaType is Schema type that expose a generic Process function
type ProcessableSchemaType interface {
	SchemaType

	Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error)
}

type ProcessFn func(dataPointer *DataPointer, action SchemaAction) (data Data, err error)

// CustomSchemaTypeID is a SchemaType with a custom SchemaTypeID
// implementation should return a constant string
type CustomSchemaTypeID interface {
	SchemaType

	SchemaTypeID() (string, error)
}

// ParentSchema is a Schema type nested inside one or more schemas in a wrapped []*Schema (SchemaList)
// the implementation must return all the *Schema direct children
type ParentSchema interface {
	SchemaType
	// Children return all direct Children of the *Schema as the original SchemaMap
	Children() SchemaList
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

var defaultSchemaList = map[string]SchemaFactory{
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
}

func DefaultSchemaUnmarshalMap() *SchemaUnmarshalMapper {
	var err error
	s := SchemaUnmarshalMap()
	for k, schemaFactory := range defaultSchemaList {
		s, err = s.Set(schemaFactory)
		if err != nil {
			panic(fmt.Errorf("cannot generate DefaultSchemaUnmarshalMap, unexpected error for schema %s: %w", k, err))
		}
	}

	return s
}

func SchemaTypeID(s SchemaType) (string, error) {
	// we must remove the first char of type, which is always a `*`
	// since SchemaType is an interface

	switch v := s.(type) {
	case CustomSchemaTypeID:
		return v.SchemaTypeID()
	case BaseSchemaType, ProcessableSchemaType:
		return reflect.TypeOf(s).String()[1:], nil
	default:
		return "", ErrInvalidSchemaType
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

func (p SchemaUnmarshalMapper) Set(schema SchemaFactory) (*SchemaUnmarshalMapper, error) {
	id, err := SchemaTypeID(schema())
	if err != nil {
		return nil, err
	}
	p.schemaElementsMap[id] = schema

	return &p, nil
}

type SchemaAction string

const (
	SchemaActionParse     SchemaAction = "PARSE"
	SchemaActionSerialize SchemaAction = "SERIALIZE"
)

// Parse is wrapper for Schema.Parse that automatically
// transforms Data into a DataPointer
func Parse(schema SchemaType, data Data) (Data, error) {
	return schema.Schema().Parse(NewDataPointer(data, schema))
}

// Serialize is wrapper for Schema.Serialize that automatically
// transforms Data into a DataPointer
func Serialize(schema SchemaType, data Data) (Data, error) {
	return schema.Schema().Serialize(NewDataPointer(data, schema))
}

// Process is wrapper for Schema.Process that automatically
// transforms Data into a DataPointer
func Process(schema SchemaType, action SchemaAction, data Data) (Data, error) {
	return schema.Schema().Process(action, NewDataPointer(data, schema))
}
