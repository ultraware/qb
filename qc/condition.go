package qc

import "git.ultraware.nl/NiseVoid/qb"

func createCondition(f1, f2 qb.Field, action func(f ...string) string) qb.Condition {
	if f1.DataType() != f2.DataType() {
		panic(`Types don't match. ` + f1.DataType() + ` and ` + f2.DataType())
	}
	return qb.Condition{Fields: []qb.Field{f1, f2}, Action: action}
}

// Eq ...
func Eq(f1, f2 qb.Field) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` = ` + f[1] })
}

// Gt ...
func Gt(f1, f2 qb.Field) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` > ` + f[1] })
}

// Gte ...
func Gte(f1, f2 qb.Field) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` >= ` + f[1] })
}

// Lt ...
func Lt(f1, f2 qb.Field) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` < ` + f[1] })
}

// Lte ...
func Lte(f1, f2 qb.Field) qb.Condition {
	return createCondition(f1, f2, func(f ...string) string { return f[0] + ` <= ` + f[1] })
}
