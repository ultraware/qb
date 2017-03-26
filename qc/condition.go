package qc

import (
	"fmt"

	"git.ultraware.nl/NiseVoid/qb"
)

func createCondition(v []interface{}, action func([]qb.Field, *qb.AliasGenerator, *qb.ValueList) string) qb.Condition {
	fields := []qb.Field{}

	for _, val := range v {
		if f, ok := val.(qb.Field); ok {
			fields = append(fields, f)
			continue
		}
		fields = append(fields, qb.Value(val))
	}

	if len(fields) == 0 {
		return qb.Condition{Action: action}
	}

	t := fields[0].DataType()
	for _, f := range fields[1:] {
		if t != f.DataType() {
			fmt.Printf("%#v\n", fields)
			panic(`Types don't match. ` + t + ` vs ` + f.DataType())
		}
	}

	return qb.Condition{Fields: fields, Action: action}
}

func concatQuery(ag *qb.AliasGenerator, vl *qb.ValueList, values ...interface{}) string {
	s := ``
	for _, val := range values {
		switch v := val.(type) {
		case (qb.Field):
			s += v.QueryString(ag, vl)
		case (string):
			s += v
		}
	}
	return s
}

// Eq ...
func Eq(f1, f2 interface{}) qb.Condition {
	return createCondition([]interface{}{f1, f2}, func(f []qb.Field, ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f[0], ` = `, f[1])
	})
}

// Gt ...
func Gt(f1, f2 interface{}) qb.Condition {
	return createCondition([]interface{}{f1, f2}, func(f []qb.Field, ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f[0], ` > `, f[1])
	})
}

// Gte ...
func Gte(f1, f2 interface{}) qb.Condition {
	return createCondition([]interface{}{f1, f2}, func(f []qb.Field, ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f[0], ` >= `, f[1])
	})
}

// Lt ...
func Lt(f1, f2 interface{}) qb.Condition {
	return createCondition([]interface{}{f1, f2}, func(f []qb.Field, ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f[0], ` < `, f[1])
	})
}

// Lte ...
func Lte(f1, f2 interface{}) qb.Condition {
	return createCondition([]interface{}{f1, f2}, func(f []qb.Field, ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f[0], ` <= `, f[1])
	})
}

// Like ...
func Like(f1 qb.Field, s string) qb.Condition {
	return createCondition([]interface{}{f1, s}, func(f []qb.Field, ag *qb.AliasGenerator, vl *qb.ValueList) string {
		return concatQuery(ag, vl, f[0], ` LIKE `, f[1])
	})
}
