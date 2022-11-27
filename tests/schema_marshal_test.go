package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/kael-k/pongo/pongo"
)

func TestSchemasMarshallUnmarshall(t *testing.T) {
	for testID, schemaType := range testsSchemaMarshall {
		testDir := schemaType.dir
		jsonSchemaPath := fmt.Sprintf("%s/%s/%s", testRootSchemas, testDir, testPongoSchemaFilename)
		jsonSchema, err := os.ReadFile(jsonSchemaPath)

		if err != nil {
			t.Errorf("error test schema schemas/unmarshall %s, error on read %s: %s", testID, jsonSchemaPath, err)
			continue
		}

		// test unmarshalling
		schema, metadata, err := pongo.UnmarshalSchemaJSONWithMapper(jsonSchema, schemaType.typeMap)
		if err != nil {
			t.Errorf("error test schema unmarshall %s, error on unmarshall JSON: %s", testID, err)
			continue
		}
		if !reflect.DeepEqual(schema, schemaType.wantSchema.Schema()) {
			t.Errorf("error test schema unmarshall %s, unmarshalled schema does not match the wanted one", testID)
			continue
		}

		// test marshalling
		testJSONSchema, err := pongo.MarshalSchemaJSONWithMetadata(schemaType.wantSchema, metadata)
		if err != nil {
			t.Errorf("error test schema schemas %s, error on schemas BaseSchemaType: %s", testID, err)
			continue
		}

		var testRawJSONSchema, wantRawJSONSchema map[string]interface{}
		err = json.Unmarshal(testJSONSchema, &testRawJSONSchema)
		if err != nil {
			t.Errorf("error test schema schemas %s, error on unmarshall serialized json BaseSchemaType to native go types: %s", testID, err)
			continue
		}
		err = json.Unmarshal(jsonSchema, &wantRawJSONSchema)
		if err != nil {
			t.Errorf("error test schema schemas %s, error on unmarshall test json file to native go types: %s", testID, err)
			continue
		}

		if !reflect.DeepEqual(testRawJSONSchema, wantRawJSONSchema) {
			t.Errorf("error test schema schemas %s, marshalled schema does not match the original one", testID)
			continue
		}
	}
}

type TestDummySchemaType struct{}

func (t TestDummySchemaType) Schema() *pongo.Schema {
	return pongo.NewBaseSchema(t)
}

func (t TestDummySchemaType) Validate(_ *pongo.DataPointer) (err error) {
	return fmt.Errorf("not implemented")
}

func (t TestDummySchemaType) Parse(_ *pongo.DataPointer) (parsedData pongo.Data, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (t TestDummySchemaType) Serialize(_ *pongo.DataPointer) (serializedData pongo.Data, err error) {
	return nil, fmt.Errorf("not implemented")
}

func TestSchemaTypeID(t *testing.T) {
	test, err := pongo.SchemaTypeID(&TestDummySchemaType{})
	if test != "tests.TestDummySchemaType" {
		t.Errorf("error test SchemaTypeID() with testing.TestDummySchemaType, expected id $testing.TestDummySchemaType, got %s", test)
	}
	if err != nil {
		t.Errorf("error test SchemaTypeID() with testing.TestDummySchemaType, unexpected error: %s", err)
	}

	test, err = pongo.SchemaTypeID(pongo.String())
	if test != "string" {
		t.Errorf("error test SchemaTypeID() with testing.TestDummySchemaType, expected id $testing.TestDummySchemaType, got %s", test)
	}
	if err != nil {
		t.Errorf("error test SchemaTypeID() with testing.TestDummySchemaType, unexpected error: %s", err)
	}
}
