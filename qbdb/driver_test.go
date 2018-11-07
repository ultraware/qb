package qbdb

import (
	"testing"

	"github.com/ultraware/qb/tests/testutil"
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
	test := `SELECT a + ?, ? FROM tbl`
	testIn := [][]interface{}{
		{`abc`, true},
		{3, nil},
	}
	testOut := []string{
		`SELECT a + @@, t FROM tbl`,
		`SELECT a + 3, NULL FROM tbl`,
	}

	for k, v := range testIn {
		out, _ := db.prepareSQL(test, v)
		testutil.Compare(t, testOut[k], out)
	}
}
