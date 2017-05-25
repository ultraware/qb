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

func checkOutput(t *testing.T, expected, out string, needsnewline bool) {
	if needsnewline {
		if !strings.HasSuffix(out, "\n") {
			t.Error(warn(`FAIL!`) + "\n\n" +
				`Got:      "` + str(out) + `"` + "\n" +
				`Expected a trailing newline`,
			)
			return
		}
		out = strings.TrimSuffix(out, "\n")
	}

	if out != expected {
		t.Error(warn(`FAIL!`) + "\n\n" +
			`Got:      "` + str(out) + `"` + "\n" +
			`Expected: "` + str(expected) + `"` + "\n",
		)
	} else {
		t.Log(okay(`PASS`)+`:`, `"`+str(out)+`"`)
	}
}

func info(t *testing.T, msg string) {
	t.Log(notice(msg))
}

// Tables

var testTable = &Table{Name: `tmp`}
var testFieldA = &TableField{Name: `colA`, Parent: testTable, Type: `int`}
var testFieldB = &TableField{Name: `colB`, Parent: testTable, Type: `int`}

var testTable2 = &Table{Name: `tmp2`}
var testFieldA2 = &TableField{Name: `colA2`, Parent: testTable2, Type: `int`}

func NewIntField(f Field) DataField {
	i := 0
	return NewDataField(f, &i)
}

func TestFrom(t *testing.T) {
	info(t, `-- Testing without alias`)
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `FROM tmp`, b.From(testTable), true)

	info(t, `-- Testing with alias`)
	b = sqlBuilder{newGenerator(), nil}

	checkOutput(t, `FROM tmp t1`, b.From(testTable), true)
}

func TestUpdate(t *testing.T) {
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `UPDATE tmp`, b.Update(testTable), true)
}

func TestJoin(t *testing.T) {
	b := sqlBuilder{newGenerator(), nil}

	_ = b.From(testTable)

	checkOutput(t,
		`INNER JOIN tmp2 t2 ON (t1.colA = t2.colA2)`,
		b.Join(Join{`INNER`, [2]Field{testFieldA, testFieldA2}, testTable2, nil}),
		true,
	)
	checkOutput(t,
		`LEFT JOIN tmp2 t2 ON (t1.colA = t2.colA2 AND a AND b)`,
		b.Join(Join{`LEFT`, [2]Field{testFieldA, testFieldA2}, testTable2, []Condition{testCondition, testCondition2}}),
		true,
	)
	checkOutput(t,
		`INNER JOIN tmp2 t2 ON (t1.colA = t2.colA2)`+"\n"+
			`LEFT JOIN tmp2 t2 ON (t1.colA = t2.colA2 AND a AND b)`,
		b.Join(
			Join{`INNER`, [2]Field{testFieldA, testFieldA2}, testTable2, nil},
			Join{`LEFT`, [2]Field{testFieldA, testFieldA2}, testTable2, []Condition{testCondition, testCondition2}},
		),
		true,
	)
}

// Fields

func TestSelect(t *testing.T) {
	f1 := NewIntField(testFieldA)
	f2 := NewIntField(testFieldB)

	info(t, `-- Testing without alias`)
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `SELECT colA, colB`, b.Select(false, f1, f2), true)

	checkOutput(t, `SELECT colA f0, colB f1`, b.Select(true, f1, f2), true)

	info(t, `-- Testing with alias`)
	b = sqlBuilder{newGenerator(), nil}

	checkOutput(t, `SELECT t1.colA, t1.colB`, b.Select(false, f1, f2), true)

	checkOutput(t, `SELECT t1.colA f0, t1.colB f1`, b.Select(true, f1, f2), true)
}

func TestSet(t *testing.T) {
	f1 := NewIntField(testFieldA)
	f2 := NewIntField(testFieldB)

	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `SET `, b.Set([]DataField{f1, f2}), true)
	checkOutput(t, `SET colA = EXCLUDED.colA, colB = EXCLUDED.colB`, b.SetExcluded([]DataField{f1, f2}), true)

	f2.Set(1)
	checkOutput(t, `SET colB = ?`, b.Set([]DataField{f1, f2}), true)
	f1.Set(3)
	checkOutput(t, `SET colA = ?, colB = ?`, b.Set([]DataField{f1, f2}), true)
}

// Conditions

var testCondition = func(_ Alias, _ *ValueList) string {
	return `a`
}

var testCondition2 = func(_ Alias, _ *ValueList) string {
	return `b`
}

func TestWhere(t *testing.T) {
	b := sqlBuilder{}

	checkOutput(t, `WHERE a`+"\n"+`AND b`, b.Where(testCondition, testCondition2), true)

	checkOutput(t, ``, b.Where(), false)
}

func TestWhereDataField(t *testing.T) {
	b := sqlBuilder{&NoAlias{}, nil}

	f1 := NewIntField(testFieldA)
	f2 := NewIntField(testFieldB)

	checkOutput(t, `WHERE colA = ?`+"\n"+`AND colB = ?`, b.WhereDataField([]DataField{f1, f2}), true)

	checkOutput(t, ``, b.WhereDataField(nil), false)
}

// Other

func TestGroupBy(t *testing.T) {
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `GROUP BY colA, colB`, b.GroupBy(testFieldA, testFieldB), true)
	checkOutput(t, ``, b.GroupBy(), false)
}

func TestOrderBy(t *testing.T) {
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `ORDER BY colA ASC, colB DESC`, b.OrderBy(Asc(testFieldA), Desc(testFieldB)), true)
	checkOutput(t, ``, b.OrderBy(), false)
}

func TestLimit(t *testing.T) {
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `LIMIT 2`, b.Limit(2), true)
	checkOutput(t, ``, b.Limit(0), false)
}

func TestOffset(t *testing.T) {
	b := sqlBuilder{&NoAlias{}, nil}

	checkOutput(t, `OFFSET 2`, b.Offset(2), true)
	checkOutput(t, ``, b.Offset(0), false)
}
