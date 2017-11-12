package qb

// DeleteBuilder builds a DELETE query
type DeleteBuilder struct {
	table *Table
	c     []Condition
}

// SQL returns a query string and a list of values
func (q DeleteBuilder) SQL(d Driver) (string, []interface{}) {
	b := newSQLBuilder(d, false)

	b.Delete(q.table)
	b.Where(q.c...)

	return b.w.String(), b.Context.Values
}
