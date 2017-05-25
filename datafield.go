package qb

import "reflect"

// Scanner is implemented by database/sql.Row or database/sql.Rows
type Scanner interface {
	Scan(dest ...interface{}) error
}

// ScanToFields scans to a set of fields
func ScanToFields(r Scanner, f []DataField) error {
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
	df.UpdateInit()
	return df
}

// UpdateInit sets the initial value to be the same as the current value
func (f *DataField) UpdateInit() {
	f.InitialValue = f.Get()
}

// Get *DEPRECATED* returns a copy of the current value
func (f *DataField) Get() interface{} {
	return reflect.ValueOf(f.Value).Elem().Interface()
}

// Set *DEPRECATED* changes the value of the value
func (f *DataField) Set(v interface{}) {
	reflect.ValueOf(f.Value).Elem().Set(reflect.ValueOf(v))
}

// Empty checks if the current field is empty (nil)
func (f *DataField) Empty() bool {
	return f.empty
}

func (f *DataField) isSet() bool {
	return f.hasChanged()
}

func (f *DataField) hasChanged() bool {
	return f.Get() != f.InitialValue
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
