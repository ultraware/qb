package qb

import (
	"strconv"
	"strings"
)

// Driver implements databse-specific features
type Driver interface {
	ValueString(int) string
	BoolString(bool) string
	UpsertSQL(*Table, []Field, Query) (string, []interface{})
	Returning(Query, []Field) (string, []interface{})
	TypeName(DataType) string
	Override() OverrideMap
}

// Query generates SQL
type Query interface {
	SQL(Driver) (string, []interface{})
}

// Alias generates table aliasses
type Alias interface {
	Get(Source) string
}

// QueryStringer returns a part of a query string
type QueryStringer interface {
	QueryString(*Context) string
}

///
/// Source
///

// Source represents a table or a subquery
type Source interface {
	QueryStringer
	aliasString() string
}

// Table represents a table in the database
type Table struct {
	Name  string
	Alias string
}

// QueryString implements QueryStringer
func (t *Table) QueryString(c *Context) string {
	alias := c.Alias(t)
	if len(alias) > 0 {
		alias = ` ` + alias
	}
	return t.Name + alias
}

func (t *Table) aliasString() string {
	if t.Alias != `` {
		return t.Alias
	}
	return strings.ToLower(t.Name[0:1])
}

// Select starts a SELECT query
func (t *Table) Select(f []Field) *SelectBuilder {
	return NewSelectBuilder(f, t)
}

// Delete starts a DELETE query
func (t *Table) Delete(c1 Condition, c ...Condition) Query {
	return DeleteBuilder{t, append(c, c1)}
}

// Update starts an UPDATE query
func (t *Table) Update() *UpdateBuilder {
	return &UpdateBuilder{t, nil, nil}
}

// Insert starts an INSERT query
func (t *Table) Insert(f []Field) *InsertBuilder {
	q := InsertBuilder{table: t, fields: f}
	return &q
}

// SubQuery represents a subquery
type SubQuery struct {
	query SelectQuery
	F     []Field
}

func newSubQuery(q SelectQuery, f []Field) *SubQuery {
	sq := &SubQuery{query: q}

	for k := range f {
		sq.F = append(sq.F, &TableField{Name: `f` + strconv.Itoa(k), Parent: sq})
	}

	return sq
}

// QueryString implements QueryStringer
func (t *SubQuery) QueryString(c *Context) string {
	alias := c.Alias(t)
	if len(alias) > 0 {
		alias = ` ` + alias
	}

	sql, v := t.query.getSQL(c.Driver, true)
	c.Add(v...)

	return getSubQuerySQL(sql) + alias
}

func getSubQuerySQL(sql string) string {
	return `(` + NEWLINE + INDENT + strings.Replace(strings.TrimSuffix(sql, "\n"), "\n", "\n"+INDENT, -1) + NEWLINE + `)`
}

func (t *SubQuery) aliasString() string {
	return `sq`
}

// Select starts a SELECT query
func (t *SubQuery) Select(f ...Field) *SelectBuilder {
	return NewSelectBuilder(f, t)
}

///
/// Field
///

// Field represents a field in a query
type Field interface {
	QueryStringer
}

// TableField represents a field in a table
type TableField struct {
	Parent   Source
	Name     string
	ReadOnly bool
	Nullable bool
	Type     DataType
	Size     int
}

// QueryString implements QueryStringer
func (f TableField) QueryString(c *Context) string {
	alias := c.Alias(f.Parent)
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

func getParent(f Field) Source {
	if v, ok := f.(*TableField); ok {
		return v.Parent
	}
	panic(`Invalid use of a non-TableField field`)
}

type valueField struct {
	Value interface{}
}

func (f valueField) QueryString(c *Context) string {
	c.Add(f.Value)
	return VALUE
}

// Value creats a new Field
func Value(v interface{}) Field {
	return valueField{v}
}

///
/// Query types
///

// Condition is used in the Where function
type Condition func(c *Context) string

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

// Context contains all the data needed to build parts of a query
type Context struct {
	Driver Driver
	alias  Alias
	Values []interface{}
}

// Add adds a value to Values
func (c *Context) Add(v ...interface{}) {
	c.Values = append(c.Values, v...)
}

// Alias returns an alias for the given Source
func (c *Context) Alias(src Source) string {
	return c.alias.Get(src)
}

// NewContext returns a new *Context
func NewContext(d Driver, a Alias) *Context {
	return &Context{d, a, nil}
}
