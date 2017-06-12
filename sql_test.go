package qb

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

var (
	warn = color.New(color.FgHiRed).SprintFunc()
	okay = color.New(color.FgHiGreen).SprintFunc()

	str = color.New(color.FgYellow).SprintFunc()

	notice = color.New(color.FgCyan).SprintFunc()
)

func quoted(s string) string {
	var n string
	if len(strings.Split(s, "\n")) > 1 {
		n = "\n"
	}
	return n + `"` + str(s) + `"`
}

func checkOutput(t *testing.T, expected, out string, needsnewline bool) {
	if needsnewline {
		if !strings.HasSuffix(out, "\n") {
			t.Error(warn(`FAIL!`) + "\n\n" +
				`Got:      ` + quoted(out) + "\n" +
				`Expected a trailing newline`,
			)
			return
		}
		out = strings.TrimSuffix(out, "\n")
	}

	if out != expected {
		t.Error(warn(`FAIL!`) + "\n\n" +
			`Got:      ` + quoted(out) + "\n" +
			`Expected: ` + quoted(expected) + "\n",
		)
	} else {
		t.Log(okay(`PASS`)+`:`, quoted(out))
	}
}

func info(t *testing.T, msg string) {
	t.Log(notice(msg))
}

func testBuilder(alias bool) sqlBuilder {
	b := sqlBuilder{nil, NoAlias(), nil}
	if alias {
		b.alias = AliasGenerator()
	}
	return b
}

// Tables

var testTable = &Table{Name: `tmp`}
var testFieldA = &TableField{Name: `colA`, Parent: testTable}
var testFieldB = &TableField{Name: `colB`, Parent: testTable}

var testTable2 = &Table{Name: `tmp2`}
var testFieldA2 = &TableField{Name: `colA2`, Parent: testTable2}

func NewIntField(f Field) DataField {
	i := 0
	return NewDataField(f, &i)
}

func TestFrom(t *testing.T) {
	info(t, `-- Testing without alias`)
	b := testBuilder(false)

	checkOutput(t, `FROM tmp`, b.From(testTable), true)

	info(t, `-- Testing with alias`)
	b = testBuilder(true)

	checkOutput(t, `FROM tmp t1`, b.From(testTable), true)
}

func TestDelete(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `DELETE FROM tmp`, b.Delete(testTable), true)
}

func TestUpdate(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `UPDATE tmp`, b.Update(testTable), true)
}

func TestInsert(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `INSERT INTO tmp ()`, b.Insert(testTable, nil), true)
	checkOutput(t, `INSERT INTO tmp (colA, colB)`,
		b.Insert(testTable, []DataField{NewIntField(testFieldA), NewIntField(testFieldB)}),
		true)
}

func TestJoin(t *testing.T) {
	b := testBuilder(true)

	_ = b.From(testTable)

	checkOutput(t,
		"\t"+`INNER JOIN tmp2 t2 ON (t1.colA = t2.colA2)`,
		b.Join(join{`INNER`, [2]Field{testFieldA, testFieldA2}, testTable2, nil}),
		true,
	)
	checkOutput(t,
		"\t"+`LEFT JOIN tmp2 t2 ON (t1.colA = t2.colA2 AND a AND b)`,
		b.Join(join{`LEFT`, [2]Field{testFieldA, testFieldA2}, testTable2, []Condition{testCondition, testCondition2}}),
		true,
	)
	checkOutput(t,
		"\t"+`INNER JOIN tmp2 t2 ON (t1.colA = t2.colA2)`+"\n\t"+
			`LEFT JOIN tmp2 t2 ON (t1.colA = t2.colA2 AND a AND b)`,
		b.Join(
			join{`INNER`, [2]Field{testFieldA, testFieldA2}, testTable2, nil},
			join{`LEFT`, [2]Field{testFieldA, testFieldA2}, testTable2, []Condition{testCondition, testCondition2}},
		),
		true,
	)
}

// Fields

func TestSelect(t *testing.T) {
	f1 := NewIntField(testFieldA)
	f2 := NewIntField(testFieldB)

	info(t, `-- Testing without alias`)
	b := testBuilder(false)

	checkOutput(t, `SELECT colA, colB`, b.Select(false, f1, f2), true)

	checkOutput(t, `SELECT colA f0, colB f1`, b.Select(true, f1, f2), true)

	info(t, `-- Testing with alias`)
	b = testBuilder(true)

	checkOutput(t, `SELECT t1.colA, t1.colB`, b.Select(false, f1, f2), true)

	checkOutput(t, `SELECT t1.colA f0, t1.colB f1`, b.Select(true, f1, f2), true)
}

func TestSet(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `SET `, b.Set([]set{}), true)

	checkOutput(t, `SET colA = ?`, b.Set([]set{{testFieldA, Value(1)}}), true)

	checkOutput(t,
		`SET colA = ?, colB = ?`,
		b.Set([]set{{testFieldA, Value(1)}, {testFieldB, Value(3)}}),
		true)
}

// Conditions

var testCondition = func(_ Driver, _ Alias, _ *ValueList) string {
	return `a`
}

var testCondition2 = func(_ Driver, _ Alias, _ *ValueList) string {
	return `b`
}

func TestWhere(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `WHERE a`+"\n\t"+`AND b`, b.Where(testCondition, testCondition2), true)

	checkOutput(t, ``, b.Where(), false)
}

func TestHaving(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `HAVING a`+"\n\t"+`AND b`, b.Having(testCondition, testCondition2), true)

	checkOutput(t, ``, b.Having(), false)
}

// Other

func TestGroupBy(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `GROUP BY colA, colB`, b.GroupBy(testFieldA, testFieldB), true)
	checkOutput(t, ``, b.GroupBy(), false)
}

func TestOrderBy(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `ORDER BY colA ASC, colB DESC`, b.OrderBy(Asc(testFieldA), Desc(testFieldB)), true)
	checkOutput(t, ``, b.OrderBy(), false)
}

func TestLimit(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `LIMIT 2`, b.Limit(2), true)
	checkOutput(t, ``, b.Limit(0), false)
}

func TestOffset(t *testing.T) {
	b := testBuilder(false)

	checkOutput(t, `OFFSET 2`, b.Offset(2), true)
	checkOutput(t, ``, b.Offset(0), false)
}

func TestValues(t *testing.T) {
	b := testBuilder(false)

	line := []Field{Value(1), Value(2), Value(3)}

	checkOutput(t, `VALUES (?, ?, ?)`, b.Values([][]Field{line}), true)
	checkOutput(t, `VALUES`+"\n\t"+
		`(?, ?, ?),`+"\n\t"+
		`(?, ?, ?)`,
		b.Values([][]Field{line, line}),
		true,
	)
}
