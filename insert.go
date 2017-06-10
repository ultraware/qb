package qb

// Default ...
type Default struct{}

// QueryString ...
func (f Default) QueryString(_ Driver, ag Alias, _ *ValueList) string {
	return `DEFAULT`
}

// Source ...
func (f Default) Source() Source {
	return nil
}

// DataType ...
func (f Default) DataType() string {
	return ``
}

func shouldDefault(v DataField) bool {
	field, ok := v.Field.(*TableField)
	if !ok {
		panic(`Cannot use non-table field in insert`)
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
		if shouldDefault(v) {
			list[k] = Default{}
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
	b := sqlBuilder{d, &NoAlias{}, nil}
	sql := b.Insert(q.table, q.fields) + b.Values(q.values)
	if q.update != nil {
		s, v := d.UpsertSQL(q.table, q.conflict, q.update)
		b.values.Append(v...)
		sql += s
	}

	return sql, b.values
}
