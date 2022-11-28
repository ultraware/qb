package qc // import "git.ultraware.nl/Ultraware/qb/qc"

import (
	"reflect"
	"strings"

	"git.ultraware.nl/Ultraware/qb"
)

func createOperatorCondition(i1, i2 interface{}, operator string) qb.Condition {
	f1 := qb.MakeField(i1)
	f2 := qb.MakeField(i2)
	return NewCondition(f1, ` `+operator+` `, f2)
}

// Eq checks if the values are equal (=)
func Eq(i1, i2 interface{}) qb.Condition {
	return useOverride(eq, i1, i2)
}

func eq(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `=`)
}

// Ne checks if the values are unequal (!=)
func Ne(i1, i2 interface{}) qb.Condition {
	return useOverride(ne, i1, i2)
}

func ne(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `!=`)
}

// Gt checks if i1 is greater than i2 (>)
func Gt(i1, i2 interface{}) qb.Condition {
	return useOverride(gt, i1, i2)
}

func gt(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `>`)
}

// Gte checks if i1 is greater or equal to i2 (>=)
func Gte(i1, i2 interface{}) qb.Condition {
	return useOverride(gte, i1, i2)
}

func gte(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `>=`)
}

// Lt checks if i1 is less than i2 (<)
func Lt(i1, i2 interface{}) qb.Condition {
	return useOverride(lt, i1, i2)
}

func lt(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `<`)
}

// Lte checks if i1 is less than or equal to i2 (<=)
func Lte(i1, i2 interface{}) qb.Condition {
	return useOverride(lte, i1, i2)
}

func lte(i1, i2 interface{}) qb.Condition {
	return createOperatorCondition(i1, i2, `<=`)
}

// Between checks if f1 is between i1 and i2
func Between(f1 qb.Field, i1, i2 interface{}) qb.Condition {
	return useOverride(between, f1, i1, i2)
}

func between(f1 qb.Field, i1, i2 interface{}) qb.Condition {
	f2 := qb.MakeField(i1)
	f3 := qb.MakeField(i2)
	return NewCondition(f1, ` BETWEEN `, f2, ` AND `, f3)
}

// IsNull checks if the field is NULL
func IsNull(f1 qb.Field) qb.Condition {
	return useOverride(isNull, f1)
}

func isNull(f1 qb.Field) qb.Condition {
	return NewCondition(f1, ` IS NULL`)
}

// NotNull checks if the field is not NULL
func NotNull(f1 qb.Field) qb.Condition {
	return useOverride(notNull, f1)
}

func notNull(f1 qb.Field) qb.Condition {
	return NewCondition(f1, ` IS NOT NULL`)
}

// Like checks if the f1 is like s
func Like(f1 qb.Field, s string) qb.Condition {
	return useOverride(like, f1, s)
}

func like(f1 qb.Field, s string) qb.Condition {
	f2 := qb.MakeField(s)
	return NewCondition(f1, ` LIKE `, f2)
}

// In checks if f1 is in the list
func In(f1 qb.Field, args ...interface{}) qb.Condition {
	if len(args) == 0 {
		panic(`Cannot call qc.In with zero in values`)
	}

	return useOverride(in, append([]interface{}{f1}, args...)...)
}

func in(f1 qb.Field, args ...interface{}) qb.Condition {
	list := strings.TrimSuffix(strings.Repeat(`?, `, len(args)), `, `)
	return func(c *qb.Context) string {
		c.Add(args...)
		return qb.ConcatQuery(c, f1, ` IN (`+list+`)`)
	}
}

// InQuery checks if f1 is in the subquery's result
func InQuery(f qb.Field, q qb.SelectQuery) qb.Condition {
	return useOverride(inQuery, f, q)
}

func inQuery(f qb.Field, q qb.SelectQuery) qb.Condition {
	return func(c *qb.Context) string {
		return qb.ConcatQuery(c, f, ` IN `, q)
	}
}

// Exists checks if the subquery returns rows
func Exists(q qb.SelectQuery) qb.Condition {
	return useOverride(exists, q)
}

func exists(q qb.SelectQuery) qb.Condition {
	return func(c *qb.Context) string {
		return qb.ConcatQuery(c, `EXISTS `, q)
	}
}

// Not reverses a boolean (!)
func Not(c qb.Condition) qb.Condition {
	return useOverride(not, c)
}

func not(c qb.Condition) qb.Condition {
	return func(ctx *qb.Context) string {
		return `NOT (` + c(ctx) + `)`
	}
}

func createLogicalCondition(operator string, conditions ...qb.Condition) qb.Condition {
	return func(ctx *qb.Context) string {
		s := strings.Builder{}
		s.WriteString(`(`)
		for k, c := range conditions {
			if k > 0 {
				s.WriteString(` ` + operator + ` `)
			}
			s.WriteString(c(ctx))
		}
		s.WriteString(`)`)
		return s.String()
	}
}

// And requires both conditions to be true
func And(c ...qb.Condition) qb.Condition {
	list := make([]interface{}, len(c))
	for k, v := range c {
		list[k] = reflect.ValueOf(v)
	}

	return useOverride(and, list...)
}

func and(c ...qb.Condition) qb.Condition {
	return createLogicalCondition(`AND`, c...)
}

// Or requires one of the conditions to be true
func Or(c ...qb.Condition) qb.Condition {
	list := make([]interface{}, len(c))
	for k, v := range c {
		list[k] = reflect.ValueOf(v)
	}

	return useOverride(or, list...)
}

func or(c ...qb.Condition) qb.Condition {
	return createLogicalCondition(`OR`, c...)
}
