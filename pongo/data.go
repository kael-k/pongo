package pongo

// Data represent a generic input root for the PonGO Schema
type Data interface{}

// DataPointer is an internal structure used by PonGO SchemaNode
// to track the root structure across the validation
// * root contains always the data structure as passed originally by the caller
// * path contains the Path which contains the target SchemaNode has to process.
//   - It also contains the stack all the Data passed across the current stack of SchemaNode
//
// Path works as a stack, which is needed to retrieve correct data to process in nested validations.
// For example, take ObjectType: which for every element in a map[string]Data, it validates every element
// with the SchemaType associated with the map key. The Data at the key of map[string]Data is then pushed
// in the Path with Push and then the new DataPointer is passed to the child SchemaType
type DataPointer struct {
	root Data
	path Path
}

// NewDataPointer construct a DataPointer
func NewDataPointer(schemaNode *SchemaNode, data Data) *DataPointer {
	dp := &DataPointer{
		root: data,
	}

	if data != nil {
		dp.path = *NewPath(*NewPathElement(schemaNode, data, ""))
	} else {
		dp.path = *NewPath()
	}

	return dp
}

// Push a new entry in the DataPointer Path stack
func (d DataPointer) Push(schemaNode *SchemaNode, data Data, key string) *DataPointer {
	d.path = *d.path.Push(schemaNode, data, key)
	return &d
}

func (d DataPointer) Path() Path {
	return d.path
}

func (d *DataPointer) Get() Data {
	if d == nil {
		return nil
	}

	return d.path.Value()
}

func (d DataPointer) GetRoot() Data {
	return d.root
}

func (d DataPointer) Clone() *DataPointer {
	return &DataPointer{
		root: d.root,
		path: *d.path.Clone(),
	}
}
