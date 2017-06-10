package qf

import "git.ultraware.nl/NiseVoid/qb"

// CalculatedField is a field created by running functions on a TableField
type CalculatedField func(qb.Driver, qb.Alias, *qb.ValueList) string

// QueryString ...
func (f CalculatedField) QueryString(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
	return f(d, ag, vl)
}

// New returns a new DataField using the given value
func (f CalculatedField) New(v interface{}) qb.DataField {
	return qb.NewDataField(f, v)
}

func newCalculatedField(args ...interface{}) CalculatedField {
	return func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string { return qb.ConcatQuery(d, ag, vl, args...) }
}
