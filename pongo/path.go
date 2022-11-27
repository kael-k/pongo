package pongo

import (
	"errors"
	"fmt"
)

type PathElement struct {
	key         string
	data        Data
	override    Data
	hasOverride bool
	schemaType  SchemaType
}

func NewPathElement(key string, data Data, schemaType SchemaType) *PathElement {
	return &PathElement{
		key:         key,
		data:        data,
		override:    nil,
		hasOverride: false,
		schemaType:  schemaType,
	}
}

func (e PathElement) Key() any {
	return e.key
}

func (e PathElement) Data() Data {
	return e.data
}

func (e PathElement) Schema() any {
	return e.schemaType
}

func (e *PathElement) SetOverride(override Data) {
	e.override = override
	e.hasOverride = true
}

func (e *PathElement) UnsetOverride() {
	e.override = nil
	e.hasOverride = false
}

const PathSeparator = "."

type Path struct {
	elements []PathElement
}

func NewPath(keys ...PathElement) *Path {
	if len(keys) > 0 {
		return &Path{keys}
	}
	return &Path{elements: []PathElement{}}
}

func (path Path) Elements() []PathElement {
	return path.elements
}

func (path Path) String() string {
	var stringPath = ""

	for _, pathElement := range path.elements {
		schemaTypeID, err := SchemaTypeID(pathElement.schemaType)
		if err != nil {
			schemaTypeID = "ErrUnknownType"
		}
		stringPath += fmt.Sprintf("%s%s<%s>", PathSeparator, pathElement.key, schemaTypeID)
	}

	return stringPath
}

func (path Path) Value() Data {
	last := path.Last()
	if last == nil {
		return nil
	}
	if last.hasOverride {
		return last.override
	}
	return last.data
}

func (path Path) OriginalValue() Data {
	last := path.Last()
	if last == nil {
		return nil
	}
	return last.data
}

func (path Path) OverwrittenValue() (Data, bool) {
	last := path.Last()
	if last == nil {
		return nil, false
	}
	return last.override, last.hasOverride
}

func (path Path) SetOverride(override Data) error {
	last := path.Last()
	if last == nil {
		return errors.New("cannot set override in Path, the path has no element")
	}
	last.SetOverride(override)
	return nil
}

func (path Path) UnsetOverride() error {
	last := path.Last()
	if last == nil {
		return errors.New("cannot unset override in Path, the path has no element")
	}
	last.UnsetOverride()
	return nil
}

func (path Path) Size() int {
	return len(path.elements)
}

// Last return the last PathElement in the Path
func (path Path) Last() *PathElement {
	if s := path.Size(); s > 0 {
		return &path.elements[s-1]
	}
	return nil
}

// Push a new PathElement in Path
func (path Path) Push(key string, data Data, schemaType SchemaType) *Path {
	path.elements = append(path.elements, *NewPathElement(key, data, schemaType))

	return &path
}

// Clone a new PathElement in Path
func (path Path) Clone() *Path {
	return &Path{
		elements: path.elements,
	}
}
