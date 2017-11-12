package qc

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qbdb"
	"git.ultraware.nl/NiseVoid/qb/tests/testutil"
)

func TestAll(t *testing.T) {
	tb := &qb.Table{Name: `test`}

	f1 := &qb.TableField{Name: `A`, Parent: tb}
	f2 := &qb.TableField{Name: `B`, Parent: tb}

	check(t, Eq(f1, f2), `A = B`)
	check(t, Ne(f1, f2), `A != B`)
	check(t, Gt(f1, f2), `A > B`)
	check(t, Gte(f1, f2), `A >= B`)
	check(t, Lt(f1, f2), `A < B`)
	check(t, Lte(f1, f2), `A <= B`)

	check(t, Between(f1, 1, 2), `A BETWEEN ? AND ?`)
	check(t, IsNull(f1), `A IS NULL`)
	check(t, NotNull(f1), `A IS NOT NULL`)
	check(t, Like(f1, `%a%`), `A LIKE ?`)
	check(t, In(f1, 1, 2, 3), `A IN (?, ?, ?)`)

	check(t, Not(Eq(f1, f2)), `NOT (A = B)`)
	check(t, And(Eq(f1, f2), NotNull(f1)), `(A = B AND A IS NOT NULL)`)
	check(t, Or(Eq(f1, f2), IsNull(f1)), `(A = B OR A IS NULL)`)

	if !fails(func() { In(f1) }) {
		t.Error(`Expected In to panic when no parameters are provided`)
	}
}

var ctx = qb.NewContext(qbdb.Driver{}, qb.NoAlias())

func check(t *testing.T, c qb.Condition, expectedSQL string) {
	sql := c(ctx)

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
