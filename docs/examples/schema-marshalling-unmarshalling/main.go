package main

import (
	"fmt"
	"github.com/kael-k/pongo/v2/pongo"
	"reflect"
)

func main() {
	var schema pongo.SchemaType
	var err error

	schema = pongo.Object(pongo.O{
		"aInt":    pongo.Int(),
		"aString": pongo.String(),
	})
	marshalledSchema, err := pongo.MarshalPongoSchema(schema)
	if err != nil {
		panic("unexpected error during marshalling: " + err.Error())
	} else {
		fmt.Printf("marshalled schema:\n%s\n", marshalledSchema)
	}

	newSchema, _, err := pongo.UnmarshalPongoSchema(marshalledSchema)
	if err != nil {
		panic("unexpected error during unmarshalling: " + err.Error())
	} else {
		fmt.Printf("unmarshalled schema: type: %s, schemaType\n", reflect.TypeOf(newSchema).String())
	}

	if reflect.DeepEqual(schema, pongo.Schema(newSchema).Type()) {
		fmt.Printf("original schema and the unmarshalled one match")
	} else {
		panic("original schema and the unmarshalled one do not match")
	}
}
