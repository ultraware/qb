package qc

import "git.ultraware.nl/NiseVoid/qb"

func createCondition(i1, i2 interface{}, action func(f ...string) string) qb.Condition {
	var f1, f2 qb.Field
	var ok bool

	if f1, ok = i1.(qb.Field); !ok {
		f1 = qb.Value(i1)
	}
	if f2, ok = i2.(qb.Field); !ok {
		f2 = qb.Value(i2)
	}

	if f1.DataType() != f2.DataType() {
		panic(`Types don't match. ` + f1.DataType() + ` and ` + f2.DataType())
	}
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: action}
}

// Eq ...
func Eq(f1, f2 interface{}) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` = ` + f[1] })
}

// Gt ...
func Gt(f1, f2 interface{}) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` > ` + f[1] })
}

// Gte ...
func Gte(f1, f2 interface{}) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` >= ` + f[1] })
}

// Lt ...
func Lt(f1, f2 interface{}) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` < ` + f[1] })
}

// Lte ...
func Lte(f1, f2 interface{}) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` <= ` + f[1] })
}
