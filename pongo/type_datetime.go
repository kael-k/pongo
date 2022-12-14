package pongo

import (
	"encoding/json"
	"fmt"
	"time"
)

type DatetimeType struct {
	Format *ActionProperty[string] `json:"format,omitempty"`
	Cast   *ActionFlagProperty     `json:"cast,omitempty"`
	Before *TimeProperty           `json:"before,omitempty"`
	After  *TimeProperty           `json:"after,omitempty"`
}

func Datetime() *DatetimeType {
	return &DatetimeType{}
}

func (d DatetimeType) Process(action SchemaAction, dataPointer *DataPointer) (data Data, err error) {
	var t time.Time
	var ok bool

	if d.Cast.GetAction(action) {
		switch r := dataPointer.Get().(type) {
		case string:
			t, err = time.Parse(d.GetFormat(action), r)
			if err != nil {
				return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s cannot cast from string: %w", dataPointer.Path(), err))
			}
		case int:
			t = time.Unix(int64(r), 0)
		case int32:
			t = time.Unix(int64(r), 0)
		case int64:
			t = time.Unix(r, 0)
		case float64:
			t = time.Unix(int64(r), 0)
		case float32:
			t = time.Unix(int64(r), 0)
		case time.Time:
			t = r
		default:
			return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s cannot cast to \"Datetime\"", dataPointer.Path()))
		}

	} else {
		t, ok = dataPointer.Get().(time.Time)
		if !ok {
			return "", NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s is not a time.Time", dataPointer.Path()))
		}
	}

	if before, ok := d.Before.Get(); ok && before.Before(t) {
		return time.Time{}, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s value is %s (Max: %s)", dataPointer.Path(), t, before))
	}
	if after, ok := d.After.Get(); ok && after.After(t) {
		return time.Time{}, NewSchemaErrorWithError(dataPointer.Path(), fmt.Errorf("schema does not validate: %s length is %s (min: %s)", dataPointer.Path(), t, after))
	}

	switch action {
	case SchemaActionSerialize:
		return t.Format(time.RFC3339Nano), nil
	case SchemaActionParse:
		return t, nil
	}

	return nil, NewSchemaErrorWithError(dataPointer.Path(), NewErrInvalidAction(d, action))
}

func (d *DatetimeType) SetFormat(f string) *DatetimeType {
	d.Format = d.Format.SetDefault(f)
	return d
}

func (d *DatetimeType) SetFormatWithAction(action SchemaAction, f string) *DatetimeType {
	d.Format = d.Format.SetAction(action, f)
	return d
}

func (d DatetimeType) GetFormat(action SchemaAction) string {
	v, ok := d.Format.GetAction(action)

	if ok {
		return v
	}

	return time.RFC3339Nano
}
func (d *DatetimeType) SetCast(cast bool) *DatetimeType {
	d.Cast = d.Cast.Set(cast)
	return d
}

func (d DatetimeType) SetCastActions(actions ...SchemaAction) *DatetimeType {
	d.Cast = d.Cast.SetActions(actions...)
	return &d
}

func (d DatetimeType) UnsetCastActions(actions ...SchemaAction) *DatetimeType {
	d.Cast.UnsetActions(actions...)
	return &d
}

func (d *DatetimeType) SetBefore(f time.Time) *DatetimeType {
	d.Before = d.Before.Set(f)
	return d
}

func (d *DatetimeType) SetAfter(f time.Time) *DatetimeType {
	d.After = d.After.Set(f)
	return d
}

func (d *DatetimeType) SchemaTypeID() string {
	return "datetime"
}

func (d DatetimeType) MarshalJSONSchema(action SchemaAction) ([]byte, error) {
	if !d.Cast.GetAction(action) {
		return nil, fmt.Errorf("%w: Cast must be enabled in order to JSONschema-marshal the type", ErrSchemaNotJSONSchemaMarshalable)
	}

	if action != SchemaActionParse {
		return nil, NewErrInvalidAction(d, action)
	}

	return json.Marshal(map[string]interface{}{
		"type":   "string",
		"format": "date-time",
	})
}
