package main

import (
	"fmt"
	"github.com/kael-k/pongo/v2/pongo"
)

func main() {
	var schema pongo.SchemaType
	var data pongo.Data
	var err error

	// expected result: success
	schema = pongo.List(pongo.String())
	data = []interface{}{"1234", "aString"}
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}

	// expected result: success
	schema = pongo.Object(pongo.O{
		"foo": pongo.String(),
	})
	data = map[string]interface{}{
		"foo": "bar",
	}
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}
}
