package qb

import "strconv"

type sqlBuilder struct {
	alias  Alias
	values ValueList
}

///// Non-statements /////

// List lists the given fields
func (b *sqlBuilder) List(f []Field, withAlias bool) string {
	s := ``
	for k, v := range f {
		if k > 0 {
			s += COMMA
		}
		s += v.QueryString(b.alias, &b.values)
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
			if !newline {
				s += ` `
			}
			s += `AND `
		}
		s += v(b.alias, &b.values)
		if newline {
			s += "\n"
		}
	}

	return s
}

func eq(f1, f2 Field) Condition {
	return func(ag Alias, vl *ValueList) string {
		return f1.QueryString(ag, vl) + ` = ` + f2.QueryString(ag, vl)
	}
}

///// SQL statements /////

func (b *sqlBuilder) Select(withAlias bool, f ...DataField) string {
	return `SELECT ` + b.ListDataFields(f, withAlias) + "\n"
}

func (b *sqlBuilder) From(src Source) string {
	return `FROM ` + src.QueryString(b.alias, &b.values) + "\n"
}

func (b *sqlBuilder) Join(j ...Join) string {
	s := ``

	for _, v := range j {
		s += v.Type + ` JOIN ` + v.New.QueryString(b.alias, &b.values) + ` ON (`

		v.Conditions = append([]Condition{eq(v.Fields[0], v.Fields[1])}, v.Conditions...)

		s += b.Conditions(v.Conditions, false)

		s += `)` + "\n"
	}

	return s
}

func (b *sqlBuilder) Where(c ...Condition) string {
	if len(c) == 0 {
		return ``
	}
	return `WHERE ` + b.Conditions(c, true)
}

func (b *sqlBuilder) WhereDataField(f []DataField) string {
	c := []Condition{}
	for _, v := range f {
		c = append(c, eq(v, Value(v.Get())))
	}

	return b.Where(c...)
}

func (b *sqlBuilder) GroupBy(f ...Field) string {
	if len(f) == 0 {
		return ``
	}
	return `GROUP BY ` + b.List(f, false) + "\n"
}

func (b *sqlBuilder) OrderBy(o ...FieldOrder) string {
	if len(o) == 0 {
		return ``
	}
	s := `ORDER BY `
	for k, v := range o {
		if k > 0 {
			s += `, `
		}
		s += v.Field.QueryString(b.alias, &b.values) + ` ` + v.Order
	}
	return s + "\n"
}

func (b *sqlBuilder) Limit(i int) string {
	if i == 0 {
		return ``
	}
	return `LIMIT ` + strconv.Itoa(i) + "\n"
}

func (b *sqlBuilder) Offset(i int) string {
	if i == 0 {
		return ``
	}
	return `OFFSET ` + strconv.Itoa(i) + "\n"
}

func (b *sqlBuilder) Update(t *Table) string {
	return `UPDATE ` + t.QueryString(b.alias, &b.values) + "\n"
}

func (b *sqlBuilder) Set(f []DataField) string {
	s := `SET `
	for k, v := range GetUpdatableFields(f) {
		if k > 0 {
			s += COMMA
		}
		s += eq(v, Value(v.Get()))(b.alias, &b.values)
	}
	return s + "\n"
}
