package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kael-k/pongo/v2/pongo"
)

func TestCastProperty(t *testing.T) {
	var test pongo.ActionFlagProperty

	if test.Get() {
		t.Errorf("error on TestCastProperty, expected Get() == false on new object, got true")
	}
	if test.GetAction(pongo.SchemaActionSerialize) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionSerialize) == false on new object, got true")
	}

	if v := test.GetActions(); v != nil {
		t.Errorf("error on TestCastProperty, expected GetActions() on new object, got %v", v)
	}

	test.SetActions(pongo.SchemaActionParse)
	if !test.GetAction(pongo.SchemaActionParse) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionParse) true, got false")
	}
	if test.GetAction(pongo.SchemaActionSerialize) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionSerialize) false, got true")
	}

	test.Set(true)
	if !test.GetAction(pongo.SchemaActionSerialize) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionSerialize) true, got false")
	}
	if !test.Get() {
		t.Errorf("error on TestCastProperty, expected Get() == true, got false")
	}
	marshaled, err := json.Marshal(test)
	if err != nil {
		t.Errorf("error marshall TestCastProperty: %s", err)
	}
	if string(marshaled) != "true" {
		t.Errorf("error marshall TestNumberProperty: expected %s, got %s", "true", string(marshaled))
	}

	test.Set(false)

	test.UnsetActions(pongo.SchemaActionSerialize)
	if !test.GetAction(pongo.SchemaActionParse) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionParse) true, got false")
	}
	if test.GetAction(pongo.SchemaActionSerialize) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionSerialize) false, got true")
	}
	marshaled, err = json.Marshal(test)
	if err != nil {
		t.Errorf("error marshall TestCastProperty: %s", err)
	}
	if string(marshaled) != ("[\"" + string(pongo.SchemaActionParse) + "\"]") {
		t.Errorf("error marshall TestNumberProperty: expected %s, got %s", "true", string(marshaled))
	}

	test.AppendActions(pongo.SchemaActionSerialize)
	if !test.GetAction(pongo.SchemaActionParse) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionParse) true, got false")
	}
	if !test.GetAction(pongo.SchemaActionSerialize) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionSerialize) true, got false")
	}

	test.ResetActions()
	if test.GetAction(pongo.SchemaActionParse) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionParse) false, got true")
	}
	if test.GetAction(pongo.SchemaActionSerialize) {
		t.Errorf("error on TestCastProperty, expected test.GetAction(pongo.SchemaActionSerialize) false, got true")
	}
}

func TestNumberProperty(t *testing.T) {
	var m pongo.NumberProperty[int]
	var unmarshall pongo.NumberProperty[int]

	const TestSet = 42

	if _, ok := m.Get(); ok {
		t.Errorf("error on TestNumberProperty, expected ok == false on new object, got true")
	}

	m.Set(TestSet)
	v, ok := m.Get()
	if !ok {
		t.Errorf("error on TestNumberProperty, expected ok == true after set, got false")
	}

	if v != TestSet {
		t.Errorf("error on TestNumberProperty, expected value == %d after set, got %d", TestSet, v)
	}

	marshaled, err := json.Marshal(m)
	if err != nil {
		t.Errorf("error marshall TestNumberProperty: %s", err)
	}
	if string(marshaled) != fmt.Sprintf("%d", TestSet) {
		t.Errorf("error marshall TestNumberProperty: expected %s, got %s", fmt.Sprintf("%d", TestSet), string(marshaled))
	}

	err = json.Unmarshal(marshaled, &unmarshall)
	if err != nil {
		t.Errorf("error unmarshall TestNumberProperty: %s", err)
	}
	if uV, uOk := unmarshall.Get(); ok != uOk || uV != v {
		t.Errorf("error unmarshall TestNumberProperty: unmarshalled (%v, %v) and original (%v, %v) mismatch", uV, uOk, v, ok)
	}

	m.Unset()

	if _, ok := m.Get(); ok {
		t.Errorf("error on TestNumberProperty, expected ok == false after unset, got true")
	}

	marshaled, err = json.Marshal(m)
	if err != nil {
		t.Errorf("error marshall TestNumberProperty: %s", err)
	}
	if string(marshaled) != "null" {
		t.Errorf("error marshall TestNumberProperty: expected %s, got %s", "null", string(marshaled))
	}

	err = json.Unmarshal(marshaled, &unmarshall)
	if err != nil {
		t.Errorf("error unmarshall TestNumberProperty: %s", err)
	}
	if _, uOk := unmarshall.Get(); uOk {
		t.Errorf("error unmarshall TestNumberProperty: expected ok == false instead ok == true")
	}
}

func TestTimeProperty(t *testing.T) {
	var m pongo.TimeProperty
	var unmarshall pongo.TimeProperty
	var TestSet = time.Unix(1663800000, 0)

	if _, ok := m.Get(); ok {
		t.Errorf("error on TestTimeProperty, expected ok == false on new object, got true")
	}

	m.Set(TestSet)
	v, ok := m.Get()
	if !ok {
		t.Errorf("error on TestTimeProperty, expected ok == true after set, got false")
	}

	if !v.Equal(TestSet) {
		t.Errorf("error on TestTimeProperty, expected value == %v after set, got %v", TestSet, v)
	}

	marshaled, err := json.Marshal(m)
	if err != nil {
		t.Errorf("error marshall TestTimeProperty: %s", err)
	}
	if want := fmt.Sprintf("\"%s\"", TestSet.Format(time.RFC3339Nano)); string(marshaled) != want {
		t.Errorf("error marshall TestTimeProperty: expected %s, got %s", want, string(marshaled))
	}

	err = json.Unmarshal(marshaled, &unmarshall)
	if err != nil {
		t.Errorf("error unmarshall TestTimeProperty: %s", err)
	}
	if uV, uOk := unmarshall.Get(); ok != uOk || !uV.Equal(v) {
		t.Errorf("error unmarshall TestTimeProperty: unmarshalled (%v, %v) and original (%v, %v) mismatch", uV, uOk, v, ok)
	}

	m.Unset()

	if _, ok := m.Get(); ok {
		t.Errorf("error on TestTimeProperty, expected ok == false after unset, got true")
	}

	marshaled, err = json.Marshal(m)
	if err != nil {
		t.Errorf("error marshall TestTimeProperty: %s", err)
	}
	if string(marshaled) != "null" {
		t.Errorf("error marshall TestTimeProperty: expected %s, got %s", "null", string(marshaled))
	}

	err = json.Unmarshal(marshaled, &unmarshall)
	if err != nil {
		t.Errorf("error unmarshall TestTimeProperty: %s", err)
	}
	if _, uOk := unmarshall.Get(); uOk {
		t.Errorf("error unmarshall TestTimeProperty: expected ok == false instead ok == true")
	}
}
