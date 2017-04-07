package qb

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// StringField is a field that can store a string
type StringField struct {
	Field
	Data string
}

func NewStringField(f Field) StringField {
	tf := StringField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *StringField) Scan(src interface{}) error {

	if v, ok := src.(string); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *StringField) Value() (driver.Value, error) {
	return f.Data, nil
}

// BoolField is a field that can store a bool
type BoolField struct {
	Field
	Data bool
}

func NewBoolField(f Field) BoolField {
	tf := BoolField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *BoolField) Scan(src interface{}) error {

	if v, ok := src.(bool); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *BoolField) Value() (driver.Value, error) {
	return f.Data, nil
}

// IntField is a field that can store a int64
type IntField struct {
	Field
	Data int64
}

func NewIntField(f Field) IntField {
	tf := IntField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *IntField) Scan(src interface{}) error {

	if v, ok := src.(int64); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *IntField) Value() (driver.Value, error) {
	return f.Data, nil
}

// FloatField is a field that can store a float64
type FloatField struct {
	Field
	Data float64
}

func NewFloatField(f Field) FloatField {
	tf := FloatField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *FloatField) Scan(src interface{}) error {

	if v, ok := src.(float64); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *FloatField) Value() (driver.Value, error) {
	return f.Data, nil
}

// BytesField is a field that can store a []byte
type BytesField struct {
	Field
	Data []byte
}

func NewBytesField(f Field) BytesField {
	tf := BytesField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *BytesField) Scan(src interface{}) error {

	if v, ok := src.([]byte); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *BytesField) Value() (driver.Value, error) {
	return f.Data, nil
}

// TimeField is a field that can store a time.Time
type TimeField struct {
	Field
	Data time.Time
}

func NewTimeField(f Field) TimeField {
	tf := TimeField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *TimeField) Scan(src interface{}) error {

	if v, ok := src.(time.Time); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *TimeField) Value() (driver.Value, error) {
	return f.Data, nil
}

// NullStringField is a field that can store a string, can be nil
type NullStringField struct {
	Field
	Data *string
}

func NewNullStringField(f Field) NullStringField {
	tf := NullStringField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *NullStringField) Scan(src interface{}) error {
	if src == nil {
		f.Data = nil
		return nil
	}
	if v, ok := src.(*string); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *NullStringField) Value() (driver.Value, error) {
	return f.Data, nil
}

// NullBoolField is a field that can store a bool, can be nil
type NullBoolField struct {
	Field
	Data *bool
}

func NewNullBoolField(f Field) NullBoolField {
	tf := NullBoolField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *NullBoolField) Scan(src interface{}) error {
	if src == nil {
		f.Data = nil
		return nil
	}
	if v, ok := src.(*bool); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *NullBoolField) Value() (driver.Value, error) {
	return f.Data, nil
}

// NullIntField is a field that can store a int64, can be nil
type NullIntField struct {
	Field
	Data *int64
}

func NewNullIntField(f Field) NullIntField {
	tf := NullIntField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *NullIntField) Scan(src interface{}) error {
	if src == nil {
		f.Data = nil
		return nil
	}
	if v, ok := src.(*int64); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *NullIntField) Value() (driver.Value, error) {
	return f.Data, nil
}

// NullFloatField is a field that can store a float64, can be nil
type NullFloatField struct {
	Field
	Data *float64
}

func NewNullFloatField(f Field) NullFloatField {
	tf := NullFloatField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *NullFloatField) Scan(src interface{}) error {
	if src == nil {
		f.Data = nil
		return nil
	}
	if v, ok := src.(*float64); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *NullFloatField) Value() (driver.Value, error) {
	return f.Data, nil
}

// NullBytesField is a field that can store a []byte, can be nil
type NullBytesField struct {
	Field
	Data *[]byte
}

func NewNullBytesField(f Field) NullBytesField {
	tf := NullBytesField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *NullBytesField) Scan(src interface{}) error {
	if src == nil {
		f.Data = nil
		return nil
	}
	if v, ok := src.(*[]byte); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *NullBytesField) Value() (driver.Value, error) {
	return f.Data, nil
}

// NullTimeField is a field that can store a time.Time, can be nil
type NullTimeField struct {
	Field
	Data *time.Time
}

func NewNullTimeField(f Field) NullTimeField {
	tf := NullTimeField{}
	tf.Field = f
	return tf
}

// Scan implements sql.Scanner
func (f *NullTimeField) Scan(src interface{}) error {
	if src == nil {
		f.Data = nil
		return nil
	}
	if v, ok := src.(*time.Time); ok {
		f.Data = v
		return nil
	}
	return fmt.Errorf(`Unsupported scan, cannot scan %T into %T`, src, f.Data)
}

// Value implements driver.Valuer
func (f *NullTimeField) Value() (driver.Value, error) {
	return f.Data, nil
}
