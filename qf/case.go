package qf

import "git.ultraware.nl/NiseVoid/qb"

type when struct {
	C qb.Condition
	F qb.Field
}

// WhenList contains a list of CASE when statements
type WhenList []when

// Case returns a type that allows you to build a when statement
func Case() WhenList {
	return WhenList{}
}

// When adds a statement to the list
func (l WhenList) When(c qb.Condition, v interface{}) WhenList {
	var (
		f  qb.Field
		ok bool
	)
	if f, ok = v.(qb.Field); !ok {
		f = qb.Value(v)
	}

	return append(l, when{C: c, F: f})
}

// Else returns a valid Field to finish the case
func (l WhenList) Else(v interface{}) CaseField {
	var (
		f  qb.Field
		ok bool
	)
	if f, ok = v.(qb.Field); !ok {
		f = qb.Value(v)
	}

	return CaseField{When: []when(l), Else: f}
}

// CaseField is a qb.Field that generates a case statement
type CaseField struct {
	When []when
	Else qb.Field
}

// QueryString returns a string for use in queries
func (f CaseField) QueryString(c *qb.Context) string {
	s := `CASE`
	for _, v := range f.When {
		s += ` WHEN ` + v.C(c) + ` THEN ` + v.F.QueryString(c)
	}
	s += ` ELSE ` + f.Else.QueryString(c) + ` END`
	return s
}

// New returns a new DataField using the given value
func (f CaseField) New(v interface{}) qb.DataField {
	return qb.NewDataField(f, v)
}
