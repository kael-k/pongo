package pongo

import (
	"encoding/json"
	"fmt"
)

/* These types are used to wrap SchemaType instances.
The wrappers function is to support the schema marshal/unmarshal process:
In SchemaType implementation, if is requested any SchemaType nesting, the implementation must:
* store the nested schema using one of SchemaNode (for 0..1 SchemaType), SchemaList (for SchemaType list) or SchemaMap (for Maps)
  this ensures that the marshaling is automatically done with SchemaNode.MarshalJSON.
  In this way, the SchemaType implementation is not responsible for the marshaling process
* if a SchemaType has any child, it MUST implement also ObjectSchema (if it has SchemaMap children),
  ParentSchema (for SchemaList children) or ParentSchema (if it has a single SchemaNode child). The implementation
  is required for the unmarshalling process to recursively pass the SchemaUnmarshalMapper and resolve the correct
  SchemaType to unmarshal
*/

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

func (s *SchemaNode) unmarshalRawJSON(mapper *SchemaUnmarshalMapper) (err error) {
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
		return fmt.Errorf("cannot unmarshall SchemaType element: SchemaType ID %s not found in SchemaUnmarshalMapper", *unmarshal.Type)
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
