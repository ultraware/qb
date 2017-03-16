package qf

import "git.ultraware.nl/NiseVoid/qb"

func makeField(i interface{}) qb.Field {
	if f, ok := i.(qb.Field); ok {
		return f
	}
	return qb.Value(i)
}
