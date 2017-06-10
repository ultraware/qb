package qf

import "git.ultraware.nl/NiseVoid/qb"

// CalculatedField is a field created by running functions on a TableField
type CalculatedField struct {
	Action func(qb.Driver, qb.Alias, *qb.ValueList) string
	S      qb.Source
	Type   string
}

// QueryString ...
func (f CalculatedField) QueryString(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string {
	return f.Action(d, ag, vl)
}

// Source ...
func (f CalculatedField) Source() qb.Source {
	return f.S
}

// DataType ...
func (f CalculatedField) DataType() string {
	return f.Type
}

// New returns a new DataField using the given value
func (f *CalculatedField) New(v interface{}) qb.DataField {
	return qb.NewDataField(f, v)
}

func newCalculatedField(src qb.Source, t string, args ...interface{}) *CalculatedField {
	return &CalculatedField{
		Action: func(d qb.Driver, ag qb.Alias, vl *qb.ValueList) string { return qb.ConcatQuery(d, ag, vl, args...) },
		S:      src,
		Type:   t,
	}
}
