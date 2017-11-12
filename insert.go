package qb

// defaultField ...
type defaultField struct{}

// Default uses a field's default value
func Default() Field {
	return defaultField{}
}

// QueryString implements QueryStringer
func (f defaultField) QueryString(_ *Context) string {
	return `DEFAULT`
}

// InsertBuilder builds an INSERT query
type InsertBuilder struct {
	table    *Table
	fields   []Field
	values   [][]Field
	update   Query
	conflict []Field
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
	q.update = query
	q.conflict = conflict

	return q
}

// SQL returns a query string and a list of values
func (q *InsertBuilder) SQL(d Driver) (string, []interface{}) {
	b := newSQLBuilder(d, false)

	b.Insert(q.table, q.fields)
	b.Values(q.values)

	sql := b.w.String()
	if q.update != nil {
		s, v := d.UpsertSQL(q.table, q.conflict, q.update)
		b.Context.Add(v...)
		sql += s
	}

	return sql, b.Context.Values
}
