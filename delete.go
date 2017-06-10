package qb

// DeleteBuilder ...
type DeleteBuilder struct {
	table *Table
	c     []Condition
}

// SQL ...
func (q DeleteBuilder) SQL(d Driver) (string, []interface{}) {
	b := sqlBuilder{d, &NoAlias{}, nil}
	return b.Delete(q.table) + b.Where(q.c...), b.values
}
