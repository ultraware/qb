package qf

import "git.ultraware.nl/NiseVoid/qb"

// CalculatedField is a field created by running functions on a TableField
type CalculatedField struct {
	Action func(qb.Alias, *qb.ValueList) string
	S      qb.Source
	Type   string
}

// QueryString ...
func (f CalculatedField) QueryString(ag qb.Alias, vl *qb.ValueList) string {
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

// New returns a new DataField using the given value
func (f *CalculatedField) New(v interface{}) qb.DataField {
	return qb.NewDataField(f, v)
}

func newCalculatedField(src qb.Source, t string, args ...interface{}) *CalculatedField {
	return &CalculatedField{
		Action: func(ag qb.Alias, vl *qb.ValueList) string { return concatQuery(ag, vl, args...) },
		S:      src,
		Type:   t,
	}
}
