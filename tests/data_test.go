package tests

import (
	"testing"

	"github.com/kael-k/pongo/pongo"
)

func TestData(t *testing.T) {
	const TestValue1 = 42
	const TestValue2 = "aRandomValue"

	var d *pongo.DataPointer
	if v := d.Get(); v != nil {
		t.Errorf("error TestData: expected nil on get on empty DataPointer, got %v", v)
	}

	s := pongo.String()
	d = pongo.NewDataPointer(TestValue1, pongo.Schema(s))
	d = d.Push("1", TestValue2, pongo.Schema(s))

	path := d.Path()
	if v := path.Size(); v != 2 {
		t.Errorf("error TestData: expected path of size 2, got %v", v)
	}

	if v := d.Get(); v != TestValue2 {
		t.Errorf("error TestData: expected d.Get return value %s, got %v", TestValue2, v)
	}

	if v := d.GetRoot(); v != TestValue1 {
		t.Errorf("error TestData: expected d.Get return value %d, got %v", TestValue1, v)
	}

	// testing clone
	clone := d.Clone()
	if clone == d {
		t.Errorf("error TestData: d.Clone returned the address of d instead a new one (so it didn't actualy cloned cloned)")
	}

	if clone.GetRoot() != d.GetRoot() {
		t.Errorf("error TestData: cloned datapointer %v is different from the original one %v", clone.GetRoot(), d.GetRoot())
	}

	clonedPath := d.Path()
	if v := clonedPath.Size(); v != 2 {
		t.Errorf("error TestData: cloned datapointer path expected a length of 2, got %d", v)
	}

	elements, clonedElements := path.Elements(), clonedPath.Elements()
	for i, v := range elements {
		if v.Key() != clonedElements[i].Key() {
			t.Errorf("error TestData: original key %v and cloned key %v mismatch", v.Key(), clonedElements[i].Key())
		}

		if v.Data() != clonedElements[i].Data() {
			t.Errorf("error TestData: original data %v and cloned data %v mismatch", v.Data(), clonedElements[i].Data())
		}

		if v != clonedElements[i] {
			t.Errorf("error TestData: original schema %v and cloned schema %v mismatch", v, clonedElements[i])
		}
	}

	// testing variable override
	d = d.Push("foo", "bar", pongo.Schema(pongo.String()))
	d2 := d.Push("foo", "baz", pongo.Schema(pongo.String()))
	if v := d2.Get(); v != "baz" {
		t.Errorf("expected value baz, got: %v", v)
	}

	err := d2.Path().SetOverride("bat")
	if err != nil {
		t.Errorf("expected no error in SetOverride, got: %s", err)
	}
	if v := d2.Get(); v != "bat" {
		t.Errorf("expected value bat, got: %v", v)
	}
	if v, ok := d2.Path().OverwrittenValue(); !ok || v != "bat" {
		t.Errorf("expected OverwrittenValue() == \"bat\", true, got: %s, %v", v, ok)
	}
	if v := d2.Path().OriginalValue(); v != "baz" {
		t.Errorf("expected OriginalValue() == \"baz\", got: %v", v)
	}

	err = d2.Path().UnsetOverride()
	if err != nil {
		t.Errorf("expected no error in UnsetOverride, got: %s", err)
	}
	if v := d2.Get(); v != "baz" {
		t.Errorf("expected value baz, got: %v", v)
	}

	if v := d.Get(); v != "bar" {
		t.Errorf("expected value bar, got: %v", v)
	}

	// check that with push Path is copied and the path on the new .Push is
	// different from the original one
	if d.Path().Size()+1 != d2.Path().Size() {
		t.Errorf("expected d2.Path().Size() == d.Path().Size()+1, got %d and %d", d.Path().Size(), d2.Path().Size())
	}
}
