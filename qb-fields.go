package qb

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// StringField is a field that can store a string
type StringField struct {
	Field
	data string
}

func NewStringField(f Field) *StringField {
	tf := StringField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *StringField) Scan(src interface{}) error {
	if v, ok := src.(string); ok {
		data := v
		f.data = data
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
	f.data = v
}

// BoolField is a field that can store a bool
type BoolField struct {
	Field
	data bool
}

func NewBoolField(f Field) *BoolField {
	tf := BoolField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *BoolField) Scan(src interface{}) error {
	if v, ok := src.(bool); ok {
		data := v
		f.data = data
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
	f.data = v
}

// IntField is a field that can store a int
type IntField struct {
	Field
	data int
}

func NewIntField(f Field) *IntField {
	tf := IntField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *IntField) Scan(src interface{}) error {
	if v, ok := src.(int64); ok {
		data := int(v)
		f.data = data
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
	f.data = v
}

// Int64Field is a field that can store a int64
type Int64Field struct {
	Field
	data int64
}

func NewInt64Field(f Field) *Int64Field {
	tf := Int64Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *Int64Field) Scan(src interface{}) error {
	if v, ok := src.(int64); ok {
		data := v
		f.data = data
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
	f.data = v
}

// Int32Field is a field that can store a int32
type Int32Field struct {
	Field
	data int32
}

func NewInt32Field(f Field) *Int32Field {
	tf := Int32Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *Int32Field) Scan(src interface{}) error {
	if v, ok := src.(int64); ok {
		data := int32(v)
		f.data = data
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
	f.data = v
}

// Float64Field is a field that can store a float64
type Float64Field struct {
	Field
	data float64
}

func NewFloat64Field(f Field) *Float64Field {
	tf := Float64Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *Float64Field) Scan(src interface{}) error {
	if v, ok := src.(float64); ok {
		data := v
		f.data = data
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
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
	f.data = v
}

// Float32Field is a field that can store a float32
type Float32Field struct {
	Field
	data float32
}

func NewFloat32Field(f Field) *Float32Field {
	tf := Float32Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *Float32Field) Scan(src interface{}) error {
	if v, ok := src.(float64); ok {
		data := float32(v)
		f.data = data
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
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
	f.data = v
}

// BytesField is a field that can store a []byte
type BytesField struct {
	Field
	data []byte
}

func NewBytesField(f Field) *BytesField {
	tf := BytesField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *BytesField) Scan(src interface{}) error {
	if v, ok := src.([]byte); ok {
		data := v
		f.data = data
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
	f.data = v
}

// TimeField is a field that can store a time.Time
type TimeField struct {
	Field
	data time.Time
}

func NewTimeField(f Field) *TimeField {
	tf := TimeField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *TimeField) Scan(src interface{}) error {
	if v, ok := src.(time.Time); ok {
		data := v
		f.data = data
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
	f.data = v
}

// NullStringField is a field that can store a string, can be nil
type NullStringField struct {
	Field
	data *string
}

func NewNullStringField(f Field) *NullStringField {
	tf := NullStringField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullStringField) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(string); ok {
		data := v
		f.data = &data
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
	f.data = v
}

// NullBoolField is a field that can store a bool, can be nil
type NullBoolField struct {
	Field
	data *bool
}

func NewNullBoolField(f Field) *NullBoolField {
	tf := NullBoolField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullBoolField) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(bool); ok {
		data := v
		f.data = &data
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
	f.data = v
}

// NullIntField is a field that can store a int, can be nil
type NullIntField struct {
	Field
	data *int
}

func NewNullIntField(f Field) *NullIntField {
	tf := NullIntField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullIntField) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(int64); ok {
		data := int(v)
		f.data = &data
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
	f.data = v
}

// NullInt64Field is a field that can store a int64, can be nil
type NullInt64Field struct {
	Field
	data *int64
}

func NewNullInt64Field(f Field) *NullInt64Field {
	tf := NullInt64Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullInt64Field) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(int64); ok {
		data := v
		f.data = &data
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
	f.data = v
}

// NullInt32Field is a field that can store a int32, can be nil
type NullInt32Field struct {
	Field
	data *int32
}

func NewNullInt32Field(f Field) *NullInt32Field {
	tf := NullInt32Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullInt32Field) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(int64); ok {
		data := int32(v)
		f.data = &data
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
	f.data = v
}

// NullFloat64Field is a field that can store a float64, can be nil
type NullFloat64Field struct {
	Field
	data *float64
}

func NewNullFloat64Field(f Field) *NullFloat64Field {
	tf := NullFloat64Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullFloat64Field) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(float64); ok {
		data := v
		f.data = &data
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
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
	f.data = v
}

// NullFloat32Field is a field that can store a float32, can be nil
type NullFloat32Field struct {
	Field
	data *float32
}

func NewNullFloat32Field(f Field) *NullFloat32Field {
	tf := NullFloat32Field{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullFloat32Field) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(float64); ok {
		data := float32(v)
		f.data = &data
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.data)
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
	f.data = v
}

// NullBytesField is a field that can store a []byte, can be nil
type NullBytesField struct {
	Field
	data *[]byte
}

func NewNullBytesField(f Field) *NullBytesField {
	tf := NullBytesField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullBytesField) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.([]byte); ok {
		data := v
		f.data = &data
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
	f.data = v
}

// NullTimeField is a field that can store a time.Time, can be nil
type NullTimeField struct {
	Field
	data *time.Time
}

func NewNullTimeField(f Field) *NullTimeField {
	tf := NullTimeField{}
	tf.Field = f
	return &tf
}

// Scan implements sql.Scanner
func (f *NullTimeField) Scan(src interface{}) error {
	if src == nil {
		f.data = nil
		return nil
	}
	if v, ok := src.(time.Time); ok {
		data := v
		f.data = &data
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
	f.data = v
}
