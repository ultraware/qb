package qc

import "git.ultraware.nl/NiseVoid/qb"

func makeField(i interface{}) qb.Field {
	if f, ok := i.(qb.Field); ok {
		return f
	}
	return qb.Value(i)
}

func concatQuery(ag qb.Alias, vl *qb.ValueList, values ...interface{}) string {
	s := ``
	for _, val := range values {
		switch v := val.(type) {
		case (qb.Field):
			s += v.QueryString(ag, vl)
		case (string):
			s += v
		}
	}
	return s
}
