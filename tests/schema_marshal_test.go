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
		schema, metadata, err := pongo.UnmarshalSchemaWithMapper(jsonSchema, schemaType.typeMap)
		if err != nil {
			t.Errorf("error test schema unmarshall %s, error on unmarshall JSON: %s", testID, err)
			continue
		}
		if !reflect.DeepEqual(schema, pongo.Schema(schemaType.wantSchema)) {
			t.Errorf("error test schema unmarshall %s, unmarshalled schema does not match the wanted one", testID)
			continue
		}

		// test marshalling
		testJSONSchema, err := pongo.MarshalPongoSchemaWithMetadata(schemaType.wantSchema, metadata)
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

func (t TestDummySchemaType) Process(_ pongo.SchemaAction, _ *pongo.DataPointer) (parsedData pongo.Data, err error) {
	return nil, fmt.Errorf("not implemented")
}

func TestSchemaTypeID(t *testing.T) {
	test := pongo.SchemaTypeID(&TestDummySchemaType{})
	if test != "tests.TestDummySchemaType" {
		t.Errorf("error test SchemaTypeID() with testing.TestDummySchemaType, expected id $testing.TestDummySchemaType, got %s", test)
	}

	test = pongo.SchemaTypeID(pongo.String())
	if test != "string" {
		t.Errorf("error test SchemaTypeID() with testing.TestDummySchemaType, expected id $testing.TestDummySchemaType, got %s", test)
	}
}
