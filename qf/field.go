package qf

import (
	"runtime"

	"git.ultraware.nl/NiseVoid/qb"
)

// CalculatedField is a field created by running functions on a TableField
type CalculatedField func(c *qb.Context) string

// QueryString implements qb.Field
func (f CalculatedField) QueryString(c *qb.Context) string {
	return f(c)
}

// NewCalculatedField returns a new CalculatedField
func NewCalculatedField(args ...interface{}) CalculatedField {
	return func(c *qb.Context) string { return qb.ConcatQuery(c, args...) }
}

func useOverride(fallback interface{}, in ...interface{}) qb.Field {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	fn := runtime.FuncForPC(pc[0])

	return CalculatedField(func(c *qb.Context) string {
		return c.Driver.Override().Field(fn, fallback, in).QueryString(c)
	})
}
