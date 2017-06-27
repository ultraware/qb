package qbdb

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
)

type driver struct{}

func (d driver) ValueString(c int) string {
	return `@@`
}

func (d driver) BoolString(v bool) string {
	if v {
		return `t`
	}
	return `f`
}

func (d driver) UpsertSQL(_ *qb.Table, _ []qb.Field, _ qb.Query) (string, []interface{}) {
	panic(`Not implemented`)
}

func (d driver) ConcatOperator() string {
	panic(`This should not be used`)
}

func (d driver) ExcludedField(string) string {
	panic(`This should not be used`)
}

func (d driver) Returning(q qb.Query, f []qb.Field) (string, []interface{}) {
	panic(`This should not be used`)
}

func (d driver) DateExtract(f string, part string) string {
	panic(`This should not be used`)
}

var db = New(driver{}, nil)

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
		if out != v {
			t.Errorf(`Print failed. Expected: "%s". Got: "%s"`, v, out)
		} else {
			t.Logf(`Print passed. "%s"`, out)
		}
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
		if out != testOut[k] {
			t.Errorf(`PrepareSQL failed. Expected: "%s". Got: "%s"`, v, out)
		} else {
			t.Logf(`PrepareSQL passed. "%s"`, out)
		}
	}
}
