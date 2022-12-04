package pongo

import (
	"encoding/json"
	"fmt"
)

type DecoratorFn func(originalType SchemaType, action SchemaAction, dataPointer *DataPointer) (data Data, err error)

type DecoratedType struct {
	OriginalType SchemaType

	handlersMap    map[SchemaAction]DecoratorFn
	defaultHandler DecoratorFn
}

func (d *DecoratedType) UnmarshalJSON(bytes []byte) (err error) {
	if d == nil {
		return fmt.Errorf("cannot UnmarshalJSON a DecoratedType with no SchemaType set")
	}
	return json.Unmarshal(bytes, d.OriginalType)
}

func (d DecoratedType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.OriginalType)
}

func Decorate(schemaType SchemaType) *DecoratedType {
	return &DecoratedType{
		OriginalType:   schemaType,
		handlersMap:    map[SchemaAction]DecoratorFn{},
		defaultHandler: nil,
	}
}

func (d *DecoratedType) SetDefaultHandler(handler DecoratorFn) *DecoratedType {
	d.defaultHandler = handler
	return d
}

func (d *DecoratedType) UnsetDefaultHandler() {
	d.defaultHandler = nil
}

func (d *DecoratedType) SetHandlers(handler DecoratorFn, actions ...SchemaAction) *DecoratedType {
	for _, action := range actions {
		d.handlersMap[action] = handler
	}
	return d
}

func (d *DecoratedType) UnsetHandlers(actions ...SchemaAction) {
	for _, action := range actions {
		delete(d.handlersMap, action)
	}
}

func (d *DecoratedType) SchemaTypeID() string {
	if d != nil {
		return SchemaTypeID(d.OriginalType)
	}

	// this should return an error
	return SchemaTypeID(nil)
}

func (d DecoratedType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	if d.handlersMap != nil {
		if handlerMap, ok := d.handlersMap[action]; ok {
			return handlerMap(d.OriginalType, action, dataPointer)
		}
	}
	if d.defaultHandler != nil {
		return d.defaultHandler(d.OriginalType, action, dataPointer)
	}

	return Schema(d.OriginalType).Process(action, dataPointer)
}
