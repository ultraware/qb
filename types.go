package qb

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

// Driver implements databse-specific features
type Driver interface {
	ValueString(int) string
	BoolString(bool) string
	UpsertSQL(*Table, []Field, Query) (string, []interface{})
	ConcatOperator() string
	ExcludedField(string) string
	Returning(Query, []Field) (string, []interface{})
	DateExtract(f string, p string) string
	TypeName(DataType) string
}

// Query ...
type Query interface {
	SQL(Driver) (string, []interface{})
}

// Alias ...
type Alias interface {
	Get(Source) string
}

// QueryStringer ...
type QueryStringer interface {
	QueryString(Driver, Alias, *ValueList) string
}

///
/// Source
///

// Source ...
type Source interface {
	QueryStringer
	aliasString() string
}

// Table ...
type Table struct {
	Name string
}

// QueryString ...
func (t *Table) QueryString(_ Driver, ag Alias, _ *ValueList) string {
	alias := ag.Get(t)
	if len(alias) > 0 {
		alias = ` ` + alias
	}
	return t.Name + alias
}

// aliasString ...
func (t *Table) aliasString() string {
	return `t`
}

// Select ...
func (t *Table) Select(f ...DataField) SelectBuilder {
	return NewSelectBuilder(f, t)
}

// Delete ...
func (t *Table) Delete(c1 Condition, c ...Condition) Query {
	return DeleteBuilder{t, append(c, c1)}
}

// Update ...
func (t *Table) Update() UpdateBuilder {
	return UpdateBuilder{t, nil, nil}
}

// Insert ...
func (t *Table) Insert(f []DataField) *InsertBuilder {
	q := InsertBuilder{table: t, fields: f}
	return &q
}

// SubQuery ...
type SubQuery struct {
	query SelectQuery
	F     []DataField
}

func newSubQuery(q SelectQuery, f []DataField) *SubQuery {
	sq := &SubQuery{query: q}

	for k := range f {
		sq.F = append(sq.F, TableField{Name: `f` + strconv.Itoa(k), Parent: sq}.New(f[k].Value))
	}

	return sq
}

// QueryString ...
func (t *SubQuery) QueryString(d Driver, ag Alias, vl *ValueList) string {
	alias := ag.Get(t)
	if len(alias) > 0 {
		alias = ` ` + alias
	}
	sql, _ := t.query.getSQL(d, true)
	return getSubQuerySQL(sql) + alias
}

func getSubQuerySQL(sql string) string {
	return `(` + NEWLINE + INDENT + strings.Replace(strings.TrimSuffix(sql, "\n"), "\n", "\n"+INDENT, -1) + NEWLINE + `)`
}

// aliasString ...
func (t *SubQuery) aliasString() string {
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
	QueryStringer
	New(interface{}) DataField
}

// TableField represents a real field in a table
type TableField struct {
	Parent     Source
	Name       string
	ReadOnly   bool
	HasDefault bool
	Primary    bool
}

// QueryString ...
func (f TableField) QueryString(_ Driver, ag Alias, _ *ValueList) string {
	alias := ag.Get(f.Parent)
	if alias != `` {
		alias += `.`
	}
	return alias + f.Name
}

// Copy creates a new instance of the field with a different Parent
func (f TableField) Copy(src Source) *TableField {
	f.Parent = src
	return &f
}

// New returns a new DataField using this field
func (f TableField) New(v interface{}) DataField {
	return NewDataField(&f, v)
}

func getParent(f Field) Source {
	if v, ok := f.(DataField); ok {
		f = v.Field
	}
	v, ok := f.(*TableField)
	if !ok {
		panic(`Invalid use of a non-TableField field`)
	}
	return v.Parent
}

type valueField struct {
	Value interface{}
}

func (f valueField) QueryString(_ Driver, _ Alias, vl *ValueList) string {
	vl.Append(f.Value)
	return VALUE
}

func (f valueField) New(_ interface{}) DataField {
	panic(`Cannot call New on a value`)
}

// Value creats a new Field
func Value(v interface{}) Field {
	return valueField{getValue(v)}
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

///
/// Query types
///

// Condition is used in the Where function
type Condition func(Driver, Alias, *ValueList) string

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
