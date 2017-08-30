package qb

// DeleteBuilder ...
type DeleteBuilder struct {
	table *Table
	c     []Condition
}

// SQL ...
func (q DeleteBuilder) SQL(d Driver) (string, []interface{}) {
	b := newSQLBuilder(d, false)

	b.Delete(q.table)
	b.Where(q.c...)

	return b.w.String(), b.Context.Values
}
