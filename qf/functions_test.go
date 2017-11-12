package qf

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/tests/testutil"
)

func TestAll(t *testing.T) { // nolint: funlen
	tb := &qb.Table{Name: `test`}

	f1 := &qb.TableField{Name: `A`, Parent: tb}
	f2 := &qb.TableField{Name: `B`, Parent: tb}
	f3 := &qb.TableField{Name: `C`, Parent: tb}

	check(t, Cast(f1, qb.Int), `CAST(A AS int)`)
	check(t, Cast(f1, qb.String), `CAST(A AS string)`)

	check(t, Distinct(f1), `DISTINCT A`)
	check(t, CountAll(), `count(1)`)
	check(t, Count(f1), `count(A)`)

	check(t, Sum(f1), `sum(A)`)
	check(t, Average(f1), `avg(A)`)
	check(t, Min(f1), `min(A)`)
	check(t, Max(f1), `max(A)`)

	check(t, Coalesce(f1, f2), `coalesce(A, B)`)

	check(t, Lower(f1), `lower(A)`)
	check(t, Concat(f3, `B`, `A`), `C || ? || ?`)
	check(t, Replace(f1, f2, `C`), `replace(A, B, ?)`)
	check(t, Substring(f1, 1, 4), `substring(A, ?, ?)`)
	check(t, Substring(f1, 1, nil), `substring(A, ?)`)

	check(t, Now(), `now()`)

	check(t, Second(f1), `EXTRACT(second FROM A)`)
	check(t, Minute(f1), `EXTRACT(minute FROM A)`)
	check(t, Hour(f1), `EXTRACT(hour FROM A)`)
	check(t, Day(f1), `EXTRACT(day FROM A)`)
	check(t, Week(f1), `EXTRACT(week FROM A)`)
	check(t, Month(f1), `EXTRACT(month FROM A)`)
	check(t, Year(f1), `EXTRACT(year FROM A)`)

	check(t, Abs(f1), `abs(A)`)
	check(t, Ceil(f1), `ceil(A)`)
	check(t, Floor(f1), `floor(A)`)
	check(t, Round(f1, 2), `round(A, ?)`)

	check(t, Add(f1, f2), `A + B`)
	check(t, Sub(f1, f2), `A - B`)
	check(t, Mult(f1, f2), `A * B`)
	check(t, Div(f1, f2), `A / B`)
	check(t, Mod(f1, f2), `A % B`)
	check(t, Pow(f1, f2), `A ^ B`)

	fails(func() {
		Excluded(f1)
	})
}

var c = qb.NewContext(qbdb.Driver{}, qb.NoAlias())

func check(t *testing.T, f qb.Field, expectedSQL string) {
	sql := f.QueryString(c)

	testutil.Compare(t, expectedSQL, sql)
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
