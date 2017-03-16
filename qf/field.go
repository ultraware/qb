package qf

import "git.ultraware.nl/NiseVoid/qb"

// CalculatedField is a field created by running functions on a TableField
type CalculatedField struct {
	Action func(*qb.AliasGenerator, *qb.ValueList) string
	S      qb.Source
	Type   string
}

// QueryString ...
func (f CalculatedField) QueryString(ag *qb.AliasGenerator, vl *qb.ValueList) string {
	return f.Action(ag, vl)
}

// Source ...
func (f CalculatedField) Source() qb.Source {
	return f.S
}

// DataType ...
func (f CalculatedField) DataType() string {
	return f.Type
}
