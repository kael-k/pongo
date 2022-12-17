package tests

// this test file aims to expose all tests file in tests/assets
// and to validate the schemas inside it

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/xeipuuv/gojsonschema"

	"github.com/kael-k/pongo/pongo"
)

const testRootSchemas = "../tests/assets/schemas"
const testPongoSchemaFilename = "pongoschema.json"
const testJSONSchemaFilename = "jsonschema.json"

type testSchemaMarshall struct {
	dir        string
	wantSchema pongo.SchemaType
	typeMap    *pongo.SchemaUnmarshalMapper
}

func DecoratedSchemaUnmarshalMapper() *pongo.SchemaUnmarshalMapper {
	return pongo.DefaultSchemaUnmarshalMap().
		Set(func() pongo.SchemaType { return pongo.Decorate(pongo.Int()) })
}

// testsSchemaMarshall string should be the test dir name wrt testRootSchemas
var testsSchemaMarshall = map[string]testSchemaMarshall{
	"example-1": {
		"example-1",
		pongo.Object(pongo.O{
			"aString": pongo.String(),
			"aInt":    pongo.Int().SetMin(10).SetCastActions(pongo.SchemaActionParse),
			"aBytes":  pongo.Bytes().SetCast(true),
			"aDouble": pongo.Float64().SetMin(1).SetCast(true),
		}),
		pongo.DefaultSchemaUnmarshalMap(),
	},
	"example-2": {
		"example-2",
		pongo.Object(pongo.O{
			"aString": pongo.Schema(pongo.String()).SetMetadata("foo", "bar"),
			"aNestedObject": pongo.Object(pongo.O{
				"aString": pongo.String(),
				"aBool":   pongo.Bool(),
			}),
		}),
		pongo.DefaultSchemaUnmarshalMap(),
	},
	"example-3": {
		"example-3",
		pongo.Object(pongo.O{
			"aString": pongo.String(),
			"aList": pongo.List(pongo.AnyOf(
				pongo.Object(pongo.O{
					"aString": pongo.String().SetCast(true),
					"aBool":   pongo.Bool(),
				}),
				pongo.String(),
			)).SetMinLen(1).SetMaxLen(3),
		}).Require("aList"),
		pongo.DefaultSchemaUnmarshalMap(),
	},
	"example-4": {
		"example-4",
		pongo.Object(pongo.O{
			"aString":             pongo.String().SetCast(true).SetMinLen(3).SetMaxLen(5),
			"aDatetime":           pongo.Datetime().SetCast(true).SetBefore(time.Unix(1754038800, 0).UTC()).SetAfter(time.Unix(1596272400, 0).UTC()),
			"aUncastableDatetime": pongo.Datetime().SetAfter(time.Unix(1817110800, 0).UTC()),
		}).Require("aString", "aDatetime"),
		pongo.DefaultSchemaUnmarshalMap(),
	},
	"example-5": {
		"example-5",
		pongo.AllOf(
			pongo.Float64().SetMin(500).SetMax(500000),
			pongo.String().SetCast(true).SetMaxLen(5),
		).SetChain(true),
		pongo.DefaultSchemaUnmarshalMap(),
	},
	"example-1_decorator": {
		"example-1",
		pongo.Object(pongo.O{
			"aString": pongo.String(),
			"aInt":    pongo.Decorate(pongo.Int().SetMin(10).SetCastActions(pongo.SchemaActionParse)),
			"aBytes":  pongo.Bytes().SetCast(true),
			"aDouble": pongo.Float64().SetMin(1).SetCast(true),
		}),
		DecoratedSchemaUnmarshalMapper(),
	},
}

func (t testSchemaMarshall) TestDirPath() string {
	return fmt.Sprintf("%s/%s", testRootSchemas, t.dir)
}

func (t testSchemaMarshall) GetRawJSONSchema() ([]byte, error) {
	rawJSONSchema, err := os.ReadFile(fmt.Sprintf("%s/%s", t.TestDirPath(), testJSONSchemaFilename))
	if err != nil {
		return nil, err
	}
	return rawJSONSchema, nil
}

