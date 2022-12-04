# PonGO Schema

![logo](https://repository-images.githubusercontent.com/574236489/645c9b8c-ecce-44c2-a82e-1351a992580e)
*logo credits to [@herbrandh](https://github.com/herbrandh)*
## Introduction

PonGO Schema is a library which provide rich and powerful validation, parsing  and serializing data.

These are some library features and use case:

* the library can be used for data serialization/deserialization and parsing
* data types set are extensible, you just need to implements `Schema`
* a schema is dynamically defined with objects, in contrast with a static schema defined by struct.
* marshall/unmarshall the schema in JSON

The library is expired by the Python [Schema library](https://pypi.org/project/schema/),
plus the option to serialize and deserialize the schema (like an "enrich" `jsonschema`) and the data (for example to
convert some "library specific" data types, such `time.Time`, into a more JSON-friendly `string`)

## Usage

### Basic validation

The library provides some `SchemaType` that parse the main go data types.
In the example below, there is an example that with a "String" as a schema and "foo" as input data.

```go
var schema pongo.SchemaType
var data pongo.Data
var err error

schema = pongo.String()

// expected result: success
data = "foo"
data, err = pongo.Parse(schema, data)
if err != nil {
    fmt.Printf("the schema does not validate for reasons: %s", err)
} else {
    fmt.Printf("the schema validate successfully, %v\n", data)
}
```

The output will show if the `Parse` process had success and what is the output of the process:

```
the schema validate successfully, foo
```

Changing the data with something that is not a `string` will produce an error

```go
data = 123
```

Output:

```
the schema does not validate for reasons: the schema encountered the followed error(s) = [path: .<string>, pathData: 123, error: schema does not validate: .<string> is not a "String"]
```

### Nested structures validation

Obviously, for more complex data structure we need more complex schemas such as a slices of string:

```go
data = []interface{}{"1234", "aString"}
```

To validate the example we need to use the `Schema` `ListType`:

```go
schema = pongo.List(pongo.String())
```

For `map[string]inteface`, we can use the `ObjectType` schema:

```go
schema = pongo.Object(pongo.O{
    "foo": pongo.String(),
})
data = map[string]interface{}{
    "foo": "bar",
}
```

Output:

```
the schema validate successfully, [1234 aString]
the schema validate successfully, map[foo:bar]
```

### Parsing and casting

When PonGO Schema `.Parse` function is called, 2 values are returned: a `pongo.Data` and `error`.
The first return element is the result of the parsing process of input data, but since the previous schema did not
manipulate data,
input and output were the same.

However, the default `SchemaType` can actually elaborate the data during the `Parse` process if configured correctly:

Depending on the schema, there are some options that can be set to change the behaviour of the schema processing.
For example, we cannot validate an `int` with `pongo.String` by default, but we can enable the casting option
for `StringType`:

```go
schema = pongo.String().SetCast(true)
data = 1234
data, err = pongo.Parse(schema, data)
if err != nil {
    fmt.Printf("the schema does not validate for reasons: %s", err)
} else {
    fmt.Printf("the schema validate successfully, %v, type: %s\n", data, reflect.TypeOf(data).String())
}
```

Another good example is `Datetime` parsing from RFC3339 string to `time.Time`

```go
schema = pongo.Datetime().SetCast(true)
data = "2022-08-09T00:00:00Z"
data, err = pongo.Parse(schema, data)
if err != nil {
    fmt.Printf("the schema does not validate for reasons: %s", err)
} else {
    fmt.Printf("the schema validate successfully, %v, type: %s\n", data, reflect.TypeOf(data).String())
}
```

Output:

```
the schema validate successfully, 1234, type: string
the schema validate successfully, 2022-08-09 00:00:00 +0000 UTC, type: time.Time
```

Please refer to the docs for the `SchemaType` implementations for all options.

### Serialization

The `Serialize` process is the reverse process of `Parse` process.

```go
schema = pongo.Datetime().SetCast(true)
data = time.Unix(1668607500, 00).UTC()
data, err = pongo.Serialize(schema, data)
if err != nil {
    fmt.Printf("the schema does not validate for reasons: %s", err)
} else {
    fmt.Printf("the schema serialized successfully, %v, type: %s\n", data, reflect.TypeOf(data).String())
}
```

Output:

```
the schema serialized successfully, 2022-11-16T14:05:00Z, type: string
```

### `AllOf`, `AnyOf` and `OneOf`

Multiple `SchemaType` can process the same type of data with different logic.

* `AllOf` process the data with all `SchemaType`, it will fail if one `SchemaType` return `error`
* `AnyOf` process the data with at least one `SchemaType`, the fist `SchemaType` running the process request with no
  error
  will return, it will fail if no `SchemaType` return `error`
* `OneOf` process the data with all `SchemaType` but only one `SchemaType` must return no error

An example

```go
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
```

Output:

```
the schema validate successfully, foo
the schema validate successfully, 123
```

### AllOf `SchemaType` chaining

By default, `AllOf` will return only the last `SchemaType` output, and every `SchemaType` input is the same input of `AllOf`.

`AllOf` has an option to pass the output of a `SchemaType` as an input to the next `SchemaType` to call. This option is called "chaining"

For example, this validation fails because the third `SchemaType` has as input `123` but the `String` has no casting enabled.
```go
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
```

Output:
```
the schema does not validate for reasons: the schema encountered the followed error(s) = [path: .<allOf>, pathData: 123, error: schema does not validate: .<allOf> is not a "String"]
```

But the following example will work, since the input of the third parameter will be `"123"` (cast as a string, from the second `SchemaType`)
```go
schema = pongo.AllOf(
    pongo.Int(),
    pongo.String().SetCast(true),
    pongo.String().SetMinLen(2),
).SetChain(true)
data = 123
data, err = pongo.Parse(schema, data)
if err != nil {
    fmt.Printf("the schema does not validate for reasons: %s\n", err)
} else {
    fmt.Printf("the schema validate successfully, %v\n", data)
}
```

Output:

```
the schema validate successfully, 123
```

### PonGO Schema Marshalling and Unmarshalling
A PonGO Schema instance is JSON-serializable and unserializable

```go
schema = pongo.Object(pongo.O{
    "aInt":    pongo.Int(),
    "aString": pongo.String(),
})
marshalledSchema, err := pongo.MarshalSchemaJSON(schema)
if err != nil {
    panic("unexpected error during marshalling: " + err.Error())
} else {
    fmt.Printf("marshalled schema:\n%s\n", marshalledSchema)
}
```

Output:
```
marshalled schema:
{"$body":{"$body":{"properties":{"aInt":{"$type":"int"},"aString":{"$type":"string"}}},"$type":"object"},"$version":"1.0"}
```

To unmarshall back the schema

```go
newSchema, _, err := pongo.UnmarshalSchemaJSON(marshalledSchema)
if err != nil {
    panic("unexpected error during unmarshalling: " + err.Error())
} else {
    fmt.Printf("unmarshalled schema: type: %s, schemaType\n", reflect.TypeOf(newSchema).String())
}

if reflect.DeepEqual(schema, pongo.Schema(newSchema).Type()) {
    fmt.Printf("original schema and the unmarshalled one match")
} else {
    panic("original schema and the unmarshalled one do not match")
}
```

Output:
```
unmarshalled schema: type: *pongo.Schema, schemaType
original schema and the unmarshalled one match
```

### Implementing a `SchemaType`

In order to create a valid "type" for the pongo library, you need to implement a `SchemaType`.

If your schema contains nested `SchemaType`(s) (such as `ListType`, `ObjectType`, `AllOfType`) you must implement `ParentSchema` 
which return a list of the nested `SchemaNode`(s). This implementation is required to correctly marshal/unmarshal the `Schema`(s).

If you need to unmarshal your custom type, you'll also need to use a custom `SchemaUnmarshalMapper` which contains your custom type.
When you'll unmarshal your PonGO schema, you must pass to the unmarshal function your custom map

```go
myCustomMapper = pongo.DefaultSchemaUnmarshalMap().set()
myPonGOSchema, err := pongo.UnmarshalSchemaJSONWithMapper(myPonGOSchemaJSON, myCustomMapper)
```