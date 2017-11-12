package qc

import (
	"runtime"

	"git.ultraware.nl/NiseVoid/qb"
)

// NewCondition returns a new Condition
func NewCondition(values ...interface{}) qb.Condition {
	return func(c *qb.Context) string {
		return qb.ConcatQuery(c, values...)
	}
}

func useOverride(fallback interface{}, in ...interface{}) qb.Condition {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	fn := runtime.FuncForPC(pc[0])

	return func(c *qb.Context) string {
		return c.Driver.Override().Condition(fn, fallback, in)(c)
	}
}
