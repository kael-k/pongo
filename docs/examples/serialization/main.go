package main

import (
	"fmt"
	"github.com/kael-k/pongo/v2/pongo"
	"reflect"
	"time"
)

func main() {
	var schema pongo.SchemaType
	var data pongo.Data
	var err error

	// expected result: success
	schema = pongo.Datetime().SetCast(true)
	data = time.Unix(1668607500, 00).UTC()
	data, err = pongo.Serialize(schema, data)
	if err != nil {
		fmt.Printf("the schema does not validate for reasons: %s", err)
	} else {
		fmt.Printf("the schema serialized successfully, %v, type: %s\n", data, reflect.TypeOf(data).String())
	}
}
