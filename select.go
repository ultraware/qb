package qb

// SelectQuery ...
type SelectQuery interface {
	SQL(Driver) (string, []interface{})
	getSQL(Driver, bool) (string, []interface{})
	SubQuery() *SubQuery
	Fields() []Field
}

type returningBuilder struct {
	query  Query
	fields []Field
}

// Returning creates a RETURNING or OUTPUT query
func Returning(q Query, f ...Field) SelectQuery {
	return returningBuilder{q, f}
}

func (q returningBuilder) SQL(d Driver) (string, []interface{}) {
	return q.getSQL(d, false)
}

func (q returningBuilder) getSQL(d Driver, aliasFields bool) (string, []interface{}) {
	b := newSQLBuilder(d, false)

	s, v := d.Returning(q.query, q.fields)
	return s, append(v, b.Context.Values...)
}

func (q returningBuilder) Fields() []Field {
	return q.fields
}

func (q returningBuilder) SubQuery() *SubQuery {
	return newSubQuery(q, q.fields)
}

// SelectBuilder ...
type SelectBuilder struct {
	source Source
	fields []Field
	where  []Condition
	joins  []join
	order  []FieldOrder
	group  []Field
	having []Condition
	tables []Source
	limit  int
	offset int
}

// NewSelectBuilder ...
func NewSelectBuilder(f []Field, src Source) *SelectBuilder {
	return &SelectBuilder{fields: f, source: src}
}

// Where ...
func (q *SelectBuilder) Where(c ...Condition) *SelectBuilder {
	q.where = append(q.where, c...)
	return q
}

// InnerJoin ...
func (q *SelectBuilder) InnerJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinInner, f1, f2, c)
}

// CrossJoin ...
func (q *SelectBuilder) CrossJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinCross, f1, f2, c)
}

// LeftJoin ...
func (q *SelectBuilder) LeftJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinLeft, f1, f2, c)
}

// RightJoin ...
func (q *SelectBuilder) RightJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinRight, f1, f2, c)
}

// ManualJoin manually joins a table
// Use this only if you know what you are doing
func (q *SelectBuilder) ManualJoin(t Join, s Source, c ...Condition) *SelectBuilder {
	q.joins = append(q.joins, join{t, s, c})

	q.tables = append(q.tables, s)

	return q
}

// join ...
func (q *SelectBuilder) join(t Join, f1, f2 Field, c []Condition) *SelectBuilder {
	if len(q.tables) == 0 {
		q.tables = []Source{q.source}
	}

	var new Source
	exists := 0
	for _, v := range q.tables {
		if src := getParent(f1); src == v {
			exists++
			new = getParent(f2)
		}
		if src := getParent(f2); src == v {
			exists++
			new = getParent(f1)
		}
	}

	if exists == 0 {
		panic(`None of these tables are present in the query`)
	}
	if exists > 1 {
		panic(`Both tables already joined`)
	}

	return q.ManualJoin(t, new, append(c, eq(f1, f2))...)
}

// GroupBy ...
func (q *SelectBuilder) GroupBy(f ...Field) *SelectBuilder {
	q.group = f
	return q
}

// Having ...
func (q *SelectBuilder) Having(c ...Condition) *SelectBuilder {
	q.having = append(q.having, c...)
	return q
}

// OrderBy ...
func (q *SelectBuilder) OrderBy(o ...FieldOrder) *SelectBuilder {
	q.order = o
	return q
}

// Limit ...
func (q *SelectBuilder) Limit(i int) *SelectBuilder {
	q.limit = i
	return q
}

// Offset ...
func (q *SelectBuilder) Offset(i int) *SelectBuilder {
	q.offset = i
	return q
}

// SQL ...
func (q *SelectBuilder) SQL(d Driver) (string, []interface{}) {
	return q.getSQL(d, false)
}

func (q *SelectBuilder) getSQL(d Driver, aliasFields bool) (string, []interface{}) {
	b := newSQLBuilder(d, true)

	for _, v := range q.tables {
		_ = b.Context.Alias(v)
	}

	b.Select(aliasFields, q.fields...)
	b.From(q.source)
	b.Join(q.joins...)
	b.Where(q.where...)
	b.GroupBy(q.group...)
	b.Having(q.having...)
	b.OrderBy(q.order...)
	b.Limit(q.limit)
	b.Offset(q.offset)

	return b.w.String(), b.Context.Values
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q *SelectBuilder) SubQuery() *SubQuery {
	return newSubQuery(q, q.fields)
}

// Fields ...
func (q *SelectBuilder) Fields() []Field {
	return q.fields
}

////////////////////////////

type combinedQuery struct {
	combineType string
	queries     []SelectQuery
}

func (q combinedQuery) getSQL(d Driver, aliasFields bool) (string, []interface{}) {
	s := ``
	values := []interface{}{}
	for k, v := range q.queries {
		var sql string
		var val []interface{}
		if k == 0 {
			sql, val = v.getSQL(d, aliasFields)
		} else {
			s += ` ` + q.combineType + ` `
			sql, val = v.getSQL(d, false)
		}
		s += getSubQuerySQL(sql)
		values = append(values, val...)
	}

	return s + NEWLINE, values
}

func (q combinedQuery) SQL(d Driver) (string, []interface{}) {
	return q.getSQL(d, false)
}

func (q combinedQuery) Fields() []Field {
	return q.queries[0].Fields()
}

func (q combinedQuery) SubQuery() *SubQuery {
	return newSubQuery(q, q.Fields())
}

////////////////////////

// UnionAll ...
func UnionAll(q ...SelectQuery) SelectQuery {
	return combinedQuery{combineType: `UNION ALL`, queries: q}
}

// Union ...
func Union(q ...SelectQuery) SelectQuery {
	return combinedQuery{combineType: `UNION`, queries: q}
}

// ExceptAll ...
func ExceptAll(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `EXCEPT ALL`, queries: []SelectQuery{q1, q2}}
}

// Except ...
func Except(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `EXCEPT`, queries: []SelectQuery{q1, q2}}
}

// IntersectAll ...
func IntersectAll(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `INTERSECT ALL`, queries: []SelectQuery{q1, q2}}
}

// Intersect ...
func Intersect(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `INTERSECT`, queries: []SelectQuery{q1, q2}}
}
