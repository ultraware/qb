package qb

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
}

// TableField represents a real field in a table
type TableField struct {
	Parent Source
	Name   string
}

// QueryString ...
func (f TableField) QueryString(alias string) string {
	return alias + `.` + f.Name
}

// Source ...
func (f TableField) Source() Source {
	return f.Parent
}

// CalculatedField is a field created by running functions on a TableField
type CalculatedField struct {
	Action func(string) string
	Field  Field
}

// QueryString ...
func (f CalculatedField) QueryString(alias string) string {
	return f.Action(f.Field.QueryString(alias))
}

// Source ...
func (f CalculatedField) Source() Source {
	return f.Field.Source()
}

// ValueField contains values supplied by the program
type ValueField struct {
	Value interface{}
}

// Value ...
func Value(v interface{}) ValueField {
	return ValueField{v}
}

// QueryString ...
func (f ValueField) QueryString(alias string) string {
	return `?`
}

// Source ...
func (f ValueField) Source() Source {
	return nil
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
