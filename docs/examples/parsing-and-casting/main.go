package main

import (
	"fmt"
	"github.com/kael-k/pongo/v2/pongo"
	"reflect"
)

func main() {
	var schema pongo.SchemaType
	var data pongo.Data
	var err error

	// expected result: success
	schema = pongo.String().SetCast(true)
	data = 1234
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v, type: %s\n", data, reflect.TypeOf(data).String())
	}

	// expected result: success
	schema = pongo.Datetime().SetCast(true)
	data = "2022-08-09T00:00:00Z"
	data, err = pongo.Parse(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema validate successfully, %v, type: %s\n", data, reflect.TypeOf(data).String())
	}
}
