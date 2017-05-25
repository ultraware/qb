package qf

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
)

var (
	c1 = func(_ qb.Alias, _ *qb.ValueList) string { return `A` }
	c2 = func(_ qb.Alias, _ *qb.ValueList) string { return `B` }
)

func TestCase(t *testing.T) {
	c := Case().When(c1, 1).When(c2, 2).Else(3)
	expected := `CASE WHEN A THEN ? WHEN B THEN ? ELSE ? END`

	values := qb.ValueList{}

	sql := c.QueryString(&qb.NoAlias{}, &values)

	if len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf(`Expected values [1, 2 3]. Got: %v`, values)
	}

	if sql != expected {
		t.Errorf(`Expected: "%s". Got: "%s"`, expected, sql)
	}

	if c.Source() != nil {
		t.Error(`Got non-nil source`)
	}

	if c.DataType() != `int` {
		t.Errorf(`Expected type int. Got type %s.`, c.DataType())
	}

	t.Logf(`Success! "%s", %v`, sql, values)
}

func TestCaseTypeChecks(t *testing.T) {
	f := func() {
		Case().When(c1, 1).When(c2, `abc`)
	}
	if !fails(f) {
		t.Error(`Mismatched types should cause a panic, but it didn't`)
	}

	f = func() {
		Case().When(c1, 1).When(c2, 2).Else(24.8)
	}
	if !fails(f) {
		t.Error(`Mismatched types should cause a panic, but it didn't`)
	}

	t.Log(`Success! Case panicked when called with mismatched types`)
}

func fails(f func()) (failed bool) {
	defer func() {
		if recover() != nil {
			failed = true
		}
	}()

	f()

	return
}
