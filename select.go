package qb

import (
	"reflect"
	"strconv"
)

// SelectQuery ...
type SelectQuery interface {
	SQL() (string, []interface{})
	getSQL(aliasFields bool) (string, []interface{})
	SubQuery() *SubQuery
	Fields() []Field
}

// SelectBuilder ...
type SelectBuilder struct {
	source Source
	fields []Field
	where  []Condition
	values []interface{}
	joins  []Join
	order  []FieldOrder
	group  []Field
	tables []Source
	limit  int
	offset int
}

// Where ...
func (q SelectBuilder) Where(c Condition) SelectBuilder {
	q.where = append(q.where, c)
	return q
}

// InnerJoin ...
func (q SelectBuilder) InnerJoin(f1, f2 Field) SelectBuilder {
	return q.join(`INNER`, f1, f2)
}

// CrossJoin ...
func (q SelectBuilder) CrossJoin(f1, f2 Field) SelectBuilder {
	return q.join(`CROSS`, f1, f2)
}

// LeftJoin ...
func (q SelectBuilder) LeftJoin(f1, f2 Field) SelectBuilder {
	return q.join(`LEFT`, f1, f2)
}

// RightJoin ...
func (q SelectBuilder) RightJoin(f1, f2 Field) SelectBuilder {
	return q.join(`RIGHT`, f1, f2)
}

// join ...
func (q SelectBuilder) join(t string, f1, f2 Field) SelectBuilder {
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
func (q SelectBuilder) GroupBy(f ...Field) SelectBuilder {
	q.group = f
	return q
}

// OrderBy ...
func (q SelectBuilder) OrderBy(o ...FieldOrder) SelectBuilder {
	q.order = o
	return q
}

// Limit ...
func (q SelectBuilder) Limit(i int) SelectBuilder {
	q.limit = i
	return q
}

// Offset ...
func (q SelectBuilder) Offset(i int) SelectBuilder {
	q.offset = i
	return q
}

// SQL ...
func (q SelectBuilder) SQL() (string, []interface{}) {
	return q.getSQL(false)
}

func (q SelectBuilder) getSQL(aliasFields bool) (string, []interface{}) {
	b := selectBuilder{newGenerator(), ValueList{}}

	s := b.selectSQL(q.fields, aliasFields)
	s += b.fromSQL(q.source)
	s += b.joinSQL(q.source, q.joins)
	s += b.whereSQL(q.where)
	s += b.groupSQL(q.group)
	s += b.orderSQL(q.order)
	if q.limit > 0 {
		s += ` LIMIT ` + strconv.Itoa(q.limit)
	}
	if q.offset > 0 {
		s += ` OFFSET ` + strconv.Itoa(q.offset)
	}

	return s, []interface{}(b.list)
}

type selectBuilder struct {
	alias AliasGenerator
	list  ValueList
}

func (b *selectBuilder) listFields(f []Field, aliasFields bool) string {
	s := ``
	for k, v := range f {
		if k > 0 {
			s += `, `
		}
		s += v.QueryString(&b.alias, &b.list)
		if aliasFields {
			s += ` f` + strconv.Itoa(k)
		}
	}
	return s
}

func (b *selectBuilder) selectSQL(f []Field, aliasFields bool) string {
	return `SELECT ` + b.listFields(f, aliasFields)
}

func (b *selectBuilder) fromSQL(s Source) string {
	return ` FROM ` + s.QueryString(&b.alias, &b.list)
}

func (b *selectBuilder) joinSQL(src Source, j []Join) string {
	s := ``

	for _, v := range j {
		f1 := v.Fields[0].QueryString(&b.alias, &b.list)
		f2 := v.Fields[1].QueryString(&b.alias, &b.list)

		s += ` ` + v.Type + ` JOIN ` + v.New.QueryString(&b.alias, &b.list) + ` ON (` + f1 + ` = ` + f2 + `)`
	}

	return s
}

func (b *selectBuilder) whereSQL(c []Condition) string {
	s := ``
	for k, v := range c {
		if k == 0 {
			s += ` WHERE `
		} else {
			s += ` AND `
		}

		s += v(&b.alias, &b.list)
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
	return ` GROUP BY ` + b.listFields(f, false)
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q SelectBuilder) SubQuery() *SubQuery {
	s, v := q.getSQL(true)

	sq := SubQuery{sql: s, values: v}

	for k, v := range q.fields {
		sq.Fields = append(sq.Fields, TableField{Name: `f` + strconv.Itoa(k), Parent: &sq, Type: v.DataType()})
	}

	return &sq
}

// Fields ...
func (q SelectBuilder) Fields() []Field {
	return q.fields
}

////////////////////////////

// CombinedQuery ...
type CombinedQuery struct {
	combineType string
	queries     []SelectQuery
}

func (q CombinedQuery) getSQL(aliasFields bool) (string, []interface{}) {
	s := ``
	values := []interface{}{}
	for k, v := range q.queries {
		var sql string
		var val []interface{}
		if k == 0 {
			sql, val = v.getSQL(aliasFields)
		} else {
			s += ` ` + q.combineType + ` `
			sql, val = v.getSQL(false)
		}
		s += `(` + sql + `)`
		values = append(values, val...)
	}

	return s, values
}

// SQL ...
func (q CombinedQuery) SQL() (string, []interface{}) {
	return q.getSQL(false)
}

// Fields ...
func (q CombinedQuery) Fields() []Field {
	return q.queries[0].Fields()
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q CombinedQuery) SubQuery() *SubQuery {
	s, v := q.getSQL(true)

	sq := SubQuery{sql: s, values: v}

	for k, v := range q.Fields() {
		sq.Fields = append(sq.Fields, TableField{Name: `f` + strconv.Itoa(k), Parent: &sq, Type: v.DataType()})
	}

	return &sq
}

////////////////////////

// UnionAll ...
func UnionAll(q ...SelectQuery) CombinedQuery {
	return CombinedQuery{combineType: `UNION ALL`, queries: q}
}

// Union ...
func Union(q ...SelectQuery) CombinedQuery {
	return CombinedQuery{combineType: `UNION`, queries: q}
}

// ExceptAll ...
func ExceptAll(q1, q2 SelectQuery) CombinedQuery {
	return CombinedQuery{combineType: `EXCEPT ALL`, queries: []SelectQuery{q1, q2}}
}

// Except ...
func Except(q1, q2 SelectQuery) CombinedQuery {
	return CombinedQuery{combineType: `EXCEPT`, queries: []SelectQuery{q1, q2}}
}

// IntersectAll ...
func IntersectAll(q1, q2 SelectQuery) CombinedQuery {
	return CombinedQuery{combineType: `INTERSECT ALL`, queries: []SelectQuery{q1, q2}}
}

// Intersect ...
func Intersect(q1, q2 SelectQuery) CombinedQuery {
	return CombinedQuery{combineType: `INTERSECT`, queries: []SelectQuery{q1, q2}}
}
