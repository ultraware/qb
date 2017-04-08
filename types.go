package qb

//go:generate go run scripts/main.go

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

///
/// Source
///

// Source ...
type Source interface {
	QueryString(*AliasGenerator, *ValueList) string
	AliasString() string
}

// Table ...
type Table struct {
	Name string
}

// QueryString ...
func (t *Table) QueryString(ag *AliasGenerator, _ *ValueList) string {
	return t.Name + ` ` + ag.Get(t)
}

// AliasString ...
func (t *Table) AliasString() string {
	return `t`
}

// Select ...
func (t *Table) Select(f ...DataField) SelectBuilder {
	return NewSelectBuilder(f, t)
}

// SubQuery ...
type SubQuery struct {
	sql    string
	values []interface{}
	fields []Field
}

// QueryString ...
func (t *SubQuery) QueryString(ag *AliasGenerator, vl *ValueList) string {
	vl.Append(t.values...)
	return `(` + t.sql + `) ` + ag.Get(t)
}

// AliasString ...
func (t *SubQuery) AliasString() string {
	return `sq`
}

// Select ...
func (t *SubQuery) Select(f ...DataField) SelectBuilder {
	return NewSelectBuilder(f, t)
}

///
/// Field
///

// Field represents a field in a query
type Field interface {
	QueryString(*AliasGenerator, *ValueList) string
	Source() Source
	DataType() string
}

// DataField is a Field that can store data
type DataField interface {
	Field
	sql.Scanner
	driver.Valuer
}

// Cursor ...
type Cursor struct {
	dst  []interface{}
	rows *sql.Rows
	err  error
}

// NewCursor returns a new Cursor
func NewCursor(f []DataField, r *sql.Rows) *Cursor {
	dst := make([]interface{}, len(f))
	for k, v := range f {
		dst[k] = v
	}
	c := &Cursor{dst: dst, rows: r}
	return c
}

// Next loads the next row into the fields
func (c Cursor) Next() bool {
	if !c.rows.Next() {
		c.Close()
		return false
	}
	err := c.rows.Scan(c.dst...)
	if err != nil {
		c.err = err
		return false
	}

	return true
}

// Close the rows object
func (c Cursor) Close() {
	c.rows.Close()
}

// Error returns the last error
func (c Cursor) Error() error {
	return c.err
}

// TableField represents a real field in a table
type TableField struct {
	Parent   Source
	Name     string
	Type     string
	ReadOnly bool
}

// QueryString ...
func (f TableField) QueryString(ag *AliasGenerator, _ *ValueList) string {
	return ag.Get(f.Parent) + `.` + f.Name
}

// Source ...
func (f TableField) Source() Source {
	return f.Parent
}

// DataType ...
func (f TableField) DataType() string {
	return f.Type
}

// ValueField contains values supplied by the program
type ValueField struct {
	Value interface{}
	Type  string
}

func getValue(v interface{}) interface{} {
	if val, ok := v.(driver.Valuer); ok {
		new, err := val.Value()
		if err == nil {
			v = new
		}
	}
	return v
}

func getType(v interface{}) string {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return `int`
	case float32, float64:
		return `float`
	case string:
		return `string`
	case bool:
		return `bool`
	case time.Time:
		return `time`
	default:
		return `unknown`
	}
}

// Value ...
func Value(v interface{}) ValueField {
	v = getValue(v)
	return ValueField{v, getType(v)}
}

// QueryString ...
func (f ValueField) QueryString(_ *AliasGenerator, vl *ValueList) string {
	vl.Append(f.Value)
	return `?`
}

// Source ...
func (f ValueField) Source() Source {
	return nil
}

// DataType ...
func (f ValueField) DataType() string {
	return f.Type
}

///
/// Query types
///

// Condition is used in the Where function
type Condition func(*AliasGenerator, *ValueList) string

// Join are used for joins on tables
type Join struct {
	Type   string
	Fields [2]Field
	New    Source
}

// FieldOrder specifies the order in which fields should be sorted
type FieldOrder struct {
	Field Field
	Order string
}

// Asc is used to sort in ascending order
func Asc(f Field) FieldOrder {
	return FieldOrder{Field: f, Order: `ASC`}
}

// Desc is used to sort in descending order
func Desc(f Field) FieldOrder {
	return FieldOrder{Field: f, Order: `DESC`}
}
