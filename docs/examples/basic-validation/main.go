package main

import (
	"fmt"
	"github.com/kael-k/pongo/pongo"
)

func main() {
	var schema pongo.SchemaType
	var data pongo.Data
	var err error

	schema = pongo.String()

	// expected result: success
	data = "foo"
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}

	// expected result: error
	data = 123
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}
}
