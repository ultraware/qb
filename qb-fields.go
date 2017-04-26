package qb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// StringField is a field that can store a string
type StringField struct {
	Field
	data    string
	changed bool
	set     bool
}

// NewStringField creates a new StringField
func NewStringField(f Field) *StringField {
	tf := StringField{Field: f}
	return &tf
}

func (f *StringField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *StringField) Scan(src interface{}) error {
	if v, ok := src.(string); ok {
		f.setData(v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *StringField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *StringField) Get() string {
	return f.data
}

// Set updates the data
func (f *StringField) Set(v string) {
	f.setData(v, true)
}

func (f *StringField) setData(v string, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *StringField) hasChanged() bool {
	return f.changed
}

func (f *StringField) isSet() bool {
	return f.set
}

// BoolField is a field that can store a bool
type BoolField struct {
	Field
	data    bool
	changed bool
	set     bool
}

// NewBoolField creates a new BoolField
func NewBoolField(f Field) *BoolField {
	tf := BoolField{Field: f}
	return &tf
}

func (f *BoolField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *BoolField) Scan(src interface{}) error {
	if v, ok := src.(bool); ok {
		f.setData(v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *BoolField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *BoolField) Get() bool {
	return f.data
}

// Set updates the data
func (f *BoolField) Set(v bool) {
	f.setData(v, true)
}

func (f *BoolField) setData(v bool, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *BoolField) hasChanged() bool {
	return f.changed
}

func (f *BoolField) isSet() bool {
	return f.set
}

// IntField is a field that can store a int
type IntField struct {
	Field
	data    int
	changed bool
	set     bool
}

// NewIntField creates a new IntField
func NewIntField(f Field) *IntField {
	tf := IntField{Field: f}
	return &tf
}

func (f *IntField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *IntField) Scan(src interface{}) error {
	if v, ok := src.(int64); ok {
		data := int(v)
		f.setData(data, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *IntField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *IntField) Get() int {
	return f.data
}

// Set updates the data
func (f *IntField) Set(v int) {
	f.setData(v, true)
}

func (f *IntField) setData(v int, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *IntField) hasChanged() bool {
	return f.changed
}

func (f *IntField) isSet() bool {
	return f.set
}

// Int64Field is a field that can store a int64
type Int64Field struct {
	Field
	data    int64
	changed bool
	set     bool
}

// NewInt64Field creates a new Int64Field
func NewInt64Field(f Field) *Int64Field {
	tf := Int64Field{Field: f}
	return &tf
}

func (f *Int64Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *Int64Field) Scan(src interface{}) error {
	if v, ok := src.(int64); ok {
		f.setData(v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *Int64Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *Int64Field) Get() int64 {
	return f.data
}

// Set updates the data
func (f *Int64Field) Set(v int64) {
	f.setData(v, true)
}

func (f *Int64Field) setData(v int64, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *Int64Field) hasChanged() bool {
	return f.changed
}

func (f *Int64Field) isSet() bool {
	return f.set
}

// Int32Field is a field that can store a int32
type Int32Field struct {
	Field
	data    int32
	changed bool
	set     bool
}

// NewInt32Field creates a new Int32Field
func NewInt32Field(f Field) *Int32Field {
	tf := Int32Field{Field: f}
	return &tf
}

func (f *Int32Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *Int32Field) Scan(src interface{}) error {
	if v, ok := src.(int64); ok {
		data := int32(v)
		f.setData(data, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *Int32Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *Int32Field) Get() int32 {
	return f.data
}

// Set updates the data
func (f *Int32Field) Set(v int32) {
	f.setData(v, true)
}

func (f *Int32Field) setData(v int32, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *Int32Field) hasChanged() bool {
	return f.changed
}

func (f *Int32Field) isSet() bool {
	return f.set
}

// Float64Field is a field that can store a float64
type Float64Field struct {
	Field
	data    float64
	changed bool
	set     bool
}

// NewFloat64Field creates a new Float64Field
func NewFloat64Field(f Field) *Float64Field {
	tf := Float64Field{Field: f}
	return &tf
}

func (f *Float64Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *Float64Field) Scan(src interface{}) error {
	nf := sql.NullFloat64{}
	err := nf.Scan(src)
	if err != nil || !nf.Valid {
		return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
	}

	v := (nf.Float64)

	f.setData(v, false)
	return nil
}

// Value implements driver.Valuer
func (f *Float64Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *Float64Field) Get() float64 {
	return f.data
}

// Set updates the data
func (f *Float64Field) Set(v float64) {
	f.setData(v, true)
}

func (f *Float64Field) setData(v float64, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *Float64Field) hasChanged() bool {
	return f.changed
}

func (f *Float64Field) isSet() bool {
	return f.set
}

// Float32Field is a field that can store a float32
type Float32Field struct {
	Field
	data    float32
	changed bool
	set     bool
}

// NewFloat32Field creates a new Float32Field
func NewFloat32Field(f Field) *Float32Field {
	tf := Float32Field{Field: f}
	return &tf
}

func (f *Float32Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *Float32Field) Scan(src interface{}) error {
	nf := sql.NullFloat64{}
	err := nf.Scan(src)
	if err != nil || !nf.Valid {
		return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
	}

	v := float32(nf.Float64)

	f.setData(v, false)
	return nil
}

// Value implements driver.Valuer
func (f *Float32Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *Float32Field) Get() float32 {
	return f.data
}

// Set updates the data
func (f *Float32Field) Set(v float32) {
	f.setData(v, true)
}

func (f *Float32Field) setData(v float32, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *Float32Field) hasChanged() bool {
	return f.changed
}

func (f *Float32Field) isSet() bool {
	return f.set
}

// BytesField is a field that can store a []byte
type BytesField struct {
	Field
	data    []byte
	changed bool
	set     bool
}

// NewBytesField creates a new BytesField
func NewBytesField(f Field) *BytesField {
	tf := BytesField{Field: f}
	return &tf
}

func (f *BytesField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *BytesField) Scan(src interface{}) error {
	if v, ok := src.([]byte); ok {
		f.setData(v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *BytesField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *BytesField) Get() []byte {
	return f.data
}

// Set updates the data
func (f *BytesField) Set(v []byte) {
	f.setData(v, true)
}

func (f *BytesField) setData(v []byte, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *BytesField) hasChanged() bool {
	return f.changed
}

func (f *BytesField) isSet() bool {
	return f.set
}

// TimeField is a field that can store a time.Time
type TimeField struct {
	Field
	data    time.Time
	changed bool
	set     bool
}

// NewTimeField creates a new TimeField
func NewTimeField(f Field) *TimeField {
	tf := TimeField{Field: f}
	return &tf
}

func (f *TimeField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *TimeField) Scan(src interface{}) error {
	if v, ok := src.(time.Time); ok {
		f.setData(v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *TimeField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *TimeField) Get() time.Time {
	return f.data
}

// Set updates the data
func (f *TimeField) Set(v time.Time) {
	f.setData(v, true)
}

func (f *TimeField) setData(v time.Time, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *TimeField) hasChanged() bool {
	return f.changed
}

func (f *TimeField) isSet() bool {
	return f.set
}

// NullStringField is a field that can store a string, can be nil
type NullStringField struct {
	Field
	data    *string
	changed bool
	set     bool
}

// NewNullStringField creates a new NullStringField
func NewNullStringField(f Field) *NullStringField {
	tf := NullStringField{Field: f}
	return &tf
}

func (f *NullStringField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullStringField) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.(string); ok {
		f.setData(&v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullStringField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullStringField) Get() *string {
	return f.data
}

// Set updates the data
func (f *NullStringField) Set(v *string) {
	f.setData(v, true)
}

func (f *NullStringField) setData(v *string, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullStringField) hasChanged() bool {
	return f.changed
}

func (f *NullStringField) isSet() bool {
	return f.set
}

// NullBoolField is a field that can store a bool, can be nil
type NullBoolField struct {
	Field
	data    *bool
	changed bool
	set     bool
}

// NewNullBoolField creates a new NullBoolField
func NewNullBoolField(f Field) *NullBoolField {
	tf := NullBoolField{Field: f}
	return &tf
}

func (f *NullBoolField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullBoolField) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.(bool); ok {
		f.setData(&v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullBoolField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullBoolField) Get() *bool {
	return f.data
}

// Set updates the data
func (f *NullBoolField) Set(v *bool) {
	f.setData(v, true)
}

func (f *NullBoolField) setData(v *bool, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullBoolField) hasChanged() bool {
	return f.changed
}

func (f *NullBoolField) isSet() bool {
	return f.set
}

// NullIntField is a field that can store a int, can be nil
type NullIntField struct {
	Field
	data    *int
	changed bool
	set     bool
}

// NewNullIntField creates a new NullIntField
func NewNullIntField(f Field) *NullIntField {
	tf := NullIntField{Field: f}
	return &tf
}

func (f *NullIntField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullIntField) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.(int64); ok {
		data := int(v)
		f.setData(&data, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullIntField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullIntField) Get() *int {
	return f.data
}

// Set updates the data
func (f *NullIntField) Set(v *int) {
	f.setData(v, true)
}

func (f *NullIntField) setData(v *int, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullIntField) hasChanged() bool {
	return f.changed
}

func (f *NullIntField) isSet() bool {
	return f.set
}

// NullInt64Field is a field that can store a int64, can be nil
type NullInt64Field struct {
	Field
	data    *int64
	changed bool
	set     bool
}

// NewNullInt64Field creates a new NullInt64Field
func NewNullInt64Field(f Field) *NullInt64Field {
	tf := NullInt64Field{Field: f}
	return &tf
}

func (f *NullInt64Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullInt64Field) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.(int64); ok {
		f.setData(&v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullInt64Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullInt64Field) Get() *int64 {
	return f.data
}

// Set updates the data
func (f *NullInt64Field) Set(v *int64) {
	f.setData(v, true)
}

func (f *NullInt64Field) setData(v *int64, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullInt64Field) hasChanged() bool {
	return f.changed
}

func (f *NullInt64Field) isSet() bool {
	return f.set
}

// NullInt32Field is a field that can store a int32, can be nil
type NullInt32Field struct {
	Field
	data    *int32
	changed bool
	set     bool
}

// NewNullInt32Field creates a new NullInt32Field
func NewNullInt32Field(f Field) *NullInt32Field {
	tf := NullInt32Field{Field: f}
	return &tf
}

func (f *NullInt32Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullInt32Field) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.(int64); ok {
		data := int32(v)
		f.setData(&data, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullInt32Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullInt32Field) Get() *int32 {
	return f.data
}

// Set updates the data
func (f *NullInt32Field) Set(v *int32) {
	f.setData(v, true)
}

func (f *NullInt32Field) setData(v *int32, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullInt32Field) hasChanged() bool {
	return f.changed
}

func (f *NullInt32Field) isSet() bool {
	return f.set
}

// NullFloat64Field is a field that can store a float64, can be nil
type NullFloat64Field struct {
	Field
	data    *float64
	changed bool
	set     bool
}

// NewNullFloat64Field creates a new NullFloat64Field
func NewNullFloat64Field(f Field) *NullFloat64Field {
	tf := NullFloat64Field{Field: f}
	return &tf
}

func (f *NullFloat64Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullFloat64Field) Scan(src interface{}) error {
	nf := sql.NullFloat64{}
	err := nf.Scan(src)
	if err != nil {
		return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
	}
	if !nf.Valid {
		f.setData(nil, false)
		return nil
	}

	v := (nf.Float64)

	f.setData(&v, false)
	return nil
}

// Value implements driver.Valuer
func (f *NullFloat64Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullFloat64Field) Get() *float64 {
	return f.data
}

// Set updates the data
func (f *NullFloat64Field) Set(v *float64) {
	f.setData(v, true)
}

func (f *NullFloat64Field) setData(v *float64, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullFloat64Field) hasChanged() bool {
	return f.changed
}

func (f *NullFloat64Field) isSet() bool {
	return f.set
}

// NullFloat32Field is a field that can store a float32, can be nil
type NullFloat32Field struct {
	Field
	data    *float32
	changed bool
	set     bool
}

// NewNullFloat32Field creates a new NullFloat32Field
func NewNullFloat32Field(f Field) *NullFloat32Field {
	tf := NullFloat32Field{Field: f}
	return &tf
}

func (f *NullFloat32Field) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullFloat32Field) Scan(src interface{}) error {
	nf := sql.NullFloat64{}
	err := nf.Scan(src)
	if err != nil {
		return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
	}
	if !nf.Valid {
		f.setData(nil, false)
		return nil
	}

	v := float32(nf.Float64)

	f.setData(&v, false)
	return nil
}

// Value implements driver.Valuer
func (f *NullFloat32Field) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullFloat32Field) Get() *float32 {
	return f.data
}

// Set updates the data
func (f *NullFloat32Field) Set(v *float32) {
	f.setData(v, true)
}

func (f *NullFloat32Field) setData(v *float32, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullFloat32Field) hasChanged() bool {
	return f.changed
}

func (f *NullFloat32Field) isSet() bool {
	return f.set
}

// NullBytesField is a field that can store a []byte, can be nil
type NullBytesField struct {
	Field
	data    *[]byte
	changed bool
	set     bool
}

// NewNullBytesField creates a new NullBytesField
func NewNullBytesField(f Field) *NullBytesField {
	tf := NullBytesField{Field: f}
	return &tf
}

func (f *NullBytesField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullBytesField) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.([]byte); ok {
		f.setData(&v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullBytesField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullBytesField) Get() *[]byte {
	return f.data
}

// Set updates the data
func (f *NullBytesField) Set(v *[]byte) {
	f.setData(v, true)
}

func (f *NullBytesField) setData(v *[]byte, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullBytesField) hasChanged() bool {
	return f.changed
}

func (f *NullBytesField) isSet() bool {
	return f.set
}

// NullTimeField is a field that can store a time.Time, can be nil
type NullTimeField struct {
	Field
	data    *time.Time
	changed bool
	set     bool
}

// NewNullTimeField creates a new NullTimeField
func NewNullTimeField(f Field) *NullTimeField {
	tf := NullTimeField{Field: f}
	return &tf
}

func (f *NullTimeField) getField() Field {
	return f.Field
}

// Scan implements sql.Scanner
func (f *NullTimeField) Scan(src interface{}) error {
	if src == nil {
		f.setData(nil, false)
		return nil
	}
	if v, ok := src.(time.Time); ok {
		f.setData(&v, false)
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
}

// Value implements driver.Valuer
func (f *NullTimeField) Value() (driver.Value, error) {
	return f.data, nil
}

// Get returns the data
func (f *NullTimeField) Get() *time.Time {
	return f.data
}

// Set updates the data
func (f *NullTimeField) Set(v *time.Time) {
	f.setData(v, true)
}

func (f *NullTimeField) setData(v *time.Time, changed bool) {
	f.changed = changed
	f.set = true
	f.data = v
}

func (f *NullTimeField) hasChanged() bool {
	return f.changed
}

func (f *NullTimeField) isSet() bool {
	return f.set
}
