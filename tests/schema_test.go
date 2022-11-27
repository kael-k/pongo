package tests

// this file only defines test types and function helpers for all other type_*_test.go files

import (
	"reflect"
	"testing"

	"github.com/kael-k/pongo/pongo"
)

type DataFactory func() pongo.Data

type testSchemaCase struct {
	desc   string
	schema pongo.SchemaType
	data   DataFactory
	want   DataFactory
	errors int
}

func testSchemaCaseParse(testCases []testSchemaCase) func(t *testing.T) {
	return func(t *testing.T) {
		for _, testCase := range testCases {
			data := testCase.data()
			p, err := pongo.Parse(testCase.schema, data)
			testSchemaCheckTest(t, testCase, data, p, err)
		}
	}
}

func testSchemaCaseSerialize(testCases []testSchemaCase) func(t *testing.T) {
	return func(t *testing.T) {
		for _, testCase := range testCases {
			data := testCase.data()
			p, err := pongo.Serialize(testCase.schema, data)
			testSchemaCheckTest(t, testCase, data, p, err)
		}
	}
}

func testSchemaCaseProcess(testCases []testSchemaCase, action pongo.SchemaAction) func(t *testing.T) {
	return func(t *testing.T) {
		for _, testCase := range testCases {
			data := testCase.data()
			p, err := pongo.Process(testCase.schema, action, data)
			testSchemaCheckTest(t, testCase, data, p, err)
		}
	}
}

func testSchemaCheckTest(t *testing.T, testCase testSchemaCase, testData pongo.Data, testValue pongo.Data, err error) {
	schemaErr, ok := err.(*pongo.SchemaError)
	if !ok && err != nil {
		t.Errorf("test schema %s: got errors that cannot be asserted to *SchemaError", testCase.desc)
	}

	if schemaErr != nil {
		if len(schemaErr.Errors) != testCase.errors {
			t.Errorf("test schema %s: expected %d erorr(s), got %d error(s): %s", testCase.desc, testCase.errors, len(schemaErr.Errors), schemaErr)
		}
	} else if schemaErr == nil && testCase.errors != 0 {
		t.Errorf("test schema %s expected error, but err == nil", testCase.desc)
		return
	}

	// test parsed data are the expected one
	if testCase.errors == 0 {
		want := testCase.want()
		if !reflect.DeepEqual(testValue, want) {
			t.Errorf("test schema %s expected validation %s, got %s", testCase.desc, want, testValue)
		}
	}

	// test that data has not changed
	clone := testCase.data()
	if !reflect.DeepEqual(testData, clone) {
		t.Errorf("test schema %s failed clone check, expected %s, got %s", testCase.desc, clone, testData)
	}
}

func TestSchemaMetadata(t *testing.T) {
	s := pongo.NewEmptySchema()
	if v, ok := s.GetMetadata("foo"); ok {
		t.Errorf("expected ok == false on new Schema GetMetadata(\"foo\"), got [%v, %v]", ok, v)
	}
	s.SetMetadata("foo", "bar")

	if v, ok := s.GetMetadata("foo"); !ok || v != "bar" {
		t.Errorf("expected ok == true and v == \"bar\" on Schema GetMetadata(\"foo\"), got [%v, %v]", ok, v)
	}

	s.SetMetadata("foo", "baz")

	if v, ok := s.GetMetadata("foo"); !ok || v != "baz" {
		t.Errorf("expected ok == true and v == \"baz\" on Schema GetMetadata(\"foo\"), got [%v, %v]", ok, v)
	}
}
