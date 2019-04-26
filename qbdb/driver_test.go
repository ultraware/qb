package qbdb

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb/internal/testutil"
)

var db = New(Driver{}, nil)

func TestPrint(t *testing.T) {
	tests := map[interface{}]string{
		2:       `2`,
		uint(3): `3`,
		4.8766:  `4.8766`,
		true:    `t`,
		`abc`:   `@@`,
		nil:     `NULL`,
	}

	c := 0
	for k, v := range tests {
		out, _ := db.printType(k, &c)
		testutil.Compare(t, v, out)
	}
}

func TestPrepareSQL(t *testing.T) {
	test := `SELECT a + ?, ? FROM tbl WHERE id = ?`
	testIn := [][]interface{}{
		{`abc`, true, 1},
		{3, nil, 123},
	}
	testOut := []string{
		`SELECT a + @@, t FROM tbl WHERE id = 1`,
		`SELECT a + 3, NULL FROM tbl WHERE id = 123`,
	}

	for k, v := range testIn {
		out, _ := db.prepareSQL(test, v)
		testutil.Compare(t, testOut[k], out)
	}
}

func BenchmarkPrepareSQL(b *testing.B) {
	test := `SELECT a + ?, ? FROM tbl WHERE str IN (?,?,?,?,?,?)`

	for i := 0; i < b.N; i++ {
		db.prepareSQL(test, []interface{}{`abc`, true, 1, `defg`, `hijk`, `lmnop`, `qrstuvw`, `xyz`})
	}
}
