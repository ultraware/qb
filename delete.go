package qb

// DeleteBuilder builds a DELETE query
type DeleteBuilder struct {
	table *Table
	c     []Condition
}

// SQL returns a query string and a list of values
func (q DeleteBuilder) SQL(b SQLBuilder) (string, []interface{}) {
	b.Delete(q.table)
	b.Where(q.c...)

	return b.w.String(), *b.Context.Values
}
