package pongo

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/* These types are used to wrap SchemaType instances.
The wrappers function is to support the schema marshal/unmarshal process:
In SchemaType implementation, if is requested any SchemaType nesting, the implementation must:
* store the nested schema using one of SchemaNode (for 0..1 SchemaType), SchemaList (for SchemaType list) or SchemaMap (for Maps)
  this ensures that the marshaling is automatically done with SchemaNode.MarshalJSON.
  In this way, the SchemaType implementation is not responsible for the marshaling process
* if a SchemaType has any child, it MUST implement also ObjectSchema (if it has SchemaMap children),
  ParentSchema (for SchemaList children) or ParentSchema (if it has a single SchemaNode child). The implementation
  is required for the unmarshalling process to recursively pass the PongoSchemaUnmarshalMapper and resolve the correct
  SchemaType to unmarshal
*/

// SchemaType is SchemaNode type that expose a generic Process function
type SchemaType interface {
	Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error)
}

type SchemaNode struct {
	SchemaType

	Metadata *Metadata
	rawJSON  []byte
}

func NewEmptySchema() *SchemaNode {
	return &SchemaNode{}
}

func Schema(schema SchemaType) *SchemaNode {
	if s, ok := schema.(*SchemaNode); ok {
		return s
	}
	return &SchemaNode{
		SchemaType: schema,
	}
}

func (s SchemaNode) Parse(data *DataPointer) (Data, error) {
	return s.Process(SchemaActionParse, data)
}

func (s SchemaNode) Serialize(data *DataPointer) (Data, error) {
	return s.Process(SchemaActionSerialize, data)
}

func (s SchemaNode) Type() (schemaType SchemaType) {
	return s.SchemaType
}

func (s *SchemaNode) SetType(schemaType SchemaType) {
	s.SchemaType = schemaType
}

func (s SchemaNode) Process(action SchemaAction, data *DataPointer) (Data, error) {
	if s.SchemaType != nil {
		return s.SchemaType.Process(action, data)
	}

	return nil, ErrNoSchemaTypeSet
}

func (s *SchemaNode) MarshalJSON() ([]byte, error) {
	schemaType := s.Type()
	k := SchemaTypeID(schemaType)
	var marshalled marshalSchemaType
	var schemaTypeJSON json.RawMessage

	marshalled.Type = &k

	var err error
	schemaTypeJSON, err = json.Marshal(schemaType)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal PongoSchema: %w", err)
	}

	if string(schemaTypeJSON) != "{}" {
		marshalled.Body = &schemaTypeJSON
	}

	marshalled.Metadata = s.Metadata

	return json.Marshal(marshalled)
}

func (s *SchemaNode) UnmarshalJSON(jsonSchema []byte) error {
	s.rawJSON = jsonSchema
	return nil
}

func (s *SchemaNode) unmarshalRawJSON(mapper *PongoSchemaUnmarshalMapper) (err error) {
	defer s.cleanRawJSON()
	var unmarshal marshalSchemaType

	err = json.Unmarshal(s.rawJSON, &unmarshal)

	if err != nil {
		return err
	}

	s.Metadata = unmarshal.Metadata

	if unmarshal.Type == nil {
		return fmt.Errorf("cannot unmarshal PongoSchema, no $type set in %s", s.rawJSON)
	}

	schemaType := mapper.Get(*unmarshal.Type)
	if schemaType == nil {
		return fmt.Errorf("cannot unmarshall SchemaType element: SchemaType ID %s not found in PongoSchemaUnmarshalMapper", *unmarshal.Type)
	}

	if unmarshal.Body != nil {
		err = json.Unmarshal(*unmarshal.Body, schemaType)
		if err != nil {
			return fmt.Errorf("cannot unmarshall $body in %v: %w", s.rawJSON, err)
		}
	}

	s.SetType(schemaType)

	children, err := s.Children()
	if err != nil {
		return err
	}

	for _, c := range children {
		err = c.unmarshalRawJSON(mapper)
		c.cleanRawJSON()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SchemaNode) cleanRawJSON() {
	s.rawJSON = nil
}

func (s SchemaNode) Children() (SchemaList, error) {
	var schemaType = s.Type()
	if schemaType == nil {
		return nil, ErrNoSchemaTypeSet
	}

	schemaTypeParent, ok := schemaType.(ParentSchema)
	if ok {
		return schemaTypeParent.Children(), nil
	}

	return SchemaList{}, nil
}

type ProcessFn func(dataPointer *DataPointer, action SchemaAction) (data Data, err error)

// ParentSchema is a SchemaType type nested inside one or more schemas in a wrapped []*SchemaNode (SchemaList)
// the implementation must return all the *SchemaNode direct children
type ParentSchema interface {
	SchemaType
	// Children return all direct Children of the *SchemaNode as the original SchemaMap
	Children() SchemaList
}

type O map[string]SchemaType
type L []SchemaType

func (o O) SchemaMap() SchemaMap {
	m := SchemaMap{}
	for k, v := range o {
		if v != nil {
			m[k] = Schema(v)
		} else {
			m[k] = nil
		}
	}

	return m
}
func (l L) SchemaList() SchemaList {
	list := SchemaList{}
	for _, v := range l {
		if v != nil {
			list = append(list, Schema(v))
		} else {
			list = append(list, nil)
		}
	}

	return list
}

type SchemaMap map[string]*SchemaNode
type SchemaList []*SchemaNode

func (m SchemaMap) Children() SchemaList {
	list := SchemaList{}
	for _, v := range m {
		if v != nil {
			list = append(list, Schema(v))
		} else {
			list = append(list, nil)
		}
	}

	return list
}

func (l SchemaList) Children() SchemaList {
	return l
}

func (s SchemaNode) GetMetadata(key string) (value string, ok bool) {
	return s.Metadata.Get(key)
}

func (s *SchemaNode) SetMetadata(key string, value string) *SchemaNode {
	s.Metadata = s.Metadata.Set(key, value)
	return s
}

type Metadata map[string]string

func (m *Metadata) Get(key string) (value string, ok bool) {
	if m == nil {
		return "", false
	}
	value, ok = (*m)[key]
	return
}

func (m *Metadata) Set(key string, value string) *Metadata {
	if m == nil {
		m = &Metadata{}
	}
	(*m)[key] = value

	return m
}

// CustomSchemaTypeID is a SchemaType with a custom SchemaTypeID
// implementation should return a constant string
type CustomSchemaTypeID interface {
	SchemaType

	SchemaTypeID() string
}

func SchemaTypeID(s SchemaType) string {
	// we must remove the first char of type, which is always a `*`
	// since SchemaType is an interface

	if customSchemaTypeID, ok := s.(CustomSchemaTypeID); ok {
		return customSchemaTypeID.SchemaTypeID()
	}

	return reflect.TypeOf(s).String()[1:]
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
