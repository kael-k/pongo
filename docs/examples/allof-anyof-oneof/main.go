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
	schema = pongo.AnyOf(
		pongo.Int(),
		pongo.String(),
	)
	data = "foo"
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}

	data = 123
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}

	// expected result: error, because the SchemaType are not chained
	schema = pongo.AllOf(
		pongo.Int(),
		pongo.String().SetCast(true),
		pongo.String().SetMinLen(2),
	)
	data = 123
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s\n", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}

	// expected result: success
	schema = pongo.AllOf(
		pongo.Int(),
		pongo.String().SetCast(true),
		pongo.String().SetMinLen(2),
	).SetChain(true)
	data = 123
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v\n", data)
	}

}