func (t testSchemaMarshall) GetJSONSchema() (gojsonschema.JSONLoader, error) {
	rawJSONSchema, err := t.GetRawJSONSchema()
	if err != nil {
		return nil, err
	}
	return gojsonschema.NewBytesLoader(rawJSONSchema), nil
}

func (t testSchemaMarshall) GetPongoSchema() (*pongo.SchemaNode, *pongo.Metadata, error) {
	rawPongoSchema, err := os.ReadFile(fmt.Sprintf("%s/%s", t.TestDirPath(), testPongoSchemaFilename))
	if err != nil {
		return nil, nil, err
	}
	schema, metadata, err := pongo.UnmarshalSchemaJSON(rawPongoSchema)
	if err != nil {
		return nil, nil, err
	}

	return schema, metadata, nil
}

func (t testSchemaMarshall) GetTestCaseFiles() (okTestsFiles []string, koTestsFiles []string, err error) {
	files, err := t.GetDirectoryFiles()
	if err != nil {
		return nil, nil, err
	}

	for _, file := range files {
		if file[:2] == "ok" && file[len(file)-5:] == ".json" {
			okTestsFiles = append(okTestsFiles, file)
		} else if file[:2] == "ko" && file[len(file)-5:] == ".json" {
			koTestsFiles = append(koTestsFiles, file)
		}
	}

	return
}

func (t testSchemaMarshall) GetDirectoryFiles() (filesName []string, err error) {
	files, err := os.ReadDir(t.TestDirPath())
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filesName = append(filesName, file.Name())
	}

	return filesName, nil
}

// validateDirectoryFiles checks if the file structure in the directory is valid
// it does not actually validate any file
func (t testSchemaMarshall) validateDirectoryFiles() error {
	okTestsFiles, koTestsFiles, err := t.GetTestCaseFiles()
	if err != nil {
		return err
	}

	files, err := t.GetDirectoryFiles()
	if err != nil {
		return err
	}

	var unknownFiles []string
	for _, file := range files {
		if file == testPongoSchemaFilename || file == testJSONSchemaFilename || file == "README.md" {
			continue
		}
		if InArray[string](file, okTestsFiles) || InArray[string](file, koTestsFiles) {
			continue
		}

		unknownFiles = append(unknownFiles, file)
	}

	if len(unknownFiles) > 0 {
		return fmt.Errorf("the following file are not known in test %s: %v", t.dir, unknownFiles)
	}

	return nil
}

