package qf

import (
	"testing"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/tests/testutil"
)

var (
	c1 = func(_ qb.Driver, _ qb.Alias, _ *qb.ValueList) string { return `A` }
	c2 = func(_ qb.Driver, _ qb.Alias, _ *qb.ValueList) string { return `B` }
)

func TestCase(t *testing.T) {
	c := Case().When(c1, 1).When(c2, 2).Else(3)
	expected := `CASE WHEN A THEN ? WHEN B THEN ? ELSE ? END`

	values := qb.ValueList{}

	sql := c.QueryString(nil, qb.NoAlias(), &values)

	if len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf(`Expected values [1, 2, 3]. Got: %v`, values)
	}

	testutil.Compare(t, expected, sql)
}
