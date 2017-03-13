package qb

import (
	"database/sql/driver"
	"time"
)

///
/// Source
///

// Source ...
type Source interface {
	QueryString() string
}

// Table ...
type Table struct {
	Name string
}

// QueryString ...
func (t *Table) QueryString() string {
	return t.Name
}

// Select ...
func (t *Table) Select(f ...Field) SelectQuery {
	return SelectQuery{source: t, fields: f}
}

///
/// Field
///

// Field represents a field in a query
type Field interface {
	QueryString(string) string
	Source() Source
	DataType() string
}

// TableField represents a real field in a table
type TableField struct {
	Parent Source
	Name   string
	Type   string
}

// QueryString ...
func (f TableField) QueryString(alias string) string {
	return alias + `.` + f.Name
}

// Source ...
func (f TableField) Source() Source {
	return f.Parent
}

// DataType ...
func (f TableField) DataType() string {
	return f.Type
}

// CalculatedField is a field created by running functions on a TableField
type CalculatedField struct {
	Action func(string) string
	Field  Field
	Type   string
}

// QueryString ...
func (f CalculatedField) QueryString(alias string) string {
	return f.Action(f.Field.QueryString(alias))
}

// Source ...
func (f CalculatedField) Source() Source {
	return f.Field.Source()
}

// DataType ...
func (f CalculatedField) DataType() string {
	return f.Type
}

// ValueField contains values supplied by the program
type ValueField struct {
	Value interface{}
	Type  string
}

func getType(v interface{}) string {
	if val, ok := v.(driver.Valuer); ok {
		new, err := val.Value()
		if err == nil {
			v = new
		}
	}

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
	return ValueField{v, getType(v)}
}

// QueryString ...
func (f ValueField) QueryString(alias string) string {
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
type Condition struct {
	Fields []Field
	Action func(...string) string
}

// Join are used for joins on tables
type Join struct {
	Type   string
	Fields [2]Field
	New    Source
}
