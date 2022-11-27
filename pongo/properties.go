package pongo

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// ActionProperty allow to set a property of generic T type that can have different value based
// on the specified action
type ActionProperty[T comparable] struct {
	Default *T                 `json:"default,omitempty"`
	Actions map[SchemaAction]T `json:"actions,omitempty"`
}

func (a *ActionProperty[T]) UnsetActions(actions ...SchemaAction) {
	if a.Actions == nil {
		return
	}
	for _, action := range actions {
		delete(a.Actions, action)
	}
}

func (a *ActionProperty[T]) SetAction(action SchemaAction, value T) {
	a.Actions[action] = value
}

func (a *ActionProperty[T]) SetDefault(value T) {
	a.Default = &value
}

func (a *ActionProperty[T]) UnsetDefault() {
	a.Default = nil
}

func (a *ActionProperty[T]) GetDefault() (value T, ok bool) {
	if a.Default != nil {
		return *a.Default, ok
	}

	return value, false
}

func (a *ActionProperty[T]) GetAction(action SchemaAction) (value T, ok bool) {
	if v, ok := a.Actions[action]; ok {
		return v, true
	}
	if a.Default != nil {
		return *a.Default, true
	}

	return value, false
}

// Empty returns true if ActionFlagProperty.all is false and if no action has been set
func (a ActionProperty[T]) Empty() bool {
	return a.Default == nil && (a.Actions == nil || len(a.Actions) == 0)
}

// ActionFlagProperty describe if the SchemaType has a specific flag given the requested action
// The cast should be enabled in 2 ways
// * all/none: using Set and Get, which enable casting for all Actions
// * only for specific action: using (Set|Unset|Reset)Actions and GetAction
// Note that if ActionFlagProperty.all has been set to true with Set (or with JSON unmarshal),
// then GetAction will always return true
type ActionFlagProperty struct {
	all     bool
	enabled map[SchemaAction]struct{}
}

func (a *ActionFlagProperty) SetActions(actions ...SchemaAction) {
	a.enabled = map[SchemaAction]struct{}{}
	for _, action := range actions {
		a.enabled[action] = struct{}{}
	}
}

func (a *ActionFlagProperty) AppendActions(actions ...SchemaAction) {
	if a.enabled == nil {
		a.enabled = map[SchemaAction]struct{}{}
	}
	for _, action := range actions {
		a.enabled[action] = struct{}{}
	}
}

func (a *ActionFlagProperty) UnsetActions(actions ...SchemaAction) {
	if a.enabled == nil {
		return
	}
	for _, action := range actions {
		delete(a.enabled, action)
	}
}

func (a *ActionFlagProperty) ResetActions() {
	a.enabled = map[SchemaAction]struct{}{}
}

func (a *ActionFlagProperty) Set(all bool) {
	a.all = all
}

func (a ActionFlagProperty) Get() bool {
	return a.all
}

func (a ActionFlagProperty) GetAction(action SchemaAction) bool {
	if a.all {
		return true
	}
	if a.enabled == nil {
		return false
	}
	_, ok := a.enabled[action]
	return ok
}

func (a ActionFlagProperty) GetActions() []SchemaAction {
	var l []SchemaAction

	for action := range a.enabled {
		l = append(l, action)
	}

	return l
}

// Empty returns true if ActionFlagProperty.all is false and if no action has been set
func (a ActionFlagProperty) Empty() bool {
	return !a.all && (a.enabled == nil || len(a.enabled) == 0)
}

func (a ActionFlagProperty) MarshalJSON() ([]byte, error) {
	if a.all {
		return json.Marshal(true)
	}
	if a.enabled == nil || len(a.enabled) == 0 {
		return json.Marshal(false)
	}

	var castOps []string
	for k := range a.enabled {
		castOps = append(castOps, string(k))
	}

	sort.Strings(castOps)

	return json.Marshal(castOps)
}

func (a *ActionFlagProperty) UnmarshalJSON(bytes []byte) error {
	var d interface{}
	err := json.Unmarshal(bytes, &d)
	if err != nil {
		return fmt.Errorf("error decoding ActionFlagProperty, got error: %w", err)
	}

	var actionsList []string

	switch v := d.(type) {
	case bool:
		a.all = v
		return nil
	case []string:
		actionsList = v
	case []interface{}:
		var ok = true
		var e string
		actionsList = []string{}
		for _, elem := range v {
			e, ok = elem.(string)
			if !ok {
				break
			}
			actionsList = append(actionsList, e)
		}
		if !ok {
			return fmt.Errorf("error decoding ActionFlagProperty, expected a bool or a list of string, got %v", v)
		}
	default:
		return fmt.Errorf("error decoding ActionFlagProperty, expected a bool or a list of string, got %v", v)
	}

	a.enabled = map[SchemaAction]struct{}{}
	for _, s := range actionsList {
		a.enabled[SchemaAction(s)] = struct{}{}
	}

	return nil
}

type NumberProperty[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	n *T
}

func (m *NumberProperty[T]) Set(i T) {
	m.n = &i
}

func (m *NumberProperty[T]) Unset() {
	m.n = nil
}

func (m *NumberProperty[T]) Get() (i T, ok bool) {
	if m.n == nil {
		return 0, false
	}
	return *m.n, true
}

func (m NumberProperty[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.n)
}

func (m *NumberProperty[T]) UnmarshalJSON(b []byte) error {
	var d *T
	err := json.Unmarshal(b, &d)
	if err != nil {
		return fmt.Errorf("error decoding NumberProperty, got error: %w", err)
	}

	m.n = d
	return nil
}

type TimeProperty struct {
	t *time.Time
}

func (b *TimeProperty) Set(i time.Time) {
	b.t = &i
}

func (b *TimeProperty) Unset() {
	b.t = nil
}

func (b *TimeProperty) Get() (i time.Time, ok bool) {
	if b.t == nil {
		return time.Time{}, false
	}
	return *b.t, true
}

func (b TimeProperty) MarshalJSON() ([]byte, error) {
	if b.t != nil {
		return json.Marshal(b.t.Format(time.RFC3339Nano))
	}
	return json.Marshal(nil)
}

func (b *TimeProperty) UnmarshalJSON(bytes []byte) error {
	var d *string
	err := json.Unmarshal(bytes, &d)
	if err != nil {
		return fmt.Errorf("error decoding NumberProperty, got error: %w", err)
	}

	if d == nil {
		b.t = nil
		return nil
	}

	var t time.Time
	t, err = time.Parse(time.RFC3339Nano, *d)
	if err != nil {
		return fmt.Errorf("error decoding NumberProperty, got error: %w", err)
	}

	b.t = &t

	return nil
}
