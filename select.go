package qb

import "strconv"

// SelectQuery ...
type SelectQuery interface {
	SQL() (string, []interface{})
	getSQL(aliasFields bool) (string, []interface{})
	SubQuery() *SubQuery
	Fields() []DataField
}

// SelectBuilder ...
type SelectBuilder struct {
	source Source
	fields []DataField
	where  []Condition
	joins  []Join
	order  []FieldOrder
	group  []Field
	tables []Source
	limit  int
	offset int
}

// NewSelectBuilder ...
func NewSelectBuilder(f []DataField, src Source) SelectBuilder {
	return SelectBuilder{fields: f, source: src}
}

// Where ...
func (q SelectBuilder) Where(c Condition) SelectBuilder {
	q.where = append(q.where, c)
	return q
}

// InnerJoin ...
func (q SelectBuilder) InnerJoin(f1, f2 Field, c ...Condition) SelectBuilder {
	return q.join(`INNER`, f1, f2, c)
}

// CrossJoin ...
func (q SelectBuilder) CrossJoin(f1, f2 Field, c ...Condition) SelectBuilder {
	return q.join(`CROSS`, f1, f2, c)
}

// LeftJoin ...
func (q SelectBuilder) LeftJoin(f1, f2 Field, c ...Condition) SelectBuilder {
	return q.join(`LEFT`, f1, f2, c)
}

// RightJoin ...
func (q SelectBuilder) RightJoin(f1, f2 Field, c ...Condition) SelectBuilder {
	return q.join(`RIGHT`, f1, f2, c)
}

// join ...
func (q SelectBuilder) join(t string, f1, f2 Field, c []Condition) SelectBuilder {
	if len(q.tables) == 0 {
		q.tables = []Source{q.source}
	}

	var new Source
	exists := 0
	for _, v := range q.tables {
		if v == f1.Source() {
			exists++
			new = f2.Source()
		}
		if v == f2.Source() {
			exists++
			new = f1.Source()
		}
	}

	if exists == 0 {
		panic(`None of these tables are present in the query`)
	}
	if exists > 1 {
		panic(`Both tables already joined`)
	}

	q.tables = append(q.tables, new)
	q.joins = append(q.joins, Join{t, [2]Field{f1, f2}, new, c})

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
	b := sqlBuilder{newGenerator(), ValueList{}}

	for _, v := range q.tables {
		_ = b.alias.Get(v)
	}

	s := b.Select(aliasFields, q.fields...) +
		b.From(q.source) +
		b.Join(q.joins...) +
		b.Where(q.where...) +
		b.GroupBy(q.group...) +
		b.OrderBy(q.order...) +
		b.Limit(q.limit) +
		b.Offset(q.offset)

	return s, []interface{}(b.values)
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q SelectBuilder) SubQuery() *SubQuery {
	s, v := q.getSQL(true)

	sq := SubQuery{sql: s, values: v}

	for k, v := range q.fields {
		sq.fields = append(sq.fields, TableField{Name: `f` + strconv.Itoa(k), Parent: &sq, Type: v.DataType()})
	}

	return &sq
}

// Fields ...
func (q SelectBuilder) Fields() []DataField {
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
func (q CombinedQuery) Fields() []DataField {
	return q.queries[0].Fields()
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q CombinedQuery) SubQuery() *SubQuery {
	s, v := q.getSQL(true)

	sq := SubQuery{sql: s, values: v}

	for k, v := range q.Fields() {
		sq.fields = append(sq.fields, TableField{Name: `f` + strconv.Itoa(k), Parent: &sq, Type: v.DataType()})
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
