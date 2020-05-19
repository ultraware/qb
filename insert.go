package qb

type defaultField struct{}

// Default uses a field's default value
func Default() Field {
	return defaultField{}
}

// QueryString implements Field
func (f defaultField) QueryString(_ *Context) string {
	return `DEFAULT`
}

// InsertBuilder builds an INSERT query
type InsertBuilder struct {
	table  *Table
	fields []Field
	values [][]Field

	conflict       []Field
	update         Query
	ignoreConflict bool
}

// Values adds values to the query
func (q *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	if len(values) != len(q.fields) {
		panic(`Number of values has to match the number of fields`)
	}

	list := make([]Field, len(values))
	for k, v := range values {
		list[k] = MakeField(v)
	}

	q.values = append(q.values, list)

	return q
}

// Upsert turns the INSERT query into an upsert query, only usable if your driver supports it
func (q *InsertBuilder) Upsert(query Query, conflict ...Field) *InsertBuilder {
	if q.ignoreConflict {
		panic(`can't upsert and ignore conflicts at the same time`)
	}

	q.update = query
	q.conflict = conflict

	return q
}

// IgnoreConflict ignores conflicts from the insert query
func (q *InsertBuilder) IgnoreConflict(conflict ...Field) *InsertBuilder {
	if q.update != nil {
		panic(`can't upsert and ignore conflicts at the same time`)
	}

	q.ignoreConflict = true
	q.conflict = conflict

	return q
}

// SQL returns a query string and a list of values
func (q *InsertBuilder) SQL(b SQLBuilder) (string, []interface{}) {
	if len(q.values) == 0 {
		panic(`Cannot exectue insert without values`)
	}

	b.Insert(q.table, q.fields)
	b.Values(q.values)

	sql := b.w.String()
	if q.update != nil {
		s, v := b.Context.Driver.UpsertSQL(q.table, q.conflict, q.update)
		b.Context.Add(v...)
		sql += s
	} else if q.ignoreConflict {
		s, v := b.Context.Driver.IgnoreConflictSQL(q.table, q.conflict)
		b.Context.Add(v...)
		sql += s
	}

	return sql, *b.Context.Values
}
