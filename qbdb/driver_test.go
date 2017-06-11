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

func (d driver) Returning(q qb.Query, f string) (string, []interface{}) {
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
