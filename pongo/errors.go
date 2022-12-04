package pongo

import (
	"fmt"
	"strings"
)

var ErrNoSchemaTypeSet = fmt.Errorf("this SchemaNode has no valid SchemaType set")

func ErrInvalidAction(schemaType SchemaType, action SchemaAction) error {
	return fmt.Errorf("cannot execute schema action %s on schema type %v", action, schemaType)
}

type SchemaError struct {
	Errors []SchemaElementError
}

// NewSchemaError create a new empty *SchemaError
func NewSchemaError() *SchemaError {
	return &SchemaError{}
}

// NewSchemaErrorWithError create a new *SchemaError and append an error
func NewSchemaErrorWithError(path Path, err error) *SchemaError {
	return (&SchemaError{}).Append(path, err)
}

// NewSchemaWithCasting create a new *SchemaError
// if err is a *SchemaError, it will return the (*SchemaError)(err)
// else it will return NewSchemaWithError(path, err)
func NewSchemaWithCasting(path Path, err error) *SchemaError {
	schemaErr, ok := err.(*SchemaError)
	if ok {
		if schemaErr != nil {
			return &SchemaError{
				Errors: schemaErr.Errors,
			}
		}
		return NewSchemaError()
	}
	return NewSchemaErrorWithError(path, err)
}

type SchemaElementError struct {
	path Path
	err  error
}

func (s SchemaElementError) Path() Path {
	return s.path
}

func (s SchemaElementError) Error() error {
	return s.err
}

func (s SchemaError) Error() string {
	var errors []string
	for _, v := range s.Errors {
		errors = append(errors, fmt.Sprintf("path: %s, pathData: %#v, error: %s", v.path, v.path.Value(), v.err))
	}

	return fmt.Sprintf("the schema encountered the followed error(s) = [%s]", strings.Join(errors, "; "))
}

func (s SchemaError) Append(path Path, err error) *SchemaError {
	s.Errors = append(s.Errors, SchemaElementError{
		path: path,
		err:  err,
	})

	return &s
}

func (s SchemaError) Merge(s2 *SchemaError) *SchemaError {
	s.Errors = append(s.Errors, s2.Errors...)
	return &s
}

// MergeWithCast will Merge err if is a *SchemaError
// or cast it to a new *SchemaError and then merge it.
func (s SchemaError) MergeWithCast(path Path, err error) *SchemaError {
	return s.Merge(NewSchemaWithCasting(path, err))
}
