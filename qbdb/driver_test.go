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

func TestPrepareSQL(t *testing.T) {
	test := `SELECT a + ?, ? FROM tbl WHERE id = ?`
	testIn := [][]interface{}{
		{`abc`, true, 1},
		{3, nil, 123},
	}
	testOut := []string{
		`SELECT a + @@, @@ FROM tbl WHERE id = @@`,
		`SELECT a + @@, @@ FROM tbl WHERE id = @@`,
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
