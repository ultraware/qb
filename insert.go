package qb

// defaultField ...
type defaultField struct{}

// QueryString ...
func (f defaultField) QueryString(_ Driver, ag Alias, _ *ValueList) string {
	return `DEFAULT`
}

// New ...
func (f defaultField) New(_ interface{}) DataField {
	panic(`Cannot call New on defaultField`)
}

func shoulddefaultField(v DataField) bool {
	field, ok := v.Field.(*TableField)
	if !ok {
		panic(`Cannot use non-TableField field in insert`)
	}
	return !v.isSet() && field.HasDefault
}

// InsertBuilder is used to create a SQL INSERT query
type InsertBuilder struct {
	table    *Table
	fields   []DataField
	values   [][]Field
	update   Query
	conflict []Field
}

// Add adds a single row to the list of rows
func (q *InsertBuilder) Add() {
	list := make([]Field, len(q.fields))
	for k, v := range q.fields {
		if shoulddefaultField(v) {
			list[k] = defaultField{}
			continue
		}
		list[k] = Value(v.GetValue())
	}
	q.values = append(q.values, list)
}

// Upsert ...
func (q *InsertBuilder) Upsert(query Query, conflict ...Field) {
	q.update = query
	q.conflict = conflict
}

// SQL ...
func (q *InsertBuilder) SQL(d Driver) (string, []interface{}) {
	b := newSQLBuilder(d, false)

	b.Insert(q.table, q.fields)
	b.Values(q.values)

	sql := b.w.String()
	if q.update != nil {
		s, v := d.UpsertSQL(q.table, q.conflict, q.update)
		b.values.Append(v...)
		sql += s
	}

	return sql, b.values
}
