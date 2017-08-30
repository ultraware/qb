package qc

import (
	"strings"

	"git.ultraware.nl/NiseVoid/qb"
)

func createOperatorCondition(i1, i2 interface{}, operator string) qb.Condition {
	f1 := qb.MakeField(i1)
	f2 := qb.MakeField(i2)
	return createCondition(f1, ` `+operator+` `, f2)
}

// Eq ...
func Eq(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `=`)
}

// Ne ...
func Ne(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `!=`)
}

// Gt ...
func Gt(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `>`)
}

// Gte ...
func Gte(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `>=`)
}

// Lt ...
func Lt(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `<`)
}

// Lte ...
func Lte(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `<=`)
}

// Between ...
func Between(f1 qb.Field, i1, i2 interface{}) qb.Condition {
	f2 := qb.MakeField(i1)
	f3 := qb.MakeField(i2)
	return createCondition(f1, ` BETWEEN `, f2, ` AND `, f3)
}

// IsNull ...
func IsNull(f1 qb.Field) qb.Condition {
	return createCondition(f1, ` IS NULL`)
}

// NotNull ...
func NotNull(f1 qb.Field) qb.Condition {
	return createCondition(f1, ` IS NOT NULL`)
}

// Like ...
func Like(f1 qb.Field, s string) qb.Condition {
	f2 := qb.MakeField(s)
	return createCondition(f1, ` LIKE `, f2)
}

// In ...
func In(f1 qb.Field, in ...interface{}) qb.Condition {
	if len(in) == 0 {
		panic(`Cannot call qc.In with zero in values`)
	}
	list := strings.TrimSuffix(strings.Repeat(`?, `, len(in)), `, `)
	return func(c *qb.Context) string {
		c.Add(in...)
		return qb.ConcatQuery(c, f1, ` IN (`+list+`)`)
	}
}

// Not ...
func Not(c qb.Condition) qb.Condition {
	return func(ctx *qb.Context) string {
		return `NOT (` + c(ctx) + `)`
	}
}

// createLogicalCondition ...
func createLogicalCondition(operator string, conditions ...qb.Condition) qb.Condition {
	return func(ctx *qb.Context) string {
		s := `(`
		for k, c := range conditions {
			if k > 0 {
				s += ` ` + operator + ` `
			}
			s += c(ctx)
		}
		s += `)`
		return s
	}
}

// And ...
func And(c ...qb.Condition) qb.Condition {
	return createLogicalCondition(`AND`, c...)
}

// Or ...
func Or(c ...qb.Condition) qb.Condition {
	return createLogicalCondition(`OR`, c...)
}
