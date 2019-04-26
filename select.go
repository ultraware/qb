package qb

// SelectQuery represents a query that returns data
type SelectQuery interface {
	Query
	getSQL(SQLBuilder, bool) (string, []interface{})
	SubQuery(...*Field) *SubQuery
	CTE(...*Field) *CTE
	Fields() []Field
}

// ReturningBuilder builds a query with a RETURNING statement
type ReturningBuilder struct {
	Query  Query
	fields []Field
}

// Returning creates a RETURNING or OUTPUT query
func Returning(q Query, f ...Field) SelectQuery {
	return ReturningBuilder{q, f}
}

// SQL returns a query string and a list of values
func (q ReturningBuilder) SQL(b SQLBuilder) (string, []interface{}) {
	return q.getSQL(b, false)
}

func (q ReturningBuilder) getSQL(b SQLBuilder, aliasFields bool) (string, []interface{}) {
	s, v := b.Context.Driver.Returning(b, q.Query, q.fields)
	return s, append(v, *b.Context.Values...)
}

// Fields returns a list of the fields used in the query
func (q ReturningBuilder) Fields() []Field {
	return q.fields
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q ReturningBuilder) SubQuery(fields ...*Field) *SubQuery {
	return newSubQuery(q, fields)
}

// CTE creates a new CTE (WITH) Query
func (q ReturningBuilder) CTE(fields ...*Field) *CTE {
	return newCTE(q, fields)
}

// SelectBuilder builds a SELECT query
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

// NewSelectBuilder retruns a new SelectBuilder
func NewSelectBuilder(f []Field, src Source) *SelectBuilder {
	return &SelectBuilder{fields: f, source: src, tables: []Source{src}}
}

// Where adds conditions to the WHERE clause
func (q *SelectBuilder) Where(c ...Condition) *SelectBuilder {
	q.where = append(q.where, c...)
	return q
}

// InnerJoin adds an INNER JOIN clause to the query
func (q *SelectBuilder) InnerJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinInner, f1, f2, c)
}

// CrossJoin adds a CROSS JOIN clause to the query
func (q *SelectBuilder) CrossJoin(s Source) *SelectBuilder {
	q.joins = append(q.joins, join{JoinCross, s, nil})

	q.tables = append(q.tables, s)

	return q
}

// LeftJoin adds a LEFT JOIN clause to the query
func (q *SelectBuilder) LeftJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinLeft, f1, f2, c)
}

// RightJoin adds a RIGHT JOIN clause to the query
func (q *SelectBuilder) RightJoin(f1, f2 Field, c ...Condition) *SelectBuilder {
	return q.join(JoinRight, f1, f2, c)
}

// ManualJoin manually joins a table
// Only use this if you know what you are doing
func (q *SelectBuilder) ManualJoin(t Join, s Source, c ...Condition) *SelectBuilder {
	q.joins = append(q.joins, join{t, s, c})

	q.tables = append(q.tables, s)

	return q
}

func (q *SelectBuilder) join(t Join, f1, f2 Field, c []Condition) *SelectBuilder {
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

// GroupBy adds a GROUP BY clause to the query
func (q *SelectBuilder) GroupBy(f ...Field) *SelectBuilder {
	q.group = f
	return q
}

// Having adds a HAVING clause to the query
func (q *SelectBuilder) Having(c ...Condition) *SelectBuilder {
	q.having = append(q.having, c...)
	return q
}

// OrderBy adds a ORDER BY clause to the query
func (q *SelectBuilder) OrderBy(o ...FieldOrder) *SelectBuilder {
	q.order = o
	return q
}

// Limit adds a LIMIT clause to the query
func (q *SelectBuilder) Limit(i int) *SelectBuilder {
	q.limit = i
	return q
}

// Offset adds a OFFSET clause to the query
func (q *SelectBuilder) Offset(i int) *SelectBuilder {
	q.offset = i
	return q
}

// CTE creates a new CTE (WITH) Query
func (q *SelectBuilder) CTE(fields ...*Field) *CTE {
	return newCTE(q, fields)
}

// SQL returns a query string and a list of values
func (q *SelectBuilder) SQL(b SQLBuilder) (string, []interface{}) {
	return q.getSQL(b, false)
}

func (q *SelectBuilder) getSQL(b SQLBuilder, aliasFields bool) (string, []interface{}) {
	oldAlias := b.Context.alias
	if _, ok := oldAlias.(*noAlias); ok {
		b.Context.alias = AliasGenerator()
	}

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

	b.Context.alias = oldAlias

	return b.w.String(), *b.Context.Values
}

// SubQuery converts the SelectQuery to a SubQuery for use in further queries
func (q *SelectBuilder) SubQuery(fields ...*Field) *SubQuery {
	return newSubQuery(q, fields)
}

// Fields returns a list of the fields used in the query
func (q *SelectBuilder) Fields() []Field {
	return q.fields
}

////////////////////////////

type combinedQuery struct {
	combineType string
	queries     []SelectQuery
}

func (q combinedQuery) getSQL(b SQLBuilder, aliasFields bool) (string, []interface{}) {
	s := ``
	for k, v := range q.queries {
		var sql string

		if k == 0 {
			sql, _ = v.getSQL(b, aliasFields)
		} else {
			s += ` ` + q.combineType + ` `
			sql, _ = v.getSQL(b, false)
		}
		s += getSubQuerySQL(sql)
	}

	return s + NEWLINE, *b.Context.Values
}

func (q combinedQuery) SQL(b SQLBuilder) (string, []interface{}) {
	return q.getSQL(b, false)
}

func (q combinedQuery) Fields() []Field {
	return q.queries[0].Fields()
}

func (q combinedQuery) CTE(fields ...*Field) *CTE {
	return newCTE(q, fields)
}

func (q combinedQuery) SubQuery(fields ...*Field) *SubQuery {
	return newSubQuery(q, fields)
}

////////////////////////

// UnionAll combines queries with an UNION ALL
func UnionAll(q ...SelectQuery) SelectQuery {
	return combinedQuery{combineType: `UNION ALL`, queries: q}
}

// Union combines queries with an UNION
func Union(q ...SelectQuery) SelectQuery {
	return combinedQuery{combineType: `UNION`, queries: q}
}

// ExceptAll combines queries with an EXCEPT ALL
func ExceptAll(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `EXCEPT ALL`, queries: []SelectQuery{q1, q2}}
}

// Except combines queries with an EXCEPT
func Except(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `EXCEPT`, queries: []SelectQuery{q1, q2}}
}

// IntersectAll combines queries with an INTERSECT ALL
func IntersectAll(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `INTERSECT ALL`, queries: []SelectQuery{q1, q2}}
}

// Intersect combines queries with an INTERSECT
func Intersect(q1, q2 SelectQuery) SelectQuery {
	return combinedQuery{combineType: `INTERSECT`, queries: []SelectQuery{q1, q2}}
}
