package qf

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
)

func TestAll(t *testing.T) {
	tb := &qb.Table{Name: `test`}

	f1 := &qb.TableField{Name: `A`, Type: `int`, Parent: tb}
	f2 := &qb.TableField{Name: `B`, Type: `int`, Parent: tb}

	check(t, Distinct(f1), `DISTINCT A`, f1.Type, tb)
	check(t, CountAll(), `count(1)`, `int`, nil)
	check(t, Count(f1), `count(A)`, `int`, tb)

	check(t, Sum(f1), `sum(A)`, f1.Type, tb)
	check(t, Average(f1), `avg(A)`, `float`, tb)
	check(t, Min(f1), `min(A)`, f1.Type, tb)
	check(t, Max(f1), `max(A)`, f1.Type, tb)

	check(t, Coalesce(f1, f2), `coalesce(A, B)`, f1.Type, tb)

	check(t, Lower(f1), `lower(A)`, `string`, tb)

	check(t, Now(), `now()`, `time`, nil)

	check(t, Abs(f1), `abs(A)`, f1.Type, tb)
	check(t, Ceil(f1), `ceil(A)`, f1.Type, tb)
	check(t, Floor(f1), `floor(A)`, f1.Type, tb)
	check(t, Round(f1, 2), `round(A, ?)`, f1.Type, tb)

	check(t, Add(f1, f2), `A + B`, f1.Type, tb)
	check(t, Sub(f1, f2), `A - B`, f1.Type, tb)
	check(t, Mult(f1, f2), `A * B`, f1.Type, tb)
	check(t, Div(f1, f2), `A / B`, f1.Type, tb)
	check(t, Mod(f1, f2), `A % B`, f1.Type, tb)
	check(t, Pow(f1, f2), `A ^ B`, f1.Type, tb)
}

func check(t *testing.T, f qb.Field, expectedSQL, expectedType string, src qb.Source) {
	sql := f.QueryString(&qb.NoAlias{}, &qb.ValueList{})

	if sql != expectedSQL {
		t.Errorf(`Incorrect SQL. Expected: "%s". Got: "%s"`, expectedSQL, sql)
		return
	}

	if f.DataType() != expectedType {
		t.Errorf(`Incorrect type. Expected: %s. Got: %s`, expectedType, f.DataType())
	}

	if f.Source() != src {
		t.Errorf(`Source was not saved correctly!`)
	}

	t.Logf(`Success! "%s"`, sql)
}
