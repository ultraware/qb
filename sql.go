package qb

import (
	"bytes"
	"strconv"
	"strings"
)

type sqlWriter struct {
	sql    bytes.Buffer
	indent int
}

func (w *sqlWriter) WriteString(s string) {
	if w.sql.Len() == 0 || w.sql.Bytes()[w.sql.Len()-1] == '\n' {
		w.sql.WriteString(strings.Repeat(INDENT, w.indent))
	}
	for k, v := range strings.Split(s, "\n") {
		if k > 0 && v != `` {
			w.sql.WriteString(NEWLINE + strings.Repeat(INDENT, w.indent))
		}
		w.sql.WriteString(v)
	}
	if s[len(s)-1] == '\n' {
		w.sql.WriteString(NEWLINE)
	}
}

func (w *sqlWriter) WriteLine(s string) {
	w.WriteString(s + NEWLINE)
}

func (w *sqlWriter) AddIndent() {
	w.indent++
}

func (w *sqlWriter) SubIndent() {
	w.indent--
}

func (w *sqlWriter) String() string {
	return w.sql.String()
}

// SQLBuilder contains data and methods to generate SQL.
// This type is not intended to be used directly
type SQLBuilder struct {
	w       sqlWriter
	Context *Context
}

// NewSQLBuilder returns a new SQLBuilder
func NewSQLBuilder(d Driver) SQLBuilder {
	alias := NoAlias()

	return SQLBuilder{sqlWriter{}, NewContext(d, alias)}
}

///// Non-statements /////

// ToSQL converts a QueryStringer to a string
func (b *SQLBuilder) ToSQL(qs QueryStringer) string {
	return qs.QueryString(b.Context)
}

// List lists the given fields
func (b *SQLBuilder) List(f []Field, withAlias bool) string {
	s := ``
	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		s += b.ToSQL(v)
		if withAlias {
			s += ` f` + strconv.Itoa(k)
		}
	}

	return s
}

// Conditions generates valid SQL for the given list of conditions
func (b *SQLBuilder) Conditions(c []Condition, newline bool) {
	fn := b.w.WriteString
	if newline {
		fn = b.w.WriteLine
	}

	if len(c) == 0 {
		return
	}
	fn(c[0](b.Context))

	if newline {
		b.w.AddIndent()
		defer b.w.SubIndent()
	}

	for _, v := range c[1:] {
		if !newline {
			fn(` `)
		}
		fn(`AND ` + v(b.Context))
	}
}

func eq(f1, f2 Field) Condition {
	return func(c *Context) string {
		return ConcatQuery(c, f1, ` = `, f2)
	}
}

///// SQL statements /////

// Select generates a SQL SELECT line
func (b *SQLBuilder) Select(withAlias bool, f ...Field) {
	b.w.WriteLine(`SELECT ` + b.List(f, withAlias))
}

// From generates a SQL FROM line
func (b *SQLBuilder) From(src Source) {
	b.w.WriteLine(`FROM ` + b.ToSQL(src))
}

// Join generates SQL JOIN lines
func (b *SQLBuilder) Join(j ...join) {
	b.w.AddIndent()
	defer b.w.SubIndent()

	for _, v := range j {
		b.w.WriteString(string(v.Type) + ` JOIN ` + b.ToSQL(v.New) + ` ON (`)
		b.Conditions(v.Conditions, false)
		b.w.WriteLine(`)`)
	}
}

// Where generates SQL WHERE/AND lines
func (b *SQLBuilder) Where(c ...Condition) {
	if len(c) == 0 {
		return
	}
	b.w.WriteString(`WHERE `)
	b.Conditions(c, true)
}

// GroupBy generates a SQL GROUP BY line
func (b *SQLBuilder) GroupBy(f ...Field) {
	if len(f) == 0 {
		return
	}
	b.w.WriteLine(`GROUP BY ` + b.List(f, false))
}

// Having generates a SQL HAVING line
func (b *SQLBuilder) Having(c ...Condition) {
	if len(c) == 0 {
		return
	}
	b.w.WriteString(`HAVING `)
	b.Conditions(c, true)
}

// OrderBy generates a SQL ORDER BY line
func (b *SQLBuilder) OrderBy(o ...FieldOrder) {
	if len(o) == 0 {
		return
	}
	s := `ORDER BY `
	for k, v := range o {
		if k > 0 {
			s += COMMA
		}
		s += b.ToSQL(v.Field) + ` ` + v.Order
	}
	b.w.WriteLine(s)
}

// Limit generates a SQL LIMIT line
func (b *SQLBuilder) Limit(i int) {
	if i == 0 {
		return
	}
	b.w.WriteLine(`LIMIT ` + strconv.Itoa(i))
}

// Offset generates a SQL OFFSET line
func (b *SQLBuilder) Offset(i int) {
	if i == 0 {
		return
	}
	b.w.WriteLine(`OFFSET ` + strconv.Itoa(i))
}

// Update generates a SQL UPDATE line
func (b *SQLBuilder) Update(t *Table) {
	_ = t.Name
	b.w.WriteLine(`UPDATE ` + b.ToSQL(t))
}

// Set generates a SQL SET line
func (b *SQLBuilder) Set(sets []set) {
	if len(sets) == 0 {
		return
	}
	if len(sets) > 1 {
		b.w.WriteLine(`SET`)
		b.w.AddIndent()
		defer b.w.SubIndent()
	} else {
		b.w.WriteString(`SET `)
	}

	for k, v := range sets {
		comma := `,`
		if k == len(sets)-1 {
			comma = ``
		}
		b.w.WriteLine(eq(v.Field, v.Value)(b.Context) + comma)
	}
}

// Delete generates a SQL DELETE FROM line
func (b *SQLBuilder) Delete(t *Table) {
	_ = t.Name
	b.w.WriteLine(`DELETE FROM ` + b.ToSQL(t))
}

// Insert generates a SQL INSERT line
func (b *SQLBuilder) Insert(t *Table, f []Field) {
	_ = t.Name
	s := ``
	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		s += b.ToSQL(v)
	}
	b.w.WriteLine(`INSERT INTO ` + b.ToSQL(t) + ` (` + s + `)`)
}

// Values generates a SQL VALUES line
func (b *SQLBuilder) Values(f [][]Field) {
	if len(f) > 1 {
		b.w.WriteLine(`VALUES`)
		b.w.AddIndent()
		defer b.w.SubIndent()
	} else {
		b.w.WriteString(`VALUES `)
	}

	for k, v := range f {
		b.valueLine(v, k != len(f)-1)
	}
}

func (b *SQLBuilder) valueLine(f []Field, addComma bool) {
	comma := `,`
	if !addComma {
		comma = ``
	}

	s := ``
	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		s += b.ToSQL(v)
	}
	b.w.WriteLine(`(` + s + `)` + comma)
}
