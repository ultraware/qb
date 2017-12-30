package qb

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb/tests/testutil"
)

func TestSQLWriter(t *testing.T) {
	w := sqlWriter{}

	w.WriteLine("0\n0")
	w.AddIndent()
	w.WriteString(`1`)
	w.WriteString("-\n1")
	w.AddIndent()
	w.WriteLine(`-`)
	w.WriteString(`2`)
	w.SubIndent()
	w.SubIndent()
	w.WriteLine(`-`)
	w.WriteString(`0`)

	testutil.Compare(t, "0\n0\n\t1-\n\t1-\n\t\t2-\n0", w.String())
}

func newCheckOutput(t *testing.T, b *SQLBuilder) func(bool, string) {
	return func(newline bool, expected string) {
		out := b.w.String()
		b.w = sqlWriter{}

		if newline {
			expected += "\n"
		}

		testutil.Compare(t, expected, out)
	}
}

func info(t *testing.T, msg string) {
	t.Log(testutil.Notice(msg))
}

func testBuilder(t *testing.T, alias bool) (*SQLBuilder, func(bool, string)) {
	b := &SQLBuilder{sqlWriter{}, NewContext(nil, NoAlias())}
	if alias {
		b.Context.alias = AliasGenerator()
	}
	return b, newCheckOutput(t, b)
}

// Tables

var testTable = &Table{Name: `tmp`}
var testFieldA = &TableField{Name: `colA`, Parent: testTable}
var testFieldB = &TableField{Name: `colB`, Parent: testTable}

var testTable2 = &Table{Name: `tmp2`}
var testFieldA2 = &TableField{Name: `colA2`, Parent: testTable2}

func TestFrom(t *testing.T) {
	info(t, `-- Testing without alias`)
	b, check := testBuilder(t, false)

	b.From(testTable)
	check(true, `FROM tmp`)

	info(t, `-- Testing with alias`)
	b, check = testBuilder(t, true)

	b.From(testTable)
	check(true, `FROM tmp t`)
}

func TestDelete(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Delete(testTable)
	check(true, `DELETE FROM tmp`)
}

func TestUpdate(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Update(testTable)
	check(true, `UPDATE tmp`)
}

func TestInsert(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Insert(testTable, nil)
	check(true, `INSERT INTO tmp ()`)

	b.Insert(testTable, []Field{testFieldA, testFieldB})
	check(true, `INSERT INTO tmp (colA, colB)`)
}

func TestJoin(t *testing.T) {
	b, check := testBuilder(t, true)

	b.From(testTable)
	b.w = sqlWriter{}

	b.Join(join{JoinInner, testTable2, []Condition{eq(testFieldA, testFieldA2)}})
	check(true,
		"\t"+`INNER JOIN tmp2 t2 ON (t.colA = t2.colA2)`,
	)

	b.Join(join{JoinLeft, testTable2, []Condition{eq(testFieldA, testFieldA2), testCondition, testCondition2}})
	check(true,
		"\t"+`LEFT JOIN tmp2 t2 ON (t.colA = t2.colA2 AND a AND b)`,
	)

	b.Join(
		join{JoinInner, testTable2, []Condition{eq(testFieldA, testFieldA2)}},
		join{JoinLeft, testTable2, []Condition{eq(testFieldA, testFieldA2), testCondition, testCondition2}},
	)
	check(true,
		"\t"+`INNER JOIN tmp2 t2 ON (t.colA = t2.colA2)`+"\n\t"+
			`LEFT JOIN tmp2 t2 ON (t.colA = t2.colA2 AND a AND b)`,
	)
}

// Fields

func TestSelect(t *testing.T) {
	f1 := testFieldA
	f2 := testFieldB

	info(t, `-- Testing without alias`)
	b, check := testBuilder(t, false)

	b.Select(false, f1, f2)
	check(true, `SELECT colA, colB`)

	b.Select(true, f1, f2)
	check(true, `SELECT colA f0, colB f1`)

	info(t, `-- Testing with alias`)
	b, check = testBuilder(t, true)

	b.Select(false, f1, f2)
	check(true, `SELECT t.colA, t.colB`)

	b.Select(true, f1, f2)
	check(true, `SELECT t.colA f0, t.colB f1`)
}

func TestSet(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Set([]set{})
	check(false, ``)

	b.Set([]set{{testFieldA, Value(1)}})
	check(true, `SET colA = ?`)

	b.Set([]set{{testFieldA, Value(1)}, {testFieldB, Value(3)}})
	check(true, "SET\n\tcolA = ?,\n\tcolB = ?")
}

// Conditions

var testCondition = func(_ *Context) string {
	return `a`
}

var testCondition2 = func(_ *Context) string {
	return `b`
}

func TestWhere(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Where(testCondition, testCondition2)
	check(true, `WHERE a`+"\n\t"+`AND b`)

	b.Where()
	check(false, ``)
}

func TestHaving(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Having(testCondition, testCondition2)
	check(true, `HAVING a`+"\n\t"+`AND b`)

	b.Having()
	check(false, ``)
}

// Other

func TestGroupBy(t *testing.T) {
	b, check := testBuilder(t, false)

	b.GroupBy(testFieldA, testFieldB)
	check(true, `GROUP BY colA, colB`)

	b.GroupBy()
	check(false, ``)
}

func TestOrderBy(t *testing.T) {
	b, check := testBuilder(t, false)

	b.OrderBy(Asc(testFieldA), Desc(testFieldB))
	check(true, `ORDER BY colA ASC, colB DESC`)

	b.OrderBy()
	check(false, ``)
}

func TestLimit(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Limit(2)
	check(true, `LIMIT 2`)

	b.Limit(0)
	check(false, ``)
}

func TestOffset(t *testing.T) {
	b, check := testBuilder(t, false)

	b.Offset(2)
	check(true, `OFFSET 2`)

	b.Offset(0)
	check(false, ``)
}

func TestValues(t *testing.T) {
	b, check := testBuilder(t, false)

	line := []Field{Value(1), Value(2), Value(3)}

	b.Values([][]Field{line})
	check(true, `VALUES (?, ?, ?)`)

	b.Values([][]Field{line, line})
	check(true, `VALUES`+"\n\t"+
		`(?, ?, ?),`+"\n\t"+
		`(?, ?, ?)`,
	)
}