func (t testSchemaMarshall) ValidateJSONSchemaMarshal() error {
	pongoSchema, _, err := t.GetPongoSchema()
	if err != nil {
		return err
	}
	testMarshalJSONSchema, err := pongo.MarshalJSONSchemaWithMetadata(pongoSchema, pongo.SchemaActionParse)
	if err != nil {
		if errors.Is(err, pongo.ErrSchemaNotJSONSchemaMarshalable) {
			return nil
		}
		return err
	}

	var wantMarshalJSONSchemaObj, testMarshalJSONSchemaObj map[string]interface{}
	wantMarshalJSONSchema, err := t.GetRawJSONSchema()
	if err != nil {
		return err
	}
	err = json.Unmarshal(wantMarshalJSONSchema, &wantMarshalJSONSchemaObj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(testMarshalJSONSchema, &testMarshalJSONSchemaObj)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(testMarshalJSONSchemaObj, wantMarshalJSONSchemaObj) {
		return fmt.Errorf(
			"cannot validate Pongo schema jsonschema marshaling for test %s/%s.\nExpected: %v\nGot: %v",
			t.TestDirPath(),
			t.dir,
			testMarshalJSONSchemaObj,
			wantMarshalJSONSchemaObj,
		)
	}

	return nil
}

func (t testSchemaMarshall) Validate() error {
	if err := t.validateDirectoryFiles(); err != nil {
		return err
	}

	pongoSchema, _, err := t.GetPongoSchema()
	if err != nil {
		return err
	}

	jsonSchema, err := t.GetJSONSchema()
	if err != nil {
		return err
	}

	okTestFiles, koTestFiles, err := t.GetTestCaseFiles()
	if err != nil {
		return err
	}

	err = t.ValidateJSONSchemaMarshal()
	if err != nil {
		return err
	}

	for _, okFile := range okTestFiles {
		var d interface{}
		rawJSON, err := os.ReadFile(fmt.Sprintf("%s/%s", t.TestDirPath(), okFile))
		if err != nil {
			return err
		}
		err = json.Unmarshal(rawJSON, &d)
		if err != nil {
			return err
		}

		_, schemaErr := pongoSchema.Parse(pongo.NewDataPointer(d, pongoSchema))
		if schemaErr != nil {
			return fmt.Errorf("cannot validate Pongo schema for ok test %s/%s: %w", t.dir, okFile, schemaErr)
		}

		res, err := gojsonschema.Validate(jsonSchema, gojsonschema.NewGoLoader(d))
		if err != nil {
			return fmt.Errorf("error validating json schema ok test %s/%s: %w", t.dir, okFile, err)
		}
		if !res.Valid() {
			return fmt.Errorf("cannot validate json schema ok test %s/%s: %v", t.dir, okFile, res.Errors())
		}
	}

	for _, koFile := range koTestFiles {
		var d interface{}
		rawJSON, err := os.ReadFile(fmt.Sprintf("%s/%s", t.TestDirPath(), koFile))
		if err != nil {
			return err
		}
		err = json.Unmarshal(rawJSON, &d)
		if err != nil {
			return err
		}

		_, schemaErr := pongoSchema.Parse(pongo.NewDataPointer(d, pongoSchema))
		if schemaErr == nil {
			return fmt.Errorf("cannot validate Pongo schema for ko test %s/%s: expected error, got no one", t.dir, koFile)
		}

		res, err := gojsonschema.Validate(jsonSchema, gojsonschema.NewGoLoader(d))
		if err == nil && res.Valid() {
			return fmt.Errorf("cannot validate JSON schema for ko test %s/%s: expected error, got no one", t.dir, koFile)
		}
	}

	return nil
}

func testAssetIntegrity(t *testing.T, testDir string, testCase testSchemaMarshall) {
	if testCase.dir != testDir {
		t.Errorf(
			"testsSchemaMarshall \"dir\" property and directory name mismatch (%s, %s)",
			testCase.dir,
			testDir,
		)
		return
	}

	if err := testCase.Validate(); err != nil {
		t.Errorf("testsSchemaMarshall validation failed: %s", err)
		return
	}
}

func TestAssetsIntegrity(t *testing.T) {
	testsDir, err := os.ReadDir(testRootSchemas)
	if err != nil {
		t.Errorf("cannot test schema schemas/unmarshall, error on read %s: %s", testRootSchemas, err)
		return
	}

	if len(testsDir) == 0 {
		t.Errorf("cannot find tests file in %s", testRootSchemas)
	}

	for _, testDir := range testsDir {
		dirName := testDir.Name()

		var testPath = fmt.Sprintf("%s/%s", testRootSchemas, dirName)

		// check if is a directory or the README.md file
		if _, err := os.ReadDir(testPath); err != nil {
			// check if is not README.md and the error is actually because it is a file
			pathErr, ok := err.(*fs.PathError)
			if ok && pathErr.Err == syscall.ENOTDIR && dirName == "README.md" {
				continue
			}

			t.Errorf("cannot read %s: %s", testPath, err)
			continue
		}

		testCase, ok := testsSchemaMarshall[dirName]
		if !ok {
			t.Errorf("testsSchemaMarshall %s not found", dirName)
			continue
		}
		testAssetIntegrity(t, dirName, testCase)
	}
}
