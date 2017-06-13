package qb

import "strconv"

type sqlBuilder struct {
	driver Driver
	alias  Alias
	values ValueList
}

///// Non-statements /////

func (b *sqlBuilder) ToSQL(qs QueryStringer) string {
	return qs.QueryString(b.driver, b.alias, &b.values)
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
func (b *sqlBuilder) Conditions(c []Condition, newline bool) string {
	s := ``
	for k, v := range c {
		if k > 0 {
			if newline {
				s += INDENT
			} else {
				s += ` `
			}
			s += `AND `
		}
		s += v(b.driver, b.alias, &b.values)
		if newline {
			s += NEWLINE
		}
	}

	return s
}

func eq(f1, f2 Field) Condition {
	return func(d Driver, ag Alias, vl *ValueList) string {
		return ConcatQuery(d, ag, vl, f1, ` = `, f2)
	}
}

///// SQL statements /////

func (b *sqlBuilder) Select(withAlias bool, f ...DataField) string {
	return `SELECT ` + b.ListDataFields(f, withAlias) + NEWLINE
}

func (b *sqlBuilder) From(src Source) string {
	return `FROM ` + b.ToSQL(src) + NEWLINE
}

func (b *sqlBuilder) Join(j ...join) string {
	s := ``

	for _, v := range j {
		s += INDENT + string(v.Type) + ` JOIN ` + b.ToSQL(v.New) +
			` ON (` + b.Conditions(v.Conditions, false) + `)` + NEWLINE
	}

	return s
}

func (b *sqlBuilder) Where(c ...Condition) string {
	if len(c) == 0 {
		return ``
	}
	return `WHERE ` + b.Conditions(c, true)
}

func (b *sqlBuilder) GroupBy(f ...Field) string {
	if len(f) == 0 {
		return ``
	}
	return `GROUP BY ` + b.List(f, false) + NEWLINE
}

func (b *sqlBuilder) Having(c ...Condition) string {
	if len(c) == 0 {
		return ``
	}
	return `HAVING ` + b.Conditions(c, true)
}

func (b *sqlBuilder) OrderBy(o ...FieldOrder) string {
	if len(o) == 0 {
		return ``
	}
	s := `ORDER BY `
	for k, v := range o {
		if k > 0 {
			s += COMMA
		}
		s += b.ToSQL(v.Field) + ` ` + v.Order
	}
	return s + NEWLINE
}

func (b *sqlBuilder) Limit(i int) string {
	if i == 0 {
		return ``
	}
	return `LIMIT ` + strconv.Itoa(i) + NEWLINE
}

func (b *sqlBuilder) Offset(i int) string {
	if i == 0 {
		return ``
	}
	return `OFFSET ` + strconv.Itoa(i) + NEWLINE
}

func (b *sqlBuilder) Update(t *Table) string {
	_ = t.Name
	return `UPDATE ` + b.ToSQL(t) + NEWLINE
}

func (b *sqlBuilder) Set(sets []set) string {
	s := `SET `
	for k, v := range sets {
		if k > 0 {
			s += COMMA
		}
		s += eq(v.Field, v.Value)(b.driver, b.alias, &b.values)
	}
	return s + NEWLINE
}

func (b *sqlBuilder) Delete(t *Table) string {
	_ = t.Name
	return `DELETE FROM ` + b.ToSQL(t) + NEWLINE
}

func (b *sqlBuilder) Insert(t *Table, f []DataField) string {
	_ = t.Name
	s := `INSERT INTO ` + b.ToSQL(t) + ` (`
	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		s += b.ToSQL(v)
	}
	return s + `)` + NEWLINE
}

func (b *sqlBuilder) Values(f [][]Field) string {
	s := `VALUES`
	for k, v := range f {
		if k != 0 {
			s += `,`
		}
		if len(f) > 1 {
			s += NEWLINE + INDENT
		} else {
			s += ` `
		}
		s += b.ValueLine(v)
	}
	return s + NEWLINE
}

func (b *sqlBuilder) ValueLine(f []Field) string {
	s := ``
	for k, v := range f {
		if k != 0 {
			s += COMMA
		}
		s += b.ToSQL(v)
	}
	return `(` + s + `)`
}
