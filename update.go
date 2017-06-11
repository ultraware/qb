package qb

type set struct {
	Field Field
	Value Field
}

func newSet(f Field, v interface{}) set {
	f1, ok := f.(*TableField)
	if !ok {
		panic(`Cannot use non-table field in update`)
	}
	if f1.ReadOnly {
		panic(`Cannot update read-only field`)
	}
	f2 := MakeField(v)
	return set{f1, f2}
}

// UpdateBuilder ...
type UpdateBuilder struct {
	t   *Table
	c   []Condition
	set []set
}

// Set ...
func (q UpdateBuilder) Set(f Field, v interface{}) UpdateBuilder {
	if v, ok := f.(DataField); ok {
		f = v.Field
	}
	q.set = append(q.set, newSet(f, v))
	return q
}

// Where ...
func (q UpdateBuilder) Where(c Condition) UpdateBuilder {
	q.c = append(q.c, c)
	return q
}

// SQL ...
func (q UpdateBuilder) SQL(d Driver) (string, []interface{}) {
	b := sqlBuilder{d, NoAlias(), nil}
	return b.Update(q.t) + b.Set(q.set) + b.Where(q.c...), b.values
}
