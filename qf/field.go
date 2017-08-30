package qf

import "git.ultraware.nl/NiseVoid/qb"

// CalculatedField is a field created by running functions on a TableField
type CalculatedField func(c *qb.Context) string

// QueryString ...
func (f CalculatedField) QueryString(c *qb.Context) string {
	return f(c)
}

// New returns a new DataField using the given value
func (f CalculatedField) New(v interface{}) qb.DataField {
	return qb.NewDataField(f, v)
}

func newCalculatedField(args ...interface{}) CalculatedField {
	return func(c *qb.Context) string { return qb.ConcatQuery(c, args...) }
}
