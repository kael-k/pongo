package pongo

import (
	"encoding/json"
	"fmt"
)

/* These pongo are used to wrap SchemaType instances.
The wrappers function is to support the schema marshal/unmarshal process:
In SchemaType implementation, if is requested any SchemaType nesting, the implementation must:
* store the nested schema using one of Schema (for 0..1 SchemaType), SchemaList (for SchemaType list) or SchemaMap (for Maps)
  this ensures that the marshaling is automatically done with Schema.MarshalJSON.
  In this way, the SchemaType implementation is not responsible for the marshaling process
* if a SchemaType has any child, it MUST implement also ObjectSchema (if it has SchemaMap children),
  ListSchema (for SchemaList children) or ParentSchema (if it has a single Schema child). The implementation
  is required for the unmarshalling process to recursively pass the SchemaUnmarshalMapper and resolve the correct
  SchemaType to unmarshal
*/

type O map[string]SchemaType
type L []SchemaType

func (o O) SchemaMap() SchemaMap {
	m := SchemaMap{}
	for k, v := range o {
		if v != nil {
			m[k] = v.Schema()
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
			list = append(list, v.Schema())
		} else {
			list = append(list, nil)
		}
	}

	return list
}

type SchemaMap map[string]*Schema
type SchemaList []*Schema

type Schema struct {
	ProcessableSchemaType
	BaseSchemaType

	Metadata *Metadata
	rawJSON  []byte
}

func NewEmptySchema() *Schema {
	return &Schema{}
}

func NewSchema(schema SchemaType) (s *Schema, err error) {
	s = NewEmptySchema()
	err = s.SetType(schema)

	return
}

func NewBaseSchema(schema BaseSchemaType) *Schema {
	return &Schema{
		BaseSchemaType: schema,
	}
}

func NewProcessableSchema(schema ProcessableSchemaType) *Schema {
	return &Schema{
		ProcessableSchemaType: schema,
	}
}

func (s Schema) IsProcessableSchema() bool {
	return s.ProcessableSchemaType != nil
}

func (s Schema) IsSchemaType() bool {
	return s.BaseSchemaType != nil
}

func (s Schema) Parse(data *DataPointer) (Data, error) {
	if s.IsSchemaType() {
		return s.BaseSchemaType.Parse(data)
	}
	if s.IsProcessableSchema() {
		return s.ProcessableSchemaType.Process(SchemaActionParse, data)
	}

	return nil, ErrNoSchemaTypeSet
}

func (s Schema) Serialize(data *DataPointer) (Data, error) {
	if s.IsSchemaType() {
		return s.BaseSchemaType.Serialize(data)
	}
	if s.IsProcessableSchema() {
		return s.ProcessableSchemaType.Process(SchemaActionSerialize, data)
	}

	return nil, ErrNoSchemaTypeSet
}

func (s Schema) Type() (schemaType SchemaType) {
	if s.IsSchemaType() {
		return s.BaseSchemaType
	}
	if s.IsProcessableSchema() {
		return s.ProcessableSchemaType
	}

	return nil
}

func (s *Schema) SetType(schemaType interface{}) error {
	switch t := schemaType.(type) {
	case BaseSchemaType:
		s.BaseSchemaType = t
		s.ProcessableSchemaType = nil
	case ProcessableSchemaType:
		s.BaseSchemaType = nil
		s.ProcessableSchemaType = t
	default:
		return ErrInvalidSchemaType
	}
	return nil
}

func (s Schema) Process(action SchemaAction, data *DataPointer) (Data, error) {
	if s.IsProcessableSchema() {
		return s.ProcessableSchemaType.Process(action, data)
	}
	if s.IsSchemaType() {
		switch action {
		case SchemaActionParse:
			return s.BaseSchemaType.Parse(data)
		case SchemaActionSerialize:
			return s.BaseSchemaType.Serialize(data)
		default:
			return nil, fmt.Errorf("cannot Process schema at path %s: cannot run action %s on SchemaType", data.Path(), action)
		}
	}

	return nil, ErrNoSchemaTypeSet
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	schemaType := s.Type()
	k, err := SchemaTypeID(schemaType)
	if err != nil {
		return nil, err
	}
	var marshalled marshalSchemaType
	var schemaTypeJSON json.RawMessage

	marshalled.Type = &k

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

func (s *Schema) UnmarshalJSON(jsonSchema []byte) error {
	s.rawJSON = jsonSchema
	return nil
}

func (s *Schema) unmarshalRawJSON(mapper *SchemaUnmarshalMapper) (err error) {
	defer s.cleanRawJSON()
	var unmarshal marshalSchemaType

	err = json.Unmarshal(s.rawJSON, &unmarshal)

	if err != nil {
		return err
	}

	s.Metadata = unmarshal.Metadata

	if unmarshal.Type == nil {
		return fmt.Errorf("cannot unmarshal PongoSchema, no $type set in %v", s.rawJSON)
	}

	schemaType := mapper.Get(*unmarshal.Type)
	if schemaType == nil {
		return fmt.Errorf("cannot unmarshall SchemaType element: SchemaType ID %s not found in SchemaUnmarshalMapper", *unmarshal.Type)
	}

	if unmarshal.Body != nil {
		err = json.Unmarshal(*unmarshal.Body, schemaType)
		if err != nil {
			return fmt.Errorf("cannot unmarshall $body in %v: %w", s.rawJSON, err)
		}
	}

	err = s.SetType(schemaType)

	if err != nil {
		return err
	}

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

func (s *Schema) cleanRawJSON() {
	s.rawJSON = nil
}

func (s *Schema) Schema() *Schema {
	return s
}

func (s Schema) Children() (SchemaList, error) {
	var schemaType = s.Type()
	if schemaType == nil {
		return nil, ErrNoSchemaTypeSet
	}

	switch schemaTypeParent := schemaType.(type) {
	case ListSchema:
		return schemaTypeParent.Children(), nil
	case ObjectSchema:
		return Values[string, *Schema](schemaTypeParent.Children()), nil
	case ParentSchema:
		return SchemaList{schemaTypeParent.Child()}, nil
	}

	return SchemaList{}, nil
}

func (m SchemaMap) Children() SchemaMap {
	return m
}

func (l SchemaList) Children() SchemaList {
	return l
}

func (s Schema) GetMetadata(key string) (value string, ok bool) {
	return s.Metadata.Get(key)
}

func (s *Schema) SetMetadata(key string, value string) *Schema {
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
