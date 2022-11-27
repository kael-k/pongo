package pongo

type Data interface{}

type DataPointer struct {
	data Data
	path Path
}

func NewDataPointer(data Data, schemaType SchemaType) *DataPointer {
	dp := &DataPointer{
		data: data,
	}

	if data != nil {
		dp.path = *NewPath(*NewPathElement("", data, schemaType))
	} else {
		dp.path = *NewPath()
	}

	return dp
}

func (d DataPointer) Push(key string, data Data, schemaType SchemaType) *DataPointer {
	d.path = *d.path.Push(key, data, schemaType)
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
	return d.data
}

func (d DataPointer) Clone() *DataPointer {
	return &DataPointer{
		data: d.data,
		path: *d.path.Clone(),
	}
}
