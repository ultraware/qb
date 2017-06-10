package qb

import "reflect"

// Scanner is implemented by database/sql.Row or database/sql.Rows
type Scanner interface {
	Scan(dest ...interface{}) error
}

// ScanToFields scans to a set of fields
func ScanToFields(f []DataField, r Scanner) error {
	dst := make([]interface{}, len(f))
	for k := range f {
		dst[k] = f[k].getScanTarget()
	}
	err := r.Scan(dst...)
	for k := range f {
		f[k].updateData(dst[k])
	}
	return err
}

// DataField ...
type DataField struct {
	Field
	*dataState
}

type dataState struct {
	Value        interface{}
	InitialValue interface{}
	empty        bool
}

// NewDataField returns a new DataField and checks if the given parameters are valid
func NewDataField(f Field, v interface{}) DataField {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		panic(`Cannot call NewDataField with non-pointer value`)
	}

	df := DataField{f, &dataState{Value: v, empty: true}}
	df.InitialValue = df.GetValue()
	return df
}

// GetValue returns the content of the current value
func (f *DataField) GetValue() interface{} {
	return reflect.ValueOf(f.Value).Elem().Interface()
}

// Empty checks if the current field is empty (nil)
func (f *DataField) Empty() bool {
	return f.empty
}

func (f *DataField) isSet() bool {
	return f.GetValue() != f.InitialValue
}

func (f *DataField) getScanTarget() interface{} {
	return reflect.New(reflect.TypeOf(f.Value)).Interface()
}

func (f *DataField) updateData(scanTarget interface{}) {
	if scanTarget == nil || reflect.ValueOf(scanTarget).Elem().IsNil() {
		reflect.ValueOf(f.Value).Elem().Set(reflect.Zero(reflect.TypeOf(f.InitialValue)))
		f.empty = true
		return
	}

	reflect.ValueOf(f.Value).Elem().Set(reflect.ValueOf(scanTarget).Elem().Elem())
	f.empty = false
}
