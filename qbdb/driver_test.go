package qbdb

import (
	"database/sql/driver"
	"testing"

	"git.ultraware.nl/Ultraware/qb/v2/internal/testutil"
)

var database = New(Driver{}, nil).(*db)

type Valuer int

func (v Valuer) Value() (driver.Value, error) {
	return v + 1, nil
}

func TestPrint(t *testing.T) {
	tests := map[interface{}]string{
		2:         `2`,
		uint(3):   `3`,
		4.8766:    `4.8766`,
		true:      `t`,
		`abc`:     `@@`,
		nil:       `NULL`,
		Valuer(5): `@@`,
	}

	c := 0
	for k, v := range tests {
		out, _ := database.printType(k, &c)
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
		out, _ := database.prepareSQL(test, v)
		testutil.Compare(t, testOut[k], out)
	}
}

func BenchmarkPrepareSQL(b *testing.B) {
	test := `SELECT a + ?, ? FROM tbl WHERE str IN (?,?,?,?,?,?)`

	for i := 0; i < b.N; i++ {
		database.prepareSQL(test, []interface{}{`abc`, true, 1, `defg`, `hijk`, `lmnop`, `qrstuvw`, `xyz`})
	}
}
