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
	qb.TableField{Parent: &tempTable, Name: `a`},
	qb.TableField{Parent: &tempTable, Name: `b`},
	qb.TableField{Parent: &tempTable, Name: `c`},
	tempTable,
}

var tbltwoTable = qb.Table{Name: `tbltwo`}
var tbltwo = struct {
	C qb.TableField
	qb.Table
}{
	qb.TableField{Parent: &tbltwoTable, Name: `c`},
	tbltwoTable,
}

var tblthreeTable = qb.Table{Name: `tbl3`}
var tblthree = struct {
	D qb.TableField
	E qb.TableField
	qb.Table
}{
	qb.TableField{Parent: &tblthreeTable, Name: `d`},
	qb.TableField{Parent: &tblthreeTable, Name: `e`},
	tblthreeTable,
}

func TestSelect(t *testing.T) {
	assert := assert.New(t)

	q := temp.Select(qf.Min(temp.A), qf.Max(temp.A), qf.Average(tbltwo.C), qf.Count(tblthree.E)).
		LeftJoin(temp.B, tbltwo.C).
		InnerJoin(tblthree.D, tbltwo.C).
		Where(qc.GreaterEqual(temp.A, qb.Value(2))).
		Where(qc.Equal(qb.Value(1), qb.Value(1))).
		OrderBy(tblthree.E).
		GroupBy(tblthree.E)

	s, v := q.SQL()
	t.Log(s, v)

	assert.Equal(`SELECT min(a1.a), max(a1.a), avg(a2.c), count(a3.e) FROM temp a1 LEFT JOIN tbltwo a2 ON (a1.b = a2.c) INNER JOIN tbl3 a3 ON (a3.d = a2.c) WHERE a1.a >= ? AND ? = ? GROUP BY a3.e ORDER BY a3.e`, s, `Incorrect query`)
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
	q := temp.Select(qf.Min(temp.A), qf.Max(temp.A), qf.Average(tbltwo.C), qf.Count(tblthree.E)).
		LeftJoin(temp.B, tbltwo.C).
		InnerJoin(tblthree.D, tbltwo.C).
		Where(qc.GreaterEqual(temp.A, qb.Value(2))).
		Where(qc.Equal(qb.Value(1), qb.Value(1))).
		OrderBy(tblthree.E).
		GroupBy(tblthree.E)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.SQL()
	}
}
