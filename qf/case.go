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

	if len(l) > 0 {
		if l[0].F.DataType() != f.DataType() {
			panic(`Return types in case don't match.`)
		}
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

	if len(l) > 0 {
		if l[0].F.DataType() != f.DataType() {
			panic(`Return types in case don't match.`)
		}
	}
	return CaseField{When: []when(l), Else: f}
}

// CaseField is a qb.Field that generates a case statement
type CaseField struct {
	When []when
	Else qb.Field
}

// QueryString returns a string for use in queries
func (f CaseField) QueryString(ag *qb.AliasGenerator, vl *qb.ValueList) string {
	s := `CASE `
	for _, v := range f.When {
		s += `WHEN ` + v.C.Action(v.C.Fields, ag, vl) + ` THEN ` + v.F.QueryString(ag, vl)
	}
	s += ` ELSE ` + f.Else.QueryString(ag, vl) + ` END`
	return s
}

// Source always returns nil
func (f CaseField) Source() qb.Source {
	return nil
}

// DataType returns the type that will be returned by the case
func (f CaseField) DataType() string {
	return f.Else.DataType()
}
