package test

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"

	"github.com/stretchr/testify/assert"
)

var tempTable = qb.Table{Name: `temp`}
var temp = struct {
	A qb.TableField
	B qb.TableField
	C qb.TableField
	qb.Table
}{
	qb.TableField{Parent: &tempTable, Name: `a`, Type: `int`},
	qb.TableField{Parent: &tempTable, Name: `b`, Type: `string`},
	qb.TableField{Parent: &tempTable, Name: `c`, Type: `string`},
	tempTable,
}

var tbltwoTable = qb.Table{Name: `tbltwo`}
var tbltwo = struct {
	C qb.TableField
	qb.Table
}{
	qb.TableField{Parent: &tbltwoTable, Name: `c`, Type: `string`},
	tbltwoTable,
}

var tblthreeTable = qb.Table{Name: `tbl3`}
var tblthree = struct {
	D qb.TableField
	E qb.TableField
	qb.Table
}{
	qb.TableField{Parent: &tblthreeTable, Name: `d`, Type: `int`},
	qb.TableField{Parent: &tblthreeTable, Name: `e`, Type: `string`},
	tblthreeTable,
}

func TestSelect(t *testing.T) {
	assert := assert.New(t)

	q := temp.Select(qf.Min(temp.A), qf.Max(temp.A), qf.Average(tbltwo.C), qf.Count(tblthree.E)).
		LeftJoin(temp.C, tbltwo.C).
		InnerJoin(tblthree.E, tbltwo.C).
		Where(qc.Gte(temp.A, 2)).
		Where(qc.Eq(1, 1)).
		OrderBy(qb.Asc(tblthree.E)).
		GroupBy(tblthree.E)

	s, v := q.SQL()
	t.Log(s, v)

	assert.Equal(`SELECT min(a1.a), max(a1.a), avg(a2.c), count(a3.e) FROM temp a1 LEFT JOIN tbltwo a2 ON (a1.c = a2.c) INNER JOIN tbl3 a3 ON (a3.e = a2.c) WHERE a1.a >= ? AND ? = ? GROUP BY a3.e ORDER BY a3.e ASC`, s, `Incorrect query`)
	assert.Equal(v, []interface{}{2, 1, 1})
}

func BenchmarkSimpleSelectInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = temp.Select(temp.A, temp.B, temp.C)
	}
}

func BenchmarkSimpleSelectSQL(b *testing.B) {
	q := temp.Select(temp.A, temp.B, temp.C)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.SQL()
	}
}

func BenchmarkCommonSelectSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := temp.Select(qf.Min(temp.A), qf.Max(temp.A), qf.Average(tbltwo.C), qf.Count(tblthree.E)).
			LeftJoin(temp.C, tbltwo.C).
			InnerJoin(tblthree.E, tbltwo.C).
			Where(qc.Gte(temp.A, 2)).
			Where(qc.Eq(1, 1)).
			OrderBy(qb.Asc(tblthree.E)).
			GroupBy(tblthree.E)
		_, _ = q.SQL()
	}
}
