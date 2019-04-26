package qb

import (
	"reflect"
	"strings"
)

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

// UpdateBuilder builds an UPDATE query
type UpdateBuilder struct {
	t   *Table
	c   []Condition
	set []set
}

// Set adds an update to the SET clause
func (q *UpdateBuilder) Set(f Field, v interface{}) *UpdateBuilder {
	q.set = append(q.set, newSet(f, v))
	return q
}

// Where adds conditions to the WHERE clause
func (q *UpdateBuilder) Where(c ...Condition) *UpdateBuilder {
	q.c = append(q.c, c...)
	return q
}

// SQL returns a query string and a list of values
func (q *UpdateBuilder) SQL(b SQLBuilder) (string, []interface{}) {
	if reflect.TypeOf(b.Context.alias) != reflect.TypeOf(NoAlias()) && strings.HasSuffix(reflect.TypeOf(b.Context.Driver).PkgPath(), `msqb`) {
		b.w.WriteLine(`UPDATE ` + q.t.aliasString())
		b.Set(q.set)
		b.w.WriteLine(`FROM ` + b.SourceToSQL(q.t))
	} else {
		b.Update(q.t)
		b.Set(q.set)
	}
	b.Where(q.c...)

	return b.w.String(), *b.Context.Values
}
