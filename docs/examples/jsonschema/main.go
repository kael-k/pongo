package main

import (
	"fmt"
	"github.com/kael-k/pongo/v2/pongo"
)

func main() {
	var schemaType pongo.SchemaType
	var err error

	schemaType = pongo.OneOf(pongo.String())
	jsonSchema, err := pongo.MarshalJSONSchemaWithMetadata(pongo.Schema(schemaType), pongo.SchemaActionSerialize)

	// expected result: success
	if err != nil {
		fmt.Printf("an error occurred during jsonschema marshaling: %s", err)
	} else {
		fmt.Printf("jsonschema-marshaled schema: %s\n", jsonSchema)
	}
}
