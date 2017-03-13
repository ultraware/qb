package qb

import "reflect"

// SelectQuery ...
type SelectQuery struct {
	source Source
	fields []Field
	where  []Condition
	values []interface{}
	joins  []Join
	order  []FieldOrder
	group  []Field
	tables []Source
}

// Where ...
func (q SelectQuery) Where(c Condition) SelectQuery {
	q.where = append(q.where, c)
	return q
}

// InnerJoin ...
func (q SelectQuery) InnerJoin(f1, f2 Field) SelectQuery {
	return q.join(`INNER`, f1, f2)
}

// CrossJoin ...
func (q SelectQuery) CrossJoin(f1, f2 Field) SelectQuery {
	return q.join(`CROSS`, f1, f2)
}

// LeftJoin ...
func (q SelectQuery) LeftJoin(f1, f2 Field) SelectQuery {
	return q.join(`LEFT`, f1, f2)
}

// RightJoin ...
func (q SelectQuery) RightJoin(f1, f2 Field) SelectQuery {
	return q.join(`RIGHT`, f1, f2)
}

// join ...
func (q SelectQuery) join(t string, f1, f2 Field) SelectQuery {
	if len(q.tables) == 0 {
		q.tables = []Source{q.source}
	}

	var new Source
	c := 0
	for _, v := range q.tables {
		if reflect.DeepEqual(v, f1.Source()) {
			c++
			new = f2.Source()
		}
		if reflect.DeepEqual(v, f2.Source()) {
			c++
			new = f1.Source()
		}
	}

	if c == 0 {
		panic(`Both tables already joined`)
	}
	if c > 1 {
		panic(`None of these tables are present in the query`)
	}

	q.tables = append(q.tables, new)
	q.joins = append(q.joins, Join{t, [2]Field{f1, f2}, new})

	return q
}

// GroupBy ...
func (q SelectQuery) GroupBy(f ...Field) SelectQuery {
	q.group = f
	return q
}

// OrderBy ...
func (q SelectQuery) OrderBy(o ...FieldOrder) SelectQuery {
	q.order = o
	return q
}

// SQL ...
func (q SelectQuery) SQL() (string, []interface{}) {

	b := selectBuilder{newGenerator(), ValueList{}}

	s := b.selectSQL(q.fields)
	s += b.fromSQL(q.source)
	s += b.joinSQL(q.source, q.joins)
	s += b.whereSQL(q.where)
	s += b.groupSQL(q.group)
	s += b.orderSQL(q.order)

	return s, []interface{}(b.list)
}

type selectBuilder struct {
	alias AliasGenerator
	list  ValueList
}

func (b *selectBuilder) listFields(f []Field) string {
	s := ``
	for k, v := range f {
		if k > 0 {
			s += `, `
		}
		s += v.QueryString(&b.alias, &b.list)
	}
	return s
}

func (b *selectBuilder) selectSQL(f []Field) string {
	return `SELECT ` + b.listFields(f)
}

func (b *selectBuilder) fromSQL(s Source) string {
	return ` FROM ` + s.QueryString() + ` ` + b.alias.Get(s) + ` `
}

func (b *selectBuilder) joinSQL(src Source, j []Join) string {
	s := ``

	for _, v := range j {
		f1 := v.Fields[0].QueryString(&b.alias, &b.list)
		f2 := v.Fields[1].QueryString(&b.alias, &b.list)

		s += v.Type + ` JOIN ` + v.New.QueryString() + ` ` + b.alias.Get(v.New) + ` ON (` + f1 + ` = ` + f2 + `) `
	}

	return s
}

func (b *selectBuilder) whereSQL(c []Condition) string {
	s := ``
	for k, v := range c {
		if k == 0 {
			s += `WHERE `
		} else {
			s += ` AND `
		}

		s += v.Action(v.Fields, &b.alias, &b.list)
	}

	return s
}

func (b *selectBuilder) orderSQL(o []FieldOrder) string {
	if len(o) == 0 {
		return ``
	}

	s := ` ORDER BY `
	for k, v := range o {
		if k > 0 {
			s += `, `
		}
		s += v.Field.QueryString(&b.alias, &b.list) + ` ` + v.Order
	}
	return s
}

func (b *selectBuilder) groupSQL(f []Field) string {
	if len(f) == 0 {
		return ``
	}
	return ` GROUP BY ` + b.listFields(f)
}