package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

func TestDecorator(t *testing.T) {
	var decoratedString = pongo.Decorate(pongo.String())
	var handlerFn = func(_ pongo.SchemaType, _ pongo.SchemaAction, _ *pongo.DataPointer) (data pongo.Data, err error) {
		return "bar", nil
	}
	decoratedString.SetHandlers(handlerFn, pongo.SchemaActionParse)

	schemaTypeID := pongo.SchemaTypeID(decoratedString)
	if schemaTypeID != "string" {
		t.Errorf("expected SchemaTypeID for decorated StringType == \"string\", got %s", schemaTypeID)
	}

	r, err := pongo.Process(decoratedString, pongo.SchemaActionParse, "foo")
	if err != nil {
		t.Errorf("error parsing decorated StringType: %s", err)
	}
	if r != "bar" {
		t.Errorf("expected parsed data == \"bar\", got: %s", r)
	}

	r, err = pongo.Process(decoratedString, pongo.SchemaActionSerialize, "foo")
	if err != nil {
		t.Errorf("error serialize decorated StringType: %s", err)
	}
	if r != "foo" {
		t.Errorf("expected serialize data == \"foo\", got: %s", r)
	}

	decoratedString.SetDefaultHandler(handlerFn)

	r, err = pongo.Process(decoratedString, pongo.SchemaActionSerialize, "bar")
	if err != nil {
		t.Errorf("error serialize decorated StringType: %s", err)
	}
	if r != "bar" {
		t.Errorf("expected serialize data == \"bar\", got: %s", r)
	}

	decoratedString.UnsetHandlers(pongo.SchemaActionSerialize, pongo.SchemaActionParse)
	decoratedString.UnsetDefaultHandler()

	r, err = pongo.Process(decoratedString, pongo.SchemaActionParse, "foo")
	if err != nil {
		t.Errorf("error parsing decorated StringType: %s", err)
	}
	if r != "foo" {
		t.Errorf("expected parsed data == \"foo\", got: %s", r)
	}
}
