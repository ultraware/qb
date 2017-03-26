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
	*qb.Table
}{
	qb.TableField{Parent: &tempTable, Name: `a`, Type: `int`},
	qb.TableField{Parent: &tempTable, Name: `b`, Type: `string`},
	qb.TableField{Parent: &tempTable, Name: `c`, Type: `string`},
	&tempTable,
}

var tbltwoTable = qb.Table{Name: `tbltwo`}
var tbltwo = struct {
	C qb.TableField
	*qb.Table
}{
	qb.TableField{Parent: &tbltwoTable, Name: `c`, Type: `string`},
	&tbltwoTable,
}

var tblthreeTable = qb.Table{Name: `tbl3`}
var tblthree = struct {
	D qb.TableField
	E qb.TableField
	*qb.Table
}{
	qb.TableField{Parent: &tblthreeTable, Name: `d`, Type: `int`},
	qb.TableField{Parent: &tblthreeTable, Name: `e`, Type: `string`},
	&tblthreeTable,
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

	assert.Equal(`SELECT min(t1.a), max(t1.a), avg(t2.c), count(t3.e) FROM temp t1 LEFT JOIN tbltwo t2 ON (t1.c = t2.c) INNER JOIN tbl3 t3 ON (t3.e = t2.c) WHERE t1.a >= ? AND ? = ? GROUP BY t3.e ORDER BY t3.e ASC`, s, `Incorrect query`)
	assert.Equal(v, []interface{}{2, 1, 1})
}

func TestSubQuery(t *testing.T) {
	assert := assert.New(t)

	sq := temp.Select(temp.A, temp.B, tbltwo.C).InnerJoin(temp.C, tbltwo.C).SubQuery()
	q := sq.Select(sq.Fields[0], sq.Fields[1], sq.Fields[2], tblthree.E).InnerJoin(sq.Fields[0], tblthree.D)
	s, _ := q.SQL()
	assert.Equal(`SELECT sq1.f0, sq1.f1, sq1.f2, t2.e FROM (SELECT t1.a f0, t1.b f1, t2.c f2 FROM temp t1 INNER JOIN tbltwo t2 ON (t1.c = t2.c)) sq1 INNER JOIN tbl3 t2 ON (sq1.f0 = t2.d)`, s, `Incorrect query`)
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
