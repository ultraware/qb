package qb

import (
	"strconv"
	"strings"
)

type sqlWriter struct {
	sql    string
	indent int
}

func (w *sqlWriter) WriteString(s string) {
	if w.sql == `` || w.sql[len(w.sql)-1] == '\n' {
		w.sql += strings.Repeat(INDENT, w.indent)
	}
	for k, v := range strings.Split(s, "\n") {
		if k > 0 && v != `` {
			w.sql += NEWLINE + strings.Repeat(INDENT, w.indent)
		}
		w.sql += v
	}
	if s[len(s)-1] == '\n' {
		w.sql += NEWLINE
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
	return w.sql
}

type sqlBuilder struct {
	w       sqlWriter
	Context *Context
}

func newSQLBuilder(d Driver, useAlias bool) sqlBuilder {
	alias := NoAlias()
	if useAlias {
		alias = AliasGenerator()
	}
	return sqlBuilder{sqlWriter{}, NewContext(d, alias)}
}

///// Non-statements /////

func (b *sqlBuilder) ToSQL(qs QueryStringer) string {
	return qs.QueryString(b.Context)
}

// List lists the given fields
func (b *sqlBuilder) List(f []Field, withAlias bool) string {
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

// ListDataFields casts DataFields to Field and returns the output of List
func (b *sqlBuilder) ListDataFields(f []DataField, withAlias bool) string {
	fields := make([]Field, len(f))
	for k, v := range f {
		fields[k] = v
	}
	return b.List(fields, withAlias)
}

// Conditions generates valid SQL for the given list of conditions
func (b *sqlBuilder) Conditions(c []Condition, newline bool) {
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

func (b *sqlBuilder) Select(withAlias bool, f ...DataField) {
	b.w.WriteLine(`SELECT ` + b.ListDataFields(f, withAlias))
}

func (b *sqlBuilder) From(src Source) {
	b.w.WriteLine(`FROM ` + b.ToSQL(src))
}

func (b *sqlBuilder) Join(j ...join) {
	b.w.AddIndent()
	defer b.w.SubIndent()

	for _, v := range j {
		b.w.WriteString(string(v.Type) + ` JOIN ` + b.ToSQL(v.New) + ` ON (`)
		b.Conditions(v.Conditions, false)
		b.w.WriteLine(`)`)
	}
}

func (b *sqlBuilder) Where(c ...Condition) {
	if len(c) == 0 {
		return
	}
	b.w.WriteString(`WHERE `)
	b.Conditions(c, true)
}

func (b *sqlBuilder) GroupBy(f ...Field) {
	if len(f) == 0 {
		return
	}
	b.w.WriteLine(`GROUP BY ` + b.List(f, false))
}

func (b *sqlBuilder) Having(c ...Condition) {
	if len(c) == 0 {
		return
	}
	b.w.WriteString(`HAVING `)
	b.Conditions(c, true)
}

func (b *sqlBuilder) OrderBy(o ...FieldOrder) {
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

func (b *sqlBuilder) Limit(i int) {
	if i == 0 {
		return
	}
	b.w.WriteLine(`LIMIT ` + strconv.Itoa(i))
}

func (b *sqlBuilder) Offset(i int) {
	if i == 0 {
		return
	}
	b.w.WriteLine(`OFFSET ` + strconv.Itoa(i))
}

func (b *sqlBuilder) Update(t *Table) {
	_ = t.Name
	b.w.WriteLine(`UPDATE ` + b.ToSQL(t))
}

func (b *sqlBuilder) Set(sets []set) {
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

func (b *sqlBuilder) Delete(t *Table) {
	_ = t.Name
	b.w.WriteLine(`DELETE FROM ` + b.ToSQL(t))
}

func (b *sqlBuilder) Insert(t *Table, f []DataField) {
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

func (b *sqlBuilder) Values(f [][]Field) {
	if len(f) > 1 {
		b.w.WriteLine(`VALUES`)
		b.w.AddIndent()
		defer b.w.SubIndent()
	} else {
		b.w.WriteString(`VALUES `)
	}

	for k, v := range f {
		b.ValueLine(v, k != len(f)-1)
	}
}

func (b *sqlBuilder) ValueLine(f []Field, addComma bool) {
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
