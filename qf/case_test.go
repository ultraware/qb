package qf

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/tests/testutil"
)

var (
	c1 = func(_ *qb.Context) string { return `A` }
	c2 = func(_ *qb.Context) string { return `B` }
)

func TestCase(t *testing.T) {
	c := Case().When(c1, 1).When(c2, 2).Else(3)
	expected := `CASE WHEN A THEN ? WHEN B THEN ? ELSE ? END`

	ctx := qb.NewContext(nil, qb.NoAlias())

	sql := c.QueryString(ctx)

	if len(*ctx.Values) != 3 || (*ctx.Values)[0] != 1 || (*ctx.Values)[1] != 2 || (*ctx.Values)[2] != 3 {
		t.Errorf(`Expected values [1, 2, 3]. Got: %v`, *ctx.Values)
	}

	testutil.Compare(t, expected, sql)
}
