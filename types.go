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
	SQL(b SQLBuilder) (string, []interface{})
}

// Alias generates table aliasses.
// This type is not intended to be used directly
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

// Table represents a table in the database.
// This type is used by qb-generator's generated code and is not intended to be used manually
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
	parts := strings.Split(t.Name, `.`)
	return strings.ToLower(parts[len(parts)-1][0:1])
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

// CTE is a type of subqueries
type CTE struct {
	query SelectQuery
	F     []Field
}

func newCTE(q SelectQuery, fields []*Field) *CTE {
	cte := &CTE{query: q}

	assignFields(&cte.F, cte, q, fields)

	return cte
}

func (cte *CTE) aliasString() string {
	return `ct`
}

// QueryString implements QueryStringer
func (cte *CTE) QueryString(c *Context) string {
	alias := c.Alias(cte)
	if len(alias) > 0 {
		alias = ` ` + alias
	}

	return c.cteName(cte) + alias
}

// With generates the SQL for a WITH statement.
// This function is not intended to be called directly
func (cte *CTE) With(b SQLBuilder) string {
	s, _ := cte.query.getSQL(b, true)
	return b.Context.cteName(cte) + ` AS ` + getSubQuerySQL(s)
}

// Select starts a SELECT query
func (cte *CTE) Select(f ...Field) *SelectBuilder {
	return NewSelectBuilder(f, cte)
}

// SubQuery represents a subquery
type SubQuery struct {
	query SelectQuery
	F     []Field
}

func newSubQuery(q SelectQuery, fields []*Field) *SubQuery {
	sq := &SubQuery{query: q}

	assignFields(&sq.F, sq, q, fields)

	return sq
}

func assignFields(dest *[]Field, parent Source, q SelectQuery, fields []*Field) {
	if fields != nil && len(q.Fields()) != len(fields) {
		panic(`Field count in CTE/SubQuery doesn't match`)
	}

	for k := range q.Fields() {
		f := &TableField{Name: `f` + strconv.Itoa(k), Parent: parent}
		*dest = append(*dest, f)

		if fields != nil {
			*fields[k] = f
		}
	}
}

// QueryString implements QueryStringer
func (t *SubQuery) QueryString(c *Context) string {
	alias := c.Alias(t)
	if len(alias) > 0 {
		alias = ` ` + alias
	}

	sql, v := t.query.getSQL(SQLBuilder{Context: c.clone(AliasGenerator())}, true)
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

// TableField represents a field in a table.
// This type is used by qb-generator's generated code and is not intended to be used manually
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

// Context contains all the data needed to build parts of a query.
// This type is not intended to be used directly
type Context struct {
	Driver   Driver
	alias    Alias
	Values   *[]interface{}
	cteNames map[*CTE]string
	cteCount *int
	CTEs     *[]*CTE
}

func (c *Context) cteName(cte *CTE) string {
	if v, ok := c.cteNames[cte]; ok {
		return v
	}

	*c.CTEs = append(*c.CTEs, cte)
	*c.cteCount++

	c.cteNames[cte] = `cte` + strconv.Itoa(*c.cteCount)
	return c.cteNames[cte]
}

// Add adds a value to Values
func (c *Context) Add(v ...interface{}) {
	*c.Values = append(*c.Values, v...)
}

// Alias returns an alias for the given Source
func (c *Context) Alias(src Source) string {
	return c.alias.Get(src)
}

func (c *Context) clone(alias Alias) *Context {
	new := *c
	new.alias = alias

	var values []interface{}
	new.Values = &values

	return &new
}

// NewContext returns a new *Context
func NewContext(d Driver, a Alias) *Context {
	values, count, ctes := []interface{}{}, 0, []*CTE{}
	return &Context{d, a, &values, make(map[*CTE]string), &count, &ctes}
}
