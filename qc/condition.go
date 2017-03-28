package qc

import (
	"git.ultraware.nl/NiseVoid/qb"
)

func createOperatorCondition(i1, i2 interface{}, operator string) qb.Condition {
	f1 := makeField(i1)
	f2 := makeField(i2)
	return func(ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f1, ` `+operator+` `, f2)
	}
}

// Eq ...
func Eq(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `=`)
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

// Like ...
func Like(f1 qb.Field, s string) qb.Condition {
	f2 := makeField(s)
	return func(ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f1, ` LIKE `, f2)
	}
}
